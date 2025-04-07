package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rBurgett/scmsh/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemDriver_generateKey(t *testing.T) {
	id, err := ulid.Parse("01JR7AVVBQXEV6Q415TT05GCTB")
	require.NoError(t, err)

	md := &MemDriver{}
	output := md.generateKey("test", id)

	assert.Equal(t, "test:01JR7AVVBQXEV6Q415TT05GCTB", output)
}

func TestMemDriver_FindAll(t *testing.T) {
	ctx := context.Background()

	items := map[string]string{
		// test namespace
		"test:01JR7AYVRSXV9AYVXV83PCCH1C": `{"some":"thing1"}`,
		"test:01JR7B1X14X7AKGFH85F88ZNVY": `{"some":"thing2"}`,
		"test:01JR7B2DEMEY1SYJ1QKM9S9BMF": `{"some":"thing3"}`,
		// test1 namespace
		"test1:01JR7B2Y4P0SCVKJJWSE5W8C2J": `{"another":"thing"}`,
	}
	md := &MemDriver{
		items: items,
	}

	expected := []string{
		`{"some":"thing1"}`,
		`{"some":"thing2"}`,
		`{"some":"thing3"}`,
	}

	output, err := md.FindAll(ctx, "test")
	require.NoError(t, err)

	assert.Len(t, output, 3)
	assert.Equal(t, expected, output)
}

func TestMemDriver_FindOne(t *testing.T) {
	ctx := context.Background()

	items := map[string]string{
		// test namespace
		"test:01JR7AYVRSXV9AYVXV83PCCH1C": `{"some":"thing1"}`,
		"test:01JR7B1X14X7AKGFH85F88ZNVY": `{"some":"thing2"}`,
		"test:01JR7B2DEMEY1SYJ1QKM9S9BMF": `{"some":"thing3"}`,
		// test1 namespace
		"test1:01JR7B2Y4P0SCVKJJWSE5W8C2J": `{"another":"thing"}`,
	}
	md := &MemDriver{
		items: items,
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
			output, err := md.FindOne(ctx, "test", tt.id)
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

func TestMemDriver_UpsertOne(t *testing.T) {
	ctx := context.Background()

	id2, err := ulid.Parse("01JR7B1X14X7AKGFH85F88ZNVY")
	require.NoError(t, err)

	items := map[string]string{
		// test namespace
		"test:01JR7AYVRSXV9AYVXV83PCCH1C": `{"some":"thing1"}`,
		"test:01JR7B1X14X7AKGFH85F88ZNVY": `{"some":"thing2"}`,
		"test:01JR7B2DEMEY1SYJ1QKM9S9BMF": `{"some":"thing3"}`,
		// test1 namespace
		"test1:01JR7B2Y4P0SCVKJJWSE5W8C2J": `{"another":"thing"}`,
	}
	md := &MemDriver{
		items: items,
	}

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
			err := md.UpsertOne(ctx, "test", tt.id, tt.value)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)

			data, ok := md.items[fmt.Sprintf("test:%s", tt.id)]
			require.True(t, ok)
			assert.Equal(t, tt.value, data)
		})
	}
}

func TestMemDriver_DeleteOne(t *testing.T) {
	ctx := context.Background()

	items := map[string]string{
		// test namespace
		"test:01JR7AYVRSXV9AYVXV83PCCH1C": `{"some":"thing1"}`,
		"test:01JR7B1X14X7AKGFH85F88ZNVY": `{"some":"thing2"}`,
		"test:01JR7B2DEMEY1SYJ1QKM9S9BMF": `{"some":"thing3"}`,
		// test1 namespace
		"test1:01JR7B2Y4P0SCVKJJWSE5W8C2J": `{"another":"thing"}`,
	}
	md := &MemDriver{
		items: items,
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
			err := md.DeleteOne(ctx, "test", tt.id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)

			_, ok := md.items[fmt.Sprintf("test:%s", tt.id)]
			assert.False(t, ok)
		})
	}
}

func TestNewMemDriver(t *testing.T) {
	output := NewMemDriver()

	require.NotNil(t, output)
}
