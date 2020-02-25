package db

type Message struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid" sqlite:"text"`
	At   string `json:"time" sqlite:"text"`
	By   string `json:"author" sqlite:"text"`
	In   string `json:"channel" sqlite:"text"`
	Body string `json:"body" sqlite:"text"`
}
