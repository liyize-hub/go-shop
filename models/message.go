package models

//简单的消息体
type Message struct {
	ProductKey string
	UserID     int64
}

//创建结构体
func NewMessage(uid int64, productkey string) *Message {
	return &Message{UserID: uid, ProductKey: productkey}
}
