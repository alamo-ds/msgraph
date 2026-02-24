package graph

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups", `{"value":[{"id":"group1","displayName":"Group 1"}]}`)
	defer server.Close()

	client := newClient(server)

	groups, err := client.Groups().Get(context.Background())
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, "group1", groups[0].ID)
}

func TestGroupByIdGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups/group1", `{"id":"group1","displayName":"Group 1"}`)
	defer server.Close()

	client := newClient(server)

	group, err := client.Groups().ById("group1").Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "group1", group.ID)
}

func TestThreadsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups/group1/threads", `{"value":[{"id":"thread1","topic":"test comment thread"}]}`)
	defer server.Close()

	client := newClient(server)

	threads, err := client.Groups().ById("group1").Threads().Get(context.Background())
	require.NoError(t, err)
	require.Len(t, threads, 1)
	require.Equal(t, "thread1", threads[0].ID)
}
