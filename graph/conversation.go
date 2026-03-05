package graph

import "time"

type Conversation struct {
	CcResipients          Recipient   `json:"ccRecipients,omitzero"`
	HasAttachments        bool        `json:"hasAttachments,omitempty"`
	ID                    string      `json:"id,omitempty"`
	IsLocked              bool        `json:"isLocked,omitempty"`
	LastDeliveredDateTime time.Time   `json:"lastDeliveredDateTime,omitzero"`
	Preview               string      `json:"preview,omitempty"`
	Topic                 string      `json:"topic,omitempty"`
	ToRecipients          []Recipient `json:"toRecipients,omitempty"`
	UniqueSenders         []string    `json:"uniqueSenders,omitempty"`
	Posts                 []Post      `json:"posts,omitempty"`
}
