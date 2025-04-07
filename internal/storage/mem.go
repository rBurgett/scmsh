package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/rBurgett/scmsh/internal/constants"
)

type MemDriver struct {
	m     sync.RWMutex
	items map[string]string
}

func (d *MemDriver) generateKey(namespace string, id ulid.ULID) string {
	return fmt.Sprintf("%s:%s", namespace, id)
}

func (d *MemDriver) FindAll(ctx context.Context, namespace string) (res []string, err error) {
	d.m.RLock()
	defer d.m.RUnlock()

	type kv struct {
		k string
		v string
	}
	var items []kv
	for k, v := range d.items {
		if strings.HasPrefix(k, fmt.Sprintf("%s:", namespace)) {
			items = append(items, kv{k, v})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].k < items[j].k
	})

	for _, i := range items {
		res = append(res, i.v)
	}

	return res, nil
}

func (d *MemDriver) FindOne(ctx context.Context, namespace string, id ulid.ULID) (res string, err error) {
	d.m.RLock()
	defer d.m.RUnlock()

	key := d.generateKey(namespace, id)
	res, ok := d.items[key]
	if !ok {
		return "", constants.ErrorNotFound
	}

	return res, nil
}

func (d *MemDriver) UpsertOne(ctx context.Context, namespace string, id ulid.ULID, value string) (err error) {
	d.m.Lock()
	defer d.m.Unlock()

	key := d.generateKey(namespace, id)
	d.items[key] = value

	return nil
}

func (d *MemDriver) DeleteOne(ctx context.Context, namespace string, id ulid.ULID) error {
	d.m.Lock()
	defer d.m.Unlock()

	key := d.generateKey(namespace, id)
	delete(d.items, key)

	return nil
}

func NewMemDriver() *MemDriver {
	return &MemDriver{}
}
