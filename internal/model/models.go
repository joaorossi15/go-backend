package model

type User struct {
	ID       int64
	Username string
	Password []byte
}

type Message struct {
	ID       int64
	SenderID int64
	RecID    int64
	Body     string
}
