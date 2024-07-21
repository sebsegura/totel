package service

import "context"

type Client interface {
	Create(ctx context.Context, in *Request) error
}

type inmem struct {
	db map[string]*Request
}

func NewInMemClient() Client {
	return &inmem{
		db: make(map[string]*Request),
	}
}

func (c *inmem) Create(ctx context.Context, in *Request) error {
	c.db[in.ID] = in
	return nil
}
