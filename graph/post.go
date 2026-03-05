package graph

import (
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Post struct {
	OdataEtag            string      `json:"@odata.etag,omitempty"`
	ID                   string      `json:"id,omitempty"`
	ChangeKey            string      `json:"changeKey,omitempty"`
	ConversationID       string      `json:"conversationId,omitempty"`
	ConversationThreadID string      `json:"conversationThreadId,omitempty"`
	HasAttachments       bool        `json:"hasAttachments,omitempty"`
	Categories           []string    `json:"categories,omitempty"`
	CreatedDateTime      time.Time   `json:"createdDateTime,omitzero"`
	LastModifiedDateTime time.Time   `json:"lastModifiedDateTime,omitzero"`
	ReceivedDateTime     time.Time   `json:"receivedDateTime,omitzero"`
	From                 Recipient   `json:"from,omitzero"`
	Sender               Recipient   `json:"sender,omitzero"`
	Body                 ItemBody    `json:"body,omitzero"`
	NewParticipants      []Recipient `json:"newParticipants,omitzero"`
}

// sexy wolf growl
func (p Post) RawBody() string {
	// NOTE: doing this with Post instead of ItemBody, since I'm not sure
	// if this logic applies to all ItemBody with type HTML
	return p.Body.rawBody()
}

func (b ItemBody) rawBody() string {
	switch strings.ToLower(b.ContentType) {
	default:
		return b.Content
	case "string":
		return strings.TrimSpace(b.Content)
	case "html":
		return extractRawText(b.Content)
	}
}

func extractRawText(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return ""
	}

	var (
		body     *html.Node
		findBody func(*html.Node)
	)

	// First get the body node, usually the 2nd tag in the response.
	findBody = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}

	findBody(doc)
	if body == nil {
		return ""
	}

	// Next, traverse the body node until the inner-most div is found.
	// This actually contains the text node: for the latest comment, this
	// is usually the second div. For all others, it is usualy the third.
	// If no such div is found, return an empty string.
	for {
		var next *html.Node

		for c := body.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "div" {
				next = c
				break
			}
		}

		if next == nil {
			break
		}

		body = next
	}

	// Lastly, extract the content. Usually this is just a raw string.
	// A table node seems to accompany it, but it is ignored.
	var content strings.Builder
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			content.WriteString(c.Data)
		}
	}

	return content.String()
}
