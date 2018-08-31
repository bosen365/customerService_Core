package model

import "time"

type Room struct {
	RoomCustomer
	RoomKf
	RoomMessages []RoomMessage
	CreateTime   time.Time
}
type RoomCustomer struct {
	CustomerId           string
	CustomerNickName     string
	CustomerHeadImgUrl   string
	CustomerPreviousKfId string
}
type RoomKf struct {
	KfId         int
	KfName       string
	KfHeadImgUrl string
	KfStatus     int
}
type RoomMessage struct {
	Id         string
	Type       string
	Msg        string
	CreateTime time.Time
}
