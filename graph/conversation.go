package graph

import "time"

type Conversation struct {
	ID                    string    `json:"id"`
	Topic                 string    `json:"topic"`
	Preview               string    `json:"preview"`
	HasAttachments        bool      `json:"hasAttachments"`
	LastDeliveredDatetime time.Time `json:"lastDeliveredDateTime"`
	UniqueSenders         []string  `json:"uniqueSenders"`
}
