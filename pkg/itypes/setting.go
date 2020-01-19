package itypes

type Setting struct {
	ID    int    `json:"id"`
	Key   string `json:"key" sqlite:"text"`
	Value string `json:"value" sqlite:"text"`
}
