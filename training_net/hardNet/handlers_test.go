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

func newTestServer() *Server {
	return &Server{store: &Store{books: make(map[int]Book)}}
}

func TestCreateBook_Success(t *testing.T) {
	srv := newTestServer()

	body := `{"title":"1984","author":"Orwell","year":1949}`
	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.handleCreateBook(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("Location"))

	var book Book
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&book))
	assert.Equal(t, 1, book.ID)
	assert.Equal(t, "1984", book.Title)
	assert.Equal(t, StatusAvailable, book.Status)
}

func TestCreateBook_InvalidJSON(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString("not json"))
	rec := httptest.NewRecorder()

	srv.handleCreateBook(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBook_MissingFields(t *testing.T) {
	srv := newTestServer()
	body := `{"title":"","author":""}`
	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	srv.handleCreateBook(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetBook_NotFound(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest("GET", "/books/999", nil)
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()

	srv.handleGetBook(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
