package graph

import (
	"context"
)

const usersResource string = "users"

type UsersRequestBuilder struct {
	c            *Client
	path         string
	selectParams []string
}

func (c *Client) Users() *UsersRequestBuilder {
	return &UsersRequestBuilder{
		c:    c,
		path: joinPath(c.BaseURL, usersResource),
	}
}

func (r *UsersRequestBuilder) Select(params ...string) *UsersRequestBuilder {
	r.selectParams = append(r.selectParams, params...)
	return r
}

type GetUsersResponse struct {
	Count int    `json:"@odatta.count"`
	Value []User `json:"value"`
}

func (r *UsersRequestBuilder) Get(ctx context.Context) ([]User, error) {
	var ret GetUsersResponse

	selectParams := userSelectParams(r.selectParams)

	if err := get(ctx, r.c, r.path+"?$select"+selectParams, &ret); err != nil {
		return nil, err
	}

	return ret.Value, nil
}

type UserRequestBuilder struct {
	Id           string
	c            *Client
	path         string
	selectParams []string
}

func (r *UsersRequestBuilder) ById(id string) *UserRequestBuilder {
	return &UserRequestBuilder{
		Id:   id,
		c:    r.c,
		path: joinPath(r.path, id),
	}
}

func (r *UserRequestBuilder) Select(params ...string) *UserRequestBuilder {
	r.selectParams = append(r.selectParams, params...)
	return r
}

func (r *UserRequestBuilder) Get(ctx context.Context) (User, error) {
	var ret User

	selectParams := userSelectParams(r.selectParams)

	// TODO: fix this
	if err := get(ctx, r.c, r.path+"?$select="+selectParams, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
