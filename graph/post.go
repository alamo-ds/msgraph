package graph

import (
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Post struct {
	OdataEtag            string      `json:"@odata.etag"`
	ID                   string      `json:"id"`
	ChangeKey            string      `json:"changeKey"`
	ConversationID       string      `json:"conversationId"`
	ConversationThreadID string      `json:"conversationThreadId"`
	HasAttachments       bool        `json:"hasAttachments"`
	Categories           []string    `json:"categories"`
	CreatedDateTime      time.Time   `json:"createdDateTime"`
	LastModifiedDateTime time.Time   `json:"lastModifiedDateTime"`
	ReceivedDateTime     time.Time   `json:"receivedDateTime"`
	From                 Recipient   `json:"from"`
	Sender               Recipient   `json:"sender"`
	Body                 ItemBody    `json:"body"`
	NewParticipants      []Recipient `json:"newParticipants"`
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
