package graph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestPlannerPlansGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups/group1/planner/plans", `{"value":[{"id":"plan1","title":"Plan 1"}]}`)
	defer server.Close()

	client := newClient(server)

	plans, err := client.Groups().ById("group1").Plans().Get(context.Background())
	require.NoError(t, err)
	require.Len(t, plans, 1)
	require.Equal(t, "plan1", plans[0].ID)
}

func TestPlanByIdGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/plans/plan1", `{"id":"plan1","title":"Plan 1"}`)
	defer server.Close()

	client := newClient(server)

	plan, err := client.Planner().ById("plan1").Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "plan1", plan.ID)
}

func TestPlanPatch(t *testing.T) {
	server := newTestServer(t, http.MethodPatch, "/planner/plans/plan1", `{"id":"plan1","title":"Updated Plan 1"}`)
	defer server.Close()

	client := newClient(server)
	client.eTagCache = map[string]string{
		server.URL + "/planner/plans/plan1": "W/\"test-etag\"",
	}

	plan, err := client.Planner().ById("plan1").Patch(context.Background(), PatchPlanParams{Title: "Updated Plan 1"})
	require.NoError(t, err)
	require.Equal(t, "Updated Plan 1", plan.Title)
}

func TestTasksGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/plans/plan1/tasks", `{"value":[{"id":"task1","title":"Task 1"}]}`)
	defer server.Close()

	client := newClient(server)

	tasks, err := client.Planner().ById("plan1").Tasks().Get(context.Background())
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	require.Equal(t, "task1", tasks[0].ID)
}

func TestTaskByIdGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/tasks/task1", `{"id":"task1","title":"Task 1"}`)
	defer server.Close()

	client := newClient(server)

	task, err := client.Planner().Tasks().ById("task1").Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "task1", task.ID)
}

func TestTaskPatch(t *testing.T) {
	server := newTestServer(t, http.MethodPatch, "/planner/tasks/task1", `{"id":"task1","title":"Updated Task 1"}`)
	defer server.Close()

	client := newClient(server)
	client.eTagCache = map[string]string{
		server.URL + "/planner/tasks/task1": "W/\"test-etag\"",
	}

	task, err := client.Planner().Tasks().ById("task1").Patch(context.Background(), PatchTaskParams{Title: "Updated Task 1"})
	require.NoError(t, err)
	require.Equal(t, "Updated Task 1", task.Title)
}

func TestTaskPost(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/planner/tasks", `{"id":"task1","title":"New Task 1"}`)
	defer server.Close()

	client := newClient(server)

	task, err := client.Planner().Tasks().Post(context.Background(), PostTaskParams{Title: "New Task 1", PlanID: "plan1"})
	require.NoError(t, err)
	require.Equal(t, "New Task 1", task.Title)
}

func TestBucketsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/plans/plan1/buckets", `{"value":[{"id":"bucket1","name":"Bucket 1"}]}`)
	defer server.Close()

	client := newClient(server)

	buckets, err := client.Planner().ById("plan1").Buckets().Get(context.Background())
	require.NoError(t, err)
	require.Len(t, buckets, 1)
	require.Equal(t, "bucket1", buckets[0].ID)
}

func TestBucketPatch(t *testing.T) {
	server := newTestServer(t, http.MethodPatch, "/planner/buckets/bucket1", `{"id":"bucket1","name":"Updated Bucket 1"}`)
	defer server.Close()

	client := newClient(server)
	client.eTagCache = map[string]string{
		server.URL + "/planner/buckets/bucket1": "W/\"test-etag\"",
	}

	bucket, err := client.Planner().Buckets().ById("bucket1").Patch(context.Background(), PatchBucketParams{Name: "Updated Bucket 1"})
	require.NoError(t, err)
	require.Equal(t, "Updated Bucket 1", bucket.Name)
}

func TestBucketPost(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/planner/buckets", `{"id":"bucket1","name":"New Bucket 1"}`)
	defer server.Close()

	client := newClient(server)

	bucket, err := client.Planner().Buckets().Post(context.Background(), PostBucketParams{Name: "New Bucket 1", PlanID: "plan1"})
	require.NoError(t, err)
	require.Equal(t, "New Bucket 1", bucket.Name)
}

func TestTaskDetailsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/tasks/task1/details", `{"id":"task1","description":"Task 1 Description"}`)
	defer server.Close()

	client := newClient(server)

	details, err := client.Planner().Tasks().ById("task1").Details().Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "task1", details.ID)
	require.Equal(t, "Task 1 Description", details.Description)
}

func TestPlannerGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner", `{"value":[{"id":"plan1","title":"Plan 1"}]}`)
	defer server.Close()

	client := newClient(server)

	plans, err := client.Planner().Get(context.Background())
	require.NoError(t, err)
	require.Len(t, plans, 1)
	require.Equal(t, "plan1", plans[0].ID)
}

func newTestServer(t *testing.T, code, path, data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, code, r.Method)
		require.Equal(t, path, r.URL.Path)

		switch code {
		case http.MethodGet, http.MethodPatch:
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		}
		w.Write([]byte(data))
	}))
}

func newClient(server *httptest.Server) *Client {
	return &Client{
		BaseURL: server.URL,
		c:       server.Client(),
		limiter: rate.NewLimiter(rate.Limit(100), 200),
	}
}
