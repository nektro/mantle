package ws

func BroadcastMessage(message map[string]string) {
	for _, item := range ConnCache {
		item.Conn.WriteJSON(message)
	}
}
