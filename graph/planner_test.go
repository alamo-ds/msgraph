package graph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestPlannerPlansGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/groups/group1/planner/plans", `{"value":[{"id":"plan1","title":"Plan 1"}]}`)
	defer server.Close()

	client := newClient(server)

	plans, err := client.Groups().ById("group1").Plans().Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(plans) != 1 {
		t.Fatalf("Expected 1 plan, got %d", len(plans))
	}
	if plans[0].ID != "plan1" {
		t.Errorf("Expected plan ID plan1, got %s", plans[0].ID)
	}
}

func TestPlanByIdGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/plans/plan1", `{"id":"plan1","title":"Plan 1"}`)
	defer server.Close()

	client := newClient(server)

	plan, err := client.Planner().ById("plan1").Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plan.ID != "plan1" {
		t.Errorf("Expected plan ID plan1, got %s", plan.ID)
	}
}

func TestPlanPatch(t *testing.T) {
	server := newTestServer(t, http.MethodPatch, "/planner/plans/plan1", `{"id":"plan1","title":"Updated Plan 1"}`)
	defer server.Close()

	client := newClient(server)
	client.eTagCache = map[string]string{
		server.URL + "/planner/plans/plan1": "W/\"test-etag\"",
	}

	plan, err := client.Planner().ById("plan1").Patch(context.Background(), PatchPlanParams{Title: "Updated Plan 1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plan.Title != "Updated Plan 1" {
		t.Errorf("Expected plan title Updated Plan 1, got %s", plan.Title)
	}
}

func TestTasksGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/plans/plan1/tasks", `{"value":[{"id":"task1","title":"Task 1"}]}`)
	defer server.Close()

	client := newClient(server)

	tasks, err := client.Planner().ById("plan1").Tasks().Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "task1" {
		t.Errorf("Expected task ID task1, got %s", tasks[0].ID)
	}
}

func TestTaskByIdGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/tasks/task1", `{"id":"task1","title":"Task 1"}`)
	defer server.Close()

	client := newClient(server)

	task, err := client.Planner().Tasks().ById("task1").Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task.ID != "task1" {
		t.Errorf("Expected task ID task1, got %s", task.ID)
	}
}

func TestTaskPatch(t *testing.T) {
	server := newTestServer(t, http.MethodPatch, "/planner/tasks/task1", `{"id":"task1","title":"Updated Task 1"}`)
	defer server.Close()

	client := newClient(server)
	client.eTagCache = map[string]string{
		server.URL + "/planner/tasks/task1": "W/\"test-etag\"",
	}

	task, err := client.Planner().Tasks().ById("task1").Patch(context.Background(), PatchTaskParams{Title: "Updated Task 1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task.Title != "Updated Task 1" {
		t.Errorf("Expected task title Updated Task 1, got %s", task.Title)
	}
}

func TestTaskPost(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/planner/tasks", `{"id":"task1","title":"New Task 1"}`)
	defer server.Close()

	client := newClient(server)

	task, err := client.Planner().Tasks().Post(context.Background(), PostTaskParams{Title: "New Task 1", PlanID: "plan1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task.Title != "New Task 1" {
		t.Errorf("Expected task title New Task 1, got %s", task.Title)
	}
}

func TestBucketsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/plans/plan1/buckets", `{"value":[{"id":"bucket1","name":"Bucket 1"}]}`)
	defer server.Close()

	client := newClient(server)

	buckets, err := client.Planner().ById("plan1").Buckets().Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(buckets) != 1 {
		t.Fatalf("Expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].ID != "bucket1" {
		t.Errorf("Expected bucket ID bucket1, got %s", buckets[0].ID)
	}
}

func TestBucketPatch(t *testing.T) {
	server := newTestServer(t, http.MethodPatch, "/planner/buckets/bucket1", `{"id":"bucket1","name":"Updated Bucket 1"}`)
	defer server.Close()

	client := newClient(server)
	client.eTagCache = map[string]string{
		server.URL + "/planner/buckets/bucket1": "W/\"test-etag\"",
	}

	bucket, err := client.Planner().Buckets().ById("bucket1").Patch(context.Background(), PatchBucketParams{Name: "Updated Bucket 1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bucket.Name != "Updated Bucket 1" {
		t.Errorf("Expected bucket name Updated Bucket 1, got %s", bucket.Name)
	}
}

func TestBucketPost(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/planner/buckets", `{"id":"bucket1","name":"New Bucket 1"}`)
	defer server.Close()

	client := newClient(server)

	bucket, err := client.Planner().Buckets().Post(context.Background(), PostBucketParams{Name: "New Bucket 1", PlanID: "plan1"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bucket.Name != "New Bucket 1" {
		t.Errorf("Expected bucket name New Bucket 1, got %s", bucket.Name)
	}
}

func TestTaskDetailsGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner/tasks/task1/details", `{"id":"task1","description":"Task 1 Description"}`)
	defer server.Close()

	client := newClient(server)

	details, err := client.Planner().Tasks().ById("task1").Details().Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if details.ID != "task1" {
		t.Errorf("Expected task details ID task1, got %s", details.ID)
	}
	if details.Description != "Task 1 Description" {
		t.Errorf("Expected task description Task 1 Description, got %s", details.Description)
	}
}

func TestPlannerGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/planner", `{"value":[{"id":"plan1","title":"Plan 1"}]}`)
	defer server.Close()

	client := newClient(server)

	plans, err := client.Planner().Get(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(plans) != 1 {
		t.Fatalf("Expected 1 plan, got %d", len(plans))
	}
	if plans[0].ID != "plan1" {
		t.Errorf("Expected plan ID plan1, got %s", plans[0].ID)
	}
}

func newTestServer(t *testing.T, code, path, data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != code {
			t.Errorf("Expected %s, got %s", code, r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("Expected %s, got %s", path, r.URL.Path)
		}

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
