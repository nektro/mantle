package itypes

type RowSetting struct {
	ID    int    `json:"id"`
	Key   string `json:"key" sqlite:"text"`
	Value string `json:"value" sqlite:"text"`
}
