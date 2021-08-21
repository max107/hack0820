package bot

type msgPayload struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	VideoURL  string `json:"video_url"`
}
