package storage

import (
	"context"
	"encoding/json"

	"github.com/oklog/ulid/v2"
)

type Driver interface {
	FindAll(ctx context.Context, namespace string) (res []string, err error)
	FindOne(ctx context.Context, namespace string, id ulid.ULID) (res string, err error)
	UpsertOne(ctx context.Context, namespace string, id ulid.ULID, value string) error
	DeleteOne(ctx context.Context, namespace string, id ulid.ULID) error
}

type Client[T any] struct {
	driver    Driver
	namespace string
}

func (c *Client[T]) FindAll(ctx context.Context) (res []T, err error) {
	data, err := c.driver.FindAll(ctx, c.namespace)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		var t T
		err = json.Unmarshal([]byte(d), &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (c *Client[T]) FindOne(ctx context.Context, id ulid.ULID) (res T, err error) {
	data, err := c.driver.FindOne(ctx, c.namespace, id)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client[T]) UpsertOne(ctx context.Context, id ulid.ULID, value T) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.driver.UpsertOne(ctx, c.namespace, id, string(encoded))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client[T]) DeleteOne(ctx context.Context, id ulid.ULID) error {
	err := c.driver.DeleteOne(ctx, c.namespace, id)
	if err != nil {
		return err
	}

	return nil
}

func NewClient[T any](driver Driver, namespace string) *Client[T] {
	return &Client[T]{
		driver:    driver,
		namespace: namespace,
	}
}
