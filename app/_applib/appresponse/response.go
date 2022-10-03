package appresponse

import (
	"nsq-demoset/app/app-services/model"
	"time"
)

type ResponseObj struct {
	ErrCode int64       `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseMessage struct {
	Id        uint64      `json:"id"`
	ToId      uint64      `json:"to_id"`
	FromId    uint64      `json:"from_id"`
	Body      string      `json:"body"`
	ToUser    *model.User `json:"to_user"`
	FromUser  *model.User `json:"from_user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
