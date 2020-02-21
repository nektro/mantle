package itypes

type ChannelPerms struct {
	ID        int    `json:"id"`
	Channel   string `json:"channel" sqlite:"text"`
	Type      int    `json:"p_type" sqlite:"int"`
	Snowflake string `json:"snowflake" sqlite:"text"`
}
