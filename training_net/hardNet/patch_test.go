package hardnet

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBookPatch_Valid(t *testing.T) {
	m := map[string]any{
		"title":  "New Title",
		"author": "New Author",
		"year":   json.Number("2000"),
		"status": "borrowed",
		"tags":   []any{"a", "b"},
		"publisher": map[string]any{
			"name":    "Penguin",
			"country": "UK",
		},
	}

	p, er := parseBookPatch(m)

	require.NoError(t, er)

	assert.Equal(t, "New Title", *p.Title)
	assert.Equal(t, "New Author", *p.Author)
	assert.Equal(t, 2000, *p.Year)
	assert.Equal(t, StatusBorrowed, *p.Status)
	assert.Equal(t, []string{"a", "b"}, p.Tags)
	assert.Equal(t, "Penguin", p.Publisher.Name)
	assert.Equal(t, "UK", p.Publisher.Country)
}

func TestParseBookPatch_Empty(t *testing.T) {
	p, err := parseBookPatch(map[string]any{})
	require.NoError(t, err)
	assert.Nil(t, p.Title)
	assert.Nil(t, p.Year)
	assert.Nil(t, p.Tags)
}

func TestParseBookPatch_Errors(t *testing.T) {
	cases := []struct {
		name string
		m    map[string]any
		want string
	}{
		{"title not string", map[string]any{"title": 123}, "title must be string"},
		{"title empty", map[string]any{"title": ""}, "title cannot be empty"},
		{"year not number", map[string]any{"year": "2000"}, "year must be number"},
		{"status unknown", map[string]any{"status": "flying"}, "unknown status"},
		{"tags not array", map[string]any{"tags": "a"}, "tags must be array"},
		{"tags not string", map[string]any{"tags": []any{1, 2}}, "tags must be string array"},
		{"publisher not object", map[string]any{"publisher": "x"}, "publisher must be object"},
		{"publisher unknown field", map[string]any{"publisher": map[string]any{"x": "y"}}, "unknown publisher field"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseBookPatch(tc.m)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}
