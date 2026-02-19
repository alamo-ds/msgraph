package graph

import (
	"context"
	"net/http"
	"testing"
)

func TestGroupsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups", `{"value":[{"id":"group1","displayName":"Group 1"}]}`)
	defer server.Close()

	client := newClient(server)

	groups, err := client.Groups().Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}
	if groups[0].ID != "group1" {
		t.Errorf("Expected group ID group1, got %s", groups[0].ID)
	}
}

func TestGroupByIdGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups/group1", `{"id":"group1","displayName":"Group 1"}`)
	defer server.Close()

	client := newClient(server)

	group, err := client.Groups().ById("group1").Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.ID != "group1" {
		t.Errorf("Expected group ID group1, got %s", group.ID)
	}
}
