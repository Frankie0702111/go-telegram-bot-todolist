package response

type Notify struct {
	TaskID     int64  `json:"task_id"`
	TelegramID int64  `json:"telegram_id"`
	Title      string `json:"title"`
	Time       string `json:"time"`
	Complete   bool   `json:"complete"`
}
