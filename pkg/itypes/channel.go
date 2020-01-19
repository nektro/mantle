package itypes

type Channel struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid" sqlite:"text"`
	Position    int    `json:"position" sqlite:"int"`
	Name        string `json:"name" sqlite:"text"`
	Description string `json:"description" sqlite:"text"`
}
