package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/conf"
	"nsq-demoset/app/app-services/internal/model"
	"nsq-demoset/app/app-services/internal/service"
	"nsq-demoset/app/app-services/internal/utils"
	"time"
)

const (
	writeWait = 1 * time.Minute

	pongWait = 1 * time.Minute

	pingPeriod = (pongWait) / 2

	maxMessageSize = 1024
)

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type wsHandler struct {
	R             *gin.Engine
	socketService *service.SocketService
	userService   model.UserService
}

func NewWSHandler(h *Handler) *wsHandler {
	return &wsHandler{
		R:             h.R,
		socketService: h.socketService,
		userService:   h.userScv,
	}
}

func (ctr *wsHandler) Register() {

	ctr.R.GET("/ws/market", ctr.serveMarketWS)
}

func (ctr *wsHandler) serveMarketWS(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Sugar.Error("Error on opening websocket connection: ", err.Error())
		return
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go ctr.keepAlive(conn)
	ctr.handleWS(conn)
}

func (ctr *wsHandler) handleWS(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	socketId, _ := uuid.NewRandom()
	uniqueId := socketId.String()

	// add new connection
	ctr.socketService.AddConn(uniqueId, conn)

	for {
		msgType, payload, err := conn.ReadMessage()
		if err != nil {
			logger.Sugar.Error(err.Error())
			ctr.socketService.DelConn(uniqueId)
			return
		}

		if msgType == websocket.TextMessage {
			logger.Sugar.Debug("RECV: ", string(payload))
			msg := &readMsg{}
			if err := json.Unmarshal(payload, &msg); err != nil {
				logger.Sugar.Error(err)
				return
			}

			switch msg.MsgType {
			case "user_connect":
				bData, err := json.Marshal(msg.Data)
				if err != nil {
					logger.Sugar.Error(err)
					return
				}

				uc := &userConnect{}
				if err := json.Unmarshal(bData, &uc); err != nil {
					logger.Sugar.Error(err)
					return
				}

				claim, err := utils.ValidateAccessToken(uc.AccessToken, conf.PublicKey)
				if err != nil {
					logger.Sugar.Error(err)
					return
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				user, err := ctr.userService.FindById(ctx, claim.User.Id)
				if err != nil {
					logger.Sugar.Error(err)
					return
				}
				user.Conn[uniqueId] = conn
			}
		}
	}
}

func (ctr *wsHandler) keepAlive(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for range ticker.C {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			logger.Sugar.Error(err.Error())
			return
		}
	}
}

type readMsg struct {
	MsgType string      `json:"msg_type"`
	Data    interface{} `json:"data"`
}

type userConnect struct {
	AccessToken string `json:"access_token"`
}
