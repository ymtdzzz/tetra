package notification

import "time"

type NotificationMsg struct {
	Message string
}

type notificationCleanTickMsg struct {
	t time.Time
}
