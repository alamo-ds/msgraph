package graph

import "context"

type PostsRequestBuilder struct {
	Id   string
	c    *Client
	path string
}

func (r *ThreadsRequestBuilder) ById(id string) *PostsRequestBuilder {
	return &PostsRequestBuilder{
		Id:   id,
		c:    r.c,
		path: joinPath(r.path, id, "posts"),
	}
}

type GetPostsResponse struct {
	Value []Post `json:"value"`
}

func (r *PostsRequestBuilder) Get(ctx context.Context) ([]Post, error) {
	var ret GetPostsResponse

	if err := get(ctx, r.c, r.path, &ret); err != nil {
		return nil, err
	}

	return ret.Value, nil
}
