package mongodb

import "time"

type MongoDb struct {
	Url         string
	DbName      string
	Users       string
	Products    string
	Receipts    string
	Sessions    string
	IdGenerator string
	Accounts    string
}

type Quote struct {
	Data    string `json:"data" bson:"data"`
	Author  string `json:"author" bson:"author"`
	Special bool   `json:"special" bson:"special"`
}

type User struct {
	Username string  `json:"username" bson:"username"`
	Password string  `json:"password" bson:"password"`
	Quotes   []Quote `json:"links" bson:"links"`
}

type Session struct {
	Username  string    `json:"username" bson:"username"`
	Token     string    `json:"token" bson:"token"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}
