package data

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
	logger "nsq-demoset/app/_applib"
	"strings"
	"sync"
	"time"
)

type MarketData struct {
	mu   sync.RWMutex
	List map[string]*MiniTicker
}

func NewMarketData() *MarketData {
	m := &MarketData{}
	m.List = make(map[string]*MiniTicker)

	go m.getData("stream.binance.com:9443")

	return m
}

type MiniTicker struct {
	EventType          string      `json:"e"`
	EventTime          uint64      `json:"E"`
	Symbol             string      `json:"s"`
	PriceChange        string      `json:"p"`
	PriceChangePercent string      `json:"P"`
	OpenPrice          string      `json:"o"`
	HighPrice          string      `json:"h"`
	LowPrice           string      `json:"l"`
	LastPrice          string      `json:"c"`
	LastTrade          interface{} `json:"L"`
	OpenTime           interface{} `json:"O"`
	CloseTime          interface{} `json:"C"`
}

func (s *MarketData) getData(addr string) {
	u := url.URL{
		Scheme: "wss",
		Host:   addr,
		Path:   "/wss",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Sugar.Error("Dail : ", err)
	}
	defer conn.Close()

	logger.Sugar.Infof("Connected to %s", u.String())

	payload := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{
			"btcusdt@ticker",
			"trxusdt@ticker",
			"ethusdt@ticker",
		},
		"id": time.Now().Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Sugar.Error(err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	conn.SetWriteDeadline(time.Now().Add(1 * time.Minute))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
			logger.Sugar.Error(err)
			return err
		}
		return nil
	})
	conn.WriteMessage(websocket.TextMessage, data)

	go func() {
		t := time.NewTicker(time.Second * 30)
		for range t.C {
			logger.Sugar.Debug("Ping Message : ")
			conn.SetWriteDeadline(time.Now().Add(time.Minute))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Sugar.Error(err)
				return
			}
		}
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Sugar.Error(err)
			return
		}

		if msgType == websocket.TextMessage {
			s.mu.Lock()
			t := &MiniTicker{}
			if err := json.Unmarshal(msg, t); err != nil {
				logger.Sugar.Error(err)
			}

			symbol := strings.ReplaceAll(t.Symbol, "USDT", "")
			s.List[symbol] = t
			s.mu.Unlock()
		}
	}
}
