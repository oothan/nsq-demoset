package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	logger "nsq-demoset/app/_applib"
)

var TeleBot *TelegramBot

type TelegramBot struct {
	Token   string
	GroupId string
	URL     string
}

func NewTelegramBot(token, groupID string) {
	TeleBot = &TelegramBot{
		Token:   token,
		GroupId: groupID,
		URL:     fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

func (s *TelegramBot) SendMessage(text string) (bool, error) {
	url := fmt.Sprintf("%s/sendMessage", s.URL)

	body, _ := json.Marshal(map[string]string{
		"chat_id": s.GroupId,
		"text":    text,
	})

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	tData, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	logger.Sugar.Debug(string(tData))

	return true, nil
}
