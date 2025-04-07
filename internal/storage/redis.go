package storage

import (
	"context"
	"fmt"
	"sort"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/rBurgett/scmsh/internal/config"
	"github.com/rBurgett/scmsh/internal/constants"
	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	client *redis.Client
}

func (d *RedisDriver) generateKey(namespace string, id ulid.ULID) string {
	return fmt.Sprintf("%s:%s", namespace, id)
}

func (d *RedisDriver) FindAll(ctx context.Context, namespace string) (res []string, err error) {
	prefix := fmt.Sprintf("%s:*", namespace)
	keys, err := d.client.Keys(ctx, prefix).Result()
	if err != nil {
		return nil, err
	}

	type kv struct {
		k string
		v string
	}
	var items []kv
	for _, k := range keys {
		v, err := d.client.Get(ctx, k).Result()
		if err != nil {
			return nil, err
		}
		items = append(items, kv{k: k, v: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].k < items[j].k
	})

	for _, i := range items {
		res = append(res, i.v)
	}

	return res, nil
}

func (d *RedisDriver) FindOne(ctx context.Context, namespace string, id ulid.ULID) (string, error) {
	key := d.generateKey(namespace, id)
	val, err := d.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", constants.ErrorNotFound
		} else {
			return "", err
		}
	}

	return val, nil
}

func (d *RedisDriver) UpsertOne(ctx context.Context, namespace string, id ulid.ULID, value string) error {
	key := d.generateKey(namespace, id)

	err := d.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (d *RedisDriver) DeleteOne(ctx context.Context, namespace string, id ulid.ULID) error {
	key := d.generateKey(namespace, id)
	err := d.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func NewRedisDriver(cfg config.Config) *RedisDriver {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	return &RedisDriver{
		client: client,
	}
}
