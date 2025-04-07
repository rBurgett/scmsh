package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/oklog/ulid/v2"
	"github.com/rBurgett/scmsh/internal/config"
	"github.com/rBurgett/scmsh/internal/constants"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisDriver_generateKey(t *testing.T) {
	id, err := ulid.Parse("01JR7AVVBQXEV6Q415TT05GCTB")
	require.NoError(t, err)

	d := &RedisDriver{}
	output := d.generateKey("test", id)

	assert.Equal(t, "test:01JR7AVVBQXEV6Q415TT05GCTB", output)
}

func TestRedisDriver_FindAll(t *testing.T) {
	ctx := context.Background()

	s := miniredis.RunT(t)
	assert.NotEmpty(t, s.Addr())

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	err := client.Set(ctx, "test:01JR7AYVRSXV9AYVXV83PCCH1C", `{"some":"thing1"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B1X14X7AKGFH85F88ZNVY", `{"some":"thing2"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B2DEMEY1SYJ1QKM9S9BMF", `{"some":"thing3"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test1:01JR7B2Y4P0SCVKJJWSE5W8C2J", `{"another":"thing"}`, 0).Err()
	require.Nil(t, err)

	d := &RedisDriver{
		client: client,
	}

	expected := []string{
		`{"some":"thing1"}`,
		`{"some":"thing2"}`,
		`{"some":"thing3"}`,
	}

	output, err := d.FindAll(ctx, "test")
	require.NoError(t, err)

	assert.Len(t, output, 3)
	assert.Equal(t, expected, output)
}

func TestRedisDriver_FindOne(t *testing.T) {
	ctx := context.Background()

	s := miniredis.RunT(t)
	assert.NotEmpty(t, s.Addr())

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	err := client.Set(ctx, "test:01JR7AYVRSXV9AYVXV83PCCH1C", `{"some":"thing1"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B1X14X7AKGFH85F88ZNVY", `{"some":"thing2"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B2DEMEY1SYJ1QKM9S9BMF", `{"some":"thing3"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test1:01JR7B2Y4P0SCVKJJWSE5W8C2J", `{"another":"thing"}`, 0).Err()
	require.Nil(t, err)

	d := &RedisDriver{
		client: client,
	}

	id2, err := ulid.Parse("01JR7B1X14X7AKGFH85F88ZNVY")
	require.NoError(t, err)

	tests := []struct {
		name          string
		id            ulid.ULID
		expected      string
		expectedError error
	}{
		{
			name:     "good id",
			id:       id2,
			expected: `{"some":"thing2"}`,
		},
		{
			name:          "bad id",
			id:            ulid.Make(),
			expectedError: constants.ErrorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := d.FindOne(ctx, "test", tt.id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestRedisDriver_UpsertOne(t *testing.T) {
	ctx := context.Background()

	s := miniredis.RunT(t)
	assert.NotEmpty(t, s.Addr())

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	err := client.Set(ctx, "test:01JR7AYVRSXV9AYVXV83PCCH1C", `{"some":"thing1"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B1X14X7AKGFH85F88ZNVY", `{"some":"thing2"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B2DEMEY1SYJ1QKM9S9BMF", `{"some":"thing3"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test1:01JR7B2Y4P0SCVKJJWSE5W8C2J", `{"another":"thing"}`, 0).Err()
	require.Nil(t, err)

	d := &RedisDriver{
		client: client,
	}

	id2, err := ulid.Parse("01JR7B1X14X7AKGFH85F88ZNVY")
	require.NoError(t, err)

	tests := []struct {
		name          string
		id            ulid.ULID
		value         string
		expectedError error
	}{
		{
			name:  "new item",
			id:    ulid.Make(),
			value: `{"some":"thing4"}`,
		},
		{
			name:  "existing item",
			id:    id2,
			value: `{"some":"thing5"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.UpsertOne(ctx, "test", tt.id, tt.value)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)

			data, err := s.Get(fmt.Sprintf("test:%s", tt.id))
			require.NoError(t, err)
			assert.Equal(t, tt.value, data)
		})
	}
}

func TestRedisDriver_DeleteOne(t *testing.T) {
	ctx := context.Background()

	s := miniredis.RunT(t)
	assert.NotEmpty(t, s.Addr())

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	err := client.Set(ctx, "test:01JR7AYVRSXV9AYVXV83PCCH1C", `{"some":"thing1"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B1X14X7AKGFH85F88ZNVY", `{"some":"thing2"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test:01JR7B2DEMEY1SYJ1QKM9S9BMF", `{"some":"thing3"}`, 0).Err()
	require.Nil(t, err)

	err = client.Set(ctx, "test1:01JR7B2Y4P0SCVKJJWSE5W8C2J", `{"another":"thing"}`, 0).Err()
	require.Nil(t, err)

	d := &RedisDriver{
		client: client,
	}

	id2, err := ulid.Parse("01JR7B1X14X7AKGFH85F88ZNVY")
	require.NoError(t, err)

	tests := []struct {
		name          string
		id            ulid.ULID
		expectedError error
	}{
		{
			name: "good id",
			id:   id2,
		},
		{
			name: "bad id",
			id:   ulid.Make(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.DeleteOne(ctx, "test", tt.id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)

			_, err = client.Get(ctx, fmt.Sprintf("test:%s", tt.id)).Result()
			require.Error(t, err)
			require.ErrorIs(t, err, redis.Nil)
		})
	}
}

func TestNewRedisDriver(t *testing.T) {
	cfg := config.Config{
		RedisAddress:  "redisaddress:1234",
		RedisPassword: "redispassword",
		RedisDatabase: 1,
	}
	output := NewRedisDriver(cfg)

	require.NotNil(t, output)
	require.NotNil(t, output.client)
}
