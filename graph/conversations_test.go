package graph

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups/group1/threads/thread1/posts", `{"value":[{"id": "thread1","createdDateTime":"2026-02-03T21:50:00Z"}]}`)
	defer server.Close()

	client := newClient(server)

	posts, err := client.Groups().ById("group1").Threads().ById("thread1").Get(context.Background())
	require.NoError(t, err)
	require.Len(t, posts, 1)
}
