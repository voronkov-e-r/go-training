package hardnet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Status string

const (
	StatusAvailable Status = "available"
	StatusBorrowed  Status = "borrowed"
	StatusLost      Status = "lost"
)

func (st *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch Status(str) {
	case StatusAvailable, StatusBorrowed, StatusLost:
		*st = Status(str)
		return nil
	default:
		return fmt.Errorf("unknown status: %q", str)
	}
}

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year"`
	Status    Status    `json:"status,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	Publisher Publisher `json:"publisher"`
}

type Publisher struct {
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
}

type Store struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID atomic.Int64
}

type Server struct {
	store *Store
}

func (s *Server) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value(requestIDKey).(string)
	log.Printf("[%s] GET /books", id)
	s.store.mu.RLock()
	books := make([]Book, 0, len(s.store.books))
	for _, b := range s.store.books {
		books = append(books, b)
	}
	s.store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (s *Server) handleGetBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	s.store.mu.RLock()

	val, ok := s.store.books[id]
	s.store.mu.RUnlock()
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(val)
}

func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&book); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" {
		http.Error(w, "Incorrect input", http.StatusBadRequest)
		return
	}

	if book.Status == "" {
		book.Status = StatusAvailable
	}

	id := int(s.store.nextID.Add(1))
	book.ID = id
	s.store.mu.Lock()
	s.store.books[id] = book
	s.store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/books/%d", id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	p, err := parseBookPatch(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.store.mu.Lock()
	book, ok := s.store.books[id]
	if !ok {
		s.store.mu.Unlock()
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	if p.Title != nil {
		book.Title = *p.Title
	}
	if p.Author != nil {
		book.Author = *p.Author
	}
	if p.Year != nil {
		book.Year = *p.Year
	}
	if p.Status != nil {
		book.Status = *p.Status
	}
	if p.Tags != nil {
		book.Tags = p.Tags
	}
	if p.Publisher != nil {
		book.Publisher = *p.Publisher
	}

	s.store.books[id] = book
	s.store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	s.store.mu.Lock()
	_, ok := s.store.books[id]
	if ok {
		delete(s.store.books, id)
	}
	s.store.mu.Unlock()

	if !ok {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func main() {
	store := &Store{books: make(map[int]Book)}
	srv := &Server{store: store}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", srv.handleGetBooks)
	mux.HandleFunc("GET /books/{id}", srv.handleGetBook)
	mux.HandleFunc("POST /books", srv.handleCreateBook)
	mux.HandleFunc("PUT /books/{id}", srv.handleUpdateBook)
	mux.HandleFunc("DELETE /books/{id}", srv.handleDeleteBook)
	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(5 * time.Second):
			fmt.Fprintln(w, "done")
		case <-r.Context().Done():
			log.Println("client disconnected")
		}
	})

	handle := Chain(mux, Logging, Recovery, RequestID)

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      handle,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server starting on :8080")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if er := httpSrv.Shutdown(ctx); er != nil {
		log.Fatalf("forced shutdown: %v", er)
	}
	log.Println("server stopped")
}
