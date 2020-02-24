package db

type Message struct {
	ID     int64  `json:"id"`
	UUID   string `json:"uuid" sqlite:"text"`
	SentAt string `json:"sent_at" sqlite:"text"`
	SentBy string `json:"sent_by" sqlite:"text"`
	Body   string `json:"body" sqlite:"text"`
}
