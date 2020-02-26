package db

import (
	"github.com/nektro/go-util/alias"
	dbstorage "github.com/nektro/go.dbstorage"
)

type Message struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid" sqlite:"text"`
	At   string `json:"time" sqlite:"text"`
	By   string `json:"author" sqlite:"text"`
	In   string `json:"channel" sqlite:"text"`
	Body string `json:"body" sqlite:"text"`
}

//
//

func CreateMessage(user *User, channel *Channel, body string) *Message {
	dbstorage.InsertsLock.Lock()
	defer dbstorage.InsertsLock.Unlock()
	m := &Message{
		db.QueryNextID(cTableMessagesPrefix + channel.UUID),
		newUUID(),
		alias.T(),
		user.UUID,
		channel.UUID,
		body,
	}
	if channel.HistoryOff {
		return m
	}
	db.QueryPrepared(true, "insert into "+cTableMessagesPrefix+channel.UUID+" values (?,?,?,?,?,?)", m.ID, m.UUID, m.At, m.By, m.In, m.Body)
	db.Build().Up(cTableChannels, "latest_message", m.UUID).Wh("uuid", channel.UUID).Exe()
	return m
}

//
//
