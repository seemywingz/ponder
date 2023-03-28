package cmd

type DISCORD_request struct {
	Content string `json:"content"`
}

type DISCORD_response struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
}
