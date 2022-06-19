package server

import (
	"context"
	"strconv"
	"time"

	api "github.com/atato/api/proto"
)

func (c *cache) Set(ctx context.Context, item *api.Item) (*api.Item, error) {
	var expiration int64
	duration, _ := time.ParseDuration(item.Expiration)
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[item.Key] = Item{
		Object:     item.Value,
		Expiration: expiration,
	}

	return item, nil
}

func (c *cache) Dump(ctx context.Context, args *api.GetKey) (*api.Item, error) {
	key := args.Key
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.items[key]
	// No key found
	if !ok {
		return nil, ErrNoKey
	}

	// This means key has some expiration
	if value.(Item).Expiration > 0 {
		if time.Now().UnixNano() > value.(Item).Expiration {
			return nil, ErrKeyExpired
		}
	}

	return &api.Item{
		Key:        key,
		Value:      value.(Item).Object.(string),
		Expiration: time.Unix(0, value.(Item).Expiration).String(),
	}, nil
}

func (c *cache) Incr(ctx context.Context, args *api.GetKey) (*api.Success, error) {
	key := args.Key
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.items[key]
	// No key found
	if !ok {
		c.items[key] = Item{
			Object: strconv.Itoa(0),
			// Expiration: int64(time.Now().Add(120).Second()),
		}
		return &api.Success{
			Success: true,
		}, nil
	}

	valueToInt, err := strconv.Atoi(value.(Item).Object.(string))
	if err != nil {
		return &api.Success{
			Success: true,
		}, ErrValueIsNotInteger
	}

	valueToInt = valueToInt + 1

	c.items[key] = Item{
		Object:     strconv.Itoa(valueToInt),
		Expiration: value.(Item).Expiration,
	}

	return &api.Success{
		Success: true,
	}, nil
}
