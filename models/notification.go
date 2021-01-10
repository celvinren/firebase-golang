package models

// Notification for notification
type Notification struct {
	Topic string `json:"topic"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
