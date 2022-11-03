package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/model"
	"sync"
)

type SocketService struct {
	mu            sync.Mutex
	ConnList      map[string]*websocket.Conn
	Users         []*model.User
	marketService *MarketService
}

type SConfig struct {
	MarketService *MarketService
}

func NewSocketService(c *SConfig) *SocketService {
	return &SocketService{
		marketService: c.MarketService,
		ConnList:      make(map[string]*websocket.Conn),
	}
}

func (s *SocketService) AddConn(socketId string, conn *websocket.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ConnList[socketId] = conn
	return nil
}

func (s *SocketService) DelConn(socketId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.ConnList, socketId)
}

// send coin information to all connected websocket connection on every 1s
func (s *SocketService) notifyAll(coin *model.CoinData) {
	for _, conn := range s.ConnList {
		s.mu.Lock()

		data, _ := json.Marshal(coin)
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			logger.Sugar.Error(err)
		}

		s.mu.Unlock()
	}
}

func (s *SocketService) NotifyAll() {
	for coin := range s.marketService.CoinChan {
		s.notifyAll(coin)
	}
}

func (s *SocketService) AddUser(user *model.User) {
	s.Users = append(s.Users, user)
}
