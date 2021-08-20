package bot

type statePayload struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
	VideoURL  string `json:"video_url"`
}

type jobPayload struct {
	ChatID    string `json:"chat_id"`
	MessageID string `json:"message_id"`
	VideoURL  string `json:"video_url"`
}
