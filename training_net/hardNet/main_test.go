package hardnet

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBook_Partial(t *testing.T) {
	srv := newTestServer()
	srv.store.books[1] = Book{ID: 1, Title: "Old", Author: "Old", Year: 1900, Status: StatusAvailable}

	body := `{"title":"New","tags":["a","b"]}`
	req := httptest.NewRequest("PUT", "/books/1", bytes.NewBufferString(body))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	srv.handleUpdateBook(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var book Book
	json.NewDecoder(rec.Body).Decode(&book)
	assert.Equal(t, "New", book.Title)
	assert.Equal(t, "Old", book.Author) // не тронули
	assert.Equal(t, 1900, book.Year)    // не тронули
	assert.Equal(t, []string{"a", "b"}, book.Tags)
}

func BenchmarkParseBookPatch(b *testing.B) {
	m := map[string]any{
		"title":  "Title",
		"author": "Author",
		"year":   json.Number("2000"),
		"status": "available",
	}

	for b.Loop() {
		_, _ = parseBookPatch(m)
	}
}
