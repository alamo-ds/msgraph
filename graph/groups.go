package graph

import (
	"context"
)

const groupResource string = "groups"

type GroupsRequestBuilder struct {
	c    *Client
	path string
}

func (c *Client) Groups() *GroupsRequestBuilder {
	return &GroupsRequestBuilder{
		c:    c,
		path: joinPath(c.BaseURL, groupResource),
	}
}

type GetGroupsResponse struct {
	Value []Group `json:"value"`
}

func (r *GroupsRequestBuilder) Get(ctx context.Context) ([]Group, error) {
	var ret GetGroupsResponse

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return nil, err
	}

	return ret.Value, nil
}

type GroupItemRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *GroupsRequestBuilder) ById(id string) *GroupItemRequestBuilder {
	return &GroupItemRequestBuilder{
		Id:   id,
		c:    r.c,
		path: joinPath(r.path, id),
	}
}

func (r *GroupItemRequestBuilder) Get(ctx context.Context) (Group, error) {
	var ret Group

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
