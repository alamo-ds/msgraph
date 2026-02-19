package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	plannerResource string = "planner"
	plansResource   string = "plans"
)

type PlannerRequestBuilder struct {
	c    *Client
	path string
}

func (c *Client) Planner() *PlannerRequestBuilder {
	return &PlannerRequestBuilder{
		c:    c,
		path: joinPath(c.BaseURL, plannerResource),
	}
}

func (r *GroupItemRequestBuilder) Plans() *PlannerRequestBuilder {
	return &PlannerRequestBuilder{
		c:    r.c,
		path: joinPath(r.path, plannerResource, plansResource),
	}
}

type PlanRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *PlannerRequestBuilder) ById(id string) *PlanRequestBuilder {
	return &PlanRequestBuilder{
		Id:   id,
		c:    r.c,
		path: joinPath(r.path, "plans", id),
	}
}

type GetPlansResponse struct {
	Count int    `json:"@odata.count"`
	Value []Plan `json:"value"`
}

func (r *PlannerRequestBuilder) Get(ctx context.Context) ([]Plan, error) {
	var ret GetPlansResponse

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return nil, err
	}

	return ret.Value, nil
}

type TasksRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *PlannerRequestBuilder) Tasks() *TasksRequestBuilder {
	return &TasksRequestBuilder{
		c:    r.c,
		path: joinPath(r.path, "tasks"),
	}
}

func (r *PlanRequestBuilder) Tasks() *TasksRequestBuilder {
	return &TasksRequestBuilder{
		Id:   r.Id,
		c:    r.c,
		path: joinPath(r.path, "tasks"),
	}
}

type TaskRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *TasksRequestBuilder) ById(id string) *TaskRequestBuilder {
	return &TaskRequestBuilder{
		Id:   id,
		c:    r.c,
		path: joinPath(r.path, id),
	}
}

func (r *TaskRequestBuilder) Get(ctx context.Context) (Task, error) {
	var ret Task
	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return ret, err
	}
	if ret.OdataEtag != "" {
		r.c.putETag(r.path, ret.OdataEtag)
	}

	return ret, nil
}

type PatchTaskParams struct {
	AppliedCategories    map[string]bool       `json:"appliedCategories,omitempty"`
	AssigneePriority     string                `json:"assigneePriority,omitempty"`
	Assignments          map[string]Assignment `json:"assignments,omitempty"`
	BucketID             string                `json:"bucketId,omitempty"`
	ConversationThreadID string                `json:"conversationThreadId,omitempty"`
	DueDateTime          time.Time             `json:"dueDateTime,omitzero"`
	OrderHint            string                `json:"orderHint,omitempty"`
	Priority             int                   `json:"priority,omitempty"`
	PercentComplete      int                   `json:"percentComplete,omitempty"`
	StartDateTime        time.Time             `json:"startDateTime,omitzero"`
	Title                string                `json:"title,omitempty"`
}

func (r *TaskRequestBuilder) Patch(ctx context.Context, params PatchTaskParams) (Task, error) {
	var ret Task

	resp, err := r.c.patch(ctx, r.path, toBody(params))
	if err != nil {
		return ret, makeReqErr(err)
	}

	if err := handlePatchPostResp(resp, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

type PostTaskParams struct {
	PlanID                   string                `json:"planId"`
	BucketID                 string                `json:"bucketId,omitempty"`
	Title                    string                `json:"title"`
	OrderHint                string                `json:"orderHint,omitempty"`
	AssigneePriority         string                `json:"assigneePriority,omitempty"`
	PercentComplete          int                   `json:"percentComplete,omitempty"`
	StartDateTime            time.Time             `json:"startDateTime,omitzero"`
	CreatedDateTime          time.Time             `json:"createdDateTime,omitzero"`
	DueDateTime              time.Time             `json:"dueDateTime,omitzero"`
	HasDescription           bool                  `json:"hasDescription,omitempty"`
	PreviewType              string                `json:"previewType,omitempty"`
	CompletedDateTime        time.Time             `json:"completedDateTime,omitzero"`
	ReferenceCount           int                   `json:"referenceCount,omitempty"`
	ChecklistItemCount       int                   `json:"checklistItemCount,omitempty"`
	ActiveChecklistItemCount int                   `json:"activeChecklistItemCount,omitempty"`
	ConversationThreadID     string                `json:"conversationThreadId,omitempty"`
	Priority                 int                   `json:"priority,omitempty"`
	CreatedBy                IdentitySet           `json:"createdBy,omitzero"`
	CompletedBy              IdentitySet           `json:"completedBy,omitzero"`
	AppliedCategories        map[string]bool       `json:"appliedCategories,omitempty"`
	Assignments              map[string]Assignment `json:"assignments,omitempty"`
}

func (r *TasksRequestBuilder) Post(ctx context.Context, params PostTaskParams) (Task, error) {
	var ret Task

	resp, err := r.c.post(ctx, r.path, toBody(params))
	if err != nil {
		return ret, err
	}

	if err := handlePatchPostResp(resp, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

// NOTE: val must be a pointer to a map or struct
func handlePatchPostResp[T any](resp *http.Response, val T) error {
	defer resp.Body.Close()

	switch resp.StatusCode {
	default:
		return requestErr(resp)
	case http.StatusNoContent:
		log.Printf("request successful but returned no content (%d)\n", resp.StatusCode)
		return nil
	case http.StatusOK, http.StatusCreated:
		if err := json.NewDecoder(resp.Body).Decode(val); err != nil {
			return fmt.Errorf("json.Decode: %v", err)
		}

		return nil
	}
}

type GetTasksResponse struct {
	Count int    `json:"@odata.count"`
	Value []Task `json:"value"`
}

func (r *TasksRequestBuilder) Get(ctx context.Context) ([]Task, error) {
	var ret GetTasksResponse

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return nil, err
	}

	return ret.Value, nil

}

type TaskItemRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *TaskRequestBuilder) Details() *TaskItemRequestBuilder {
	return &TaskItemRequestBuilder{
		Id:   r.Id,
		c:    r.c,
		path: joinPath(r.path, "details"),
	}
}

func (r *TaskItemRequestBuilder) Get(ctx context.Context) (TaskDetails, error) {
	var ret TaskDetails

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

type PatchPlanParams struct {
	Title string `json:"title,omitempty"`
}

func (r *PlanRequestBuilder) Get(ctx context.Context) (Plan, error) {
	var ret Plan

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

func (r *PlanRequestBuilder) Patch(ctx context.Context, params PatchPlanParams) (Plan, error) {
	var ret Plan

	resp, err := r.c.patch(ctx, r.path, toBody(params))
	if err != nil {
		return ret, makeReqErr(err)
	}

	if err := handlePatchPostResp(resp, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

type BucketsRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *PlanRequestBuilder) Buckets() *BucketsRequestBuilder {
	return &BucketsRequestBuilder{
		Id:   r.Id,
		c:    r.c,
		path: joinPath(r.path, "buckets"),
	}
}

type GetBucketsResponse struct {
	Count int      `json:"@odata.count"`
	Value []Bucket `json:"value"`
}

func (r *BucketsRequestBuilder) Get(ctx context.Context) ([]Bucket, error) {
	var ret GetBucketsResponse

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return nil, err
	}

	return ret.Value, nil
}

type BucketItemRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *PlannerRequestBuilder) Buckets() *BucketItemRequestBuilder {
	return &BucketItemRequestBuilder{
		c:    r.c,
		path: joinPath(r.path, "buckets"),
	}
}

func (r *BucketItemRequestBuilder) ById(id string) *BucketItemRequestBuilder {
	r.Id = id
	r.path = joinPath(r.path, id)
	return r
}

type PatchBucketParams struct {
	Name      string `json:"name,omitempty"`
	OrderHint string `json:"orderHint,omitempty"`
}

func (r *BucketItemRequestBuilder) Patch(ctx context.Context, params PatchBucketParams) (Bucket, error) {
	var ret Bucket
	// TODO: add this under r.c.patch to catch all instances
	// where the id may not be set
	if r.Id == "" {
		// something among these lines
		return ret, fmt.Errorf("id for resource type %T not set", ret)
	}

	resp, err := r.c.patch(ctx, r.path, toBody(params))
	if err != nil {
		return ret, makeReqErr(err)
	}

	if err := handlePatchPostResp(resp, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

type PostBucketParams struct {
	Name      string `json:"name"`
	OrderHint string `json:"orderHint,omitempty"`
	PlanID    string `json:"planId"`
}

func (r *BucketItemRequestBuilder) Post(ctx context.Context, params PostBucketParams) (Bucket, error) {
	var ret Bucket

	resp, err := r.c.post(ctx, r.path, toBody(params))
	if err != nil {
		return ret, err
	}

	if err := handlePatchPostResp(resp, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
