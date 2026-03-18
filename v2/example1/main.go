// example1 demonstrates the v2 CRUD package with an in-memory book library.
// Run with: go run main.go
// Then open: http://localhost:8080/books
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	crud "github.com/dracory/crud/v2"
)

// Book is our example entity.
type Book struct {
	ID     string
	Title  string
	Author string
	Genre  string
	Status string
	Notes  string
}

// --- In-memory store --------------------------------------------------------

var (
	mu      sync.Mutex
	counter int
	books   = map[string]*Book{
		"1": {ID: "1", Title: "The Go Programming Language", Author: "Donovan & Kernighan", Genre: "technology", Status: "available", Notes: "Classic Go reference."},
		"2": {ID: "2", Title: "Clean Code", Author: "Robert C. Martin", Genre: "technology", Status: "borrowed", Notes: ""},
		"3": {ID: "3", Title: "Dune", Author: "Frank Herbert", Genre: "sci-fi", Status: "available", Notes: "Epic sci-fi saga."},
	}
)

func nextID() string {
	counter++
	return fmt.Sprintf("%d", 100+counter)
}

// --- Genre & Status options (reused across fields) --------------------------

var genreOptions = []crud.FieldOption{
	{Key: "technology", Value: "Technology"},
	{Key: "sci-fi", Value: "Sci-Fi"},
	{Key: "fantasy", Value: "Fantasy"},
	{Key: "history", Value: "History"},
	{Key: "other", Value: "Other"},
}

var statusOptions = []crud.FieldOption{
	{Key: "available", Value: "Available"},
	{Key: "borrowed", Value: "Borrowed"},
	{Key: "reserved", Value: "Reserved"},
}

// --- CRUD callbacks ---------------------------------------------------------

func funcRows(r *http.Request) ([]crud.Row, error) {
	mu.Lock()
	defer mu.Unlock()

	var rows []crud.Row
	for _, b := range books {
		statusLabel := b.Status
		for _, o := range statusOptions {
			if o.Key == b.Status {
				statusLabel = o.Value
			}
		}
		genreLabel := b.Genre
		for _, o := range genreOptions {
			if o.Key == b.Genre {
				genreLabel = o.Value
			}
		}
		rows = append(rows, crud.Row{
			ID:   b.ID,
			Data: []string{b.Title, b.Author, genreLabel, statusLabel},
		})
	}
	return rows, nil
}

func funcCreate(r *http.Request, data map[string]string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	id := nextID()
	books[id] = &Book{
		ID:     id,
		Title:  data["title"],
		Author: data["author"],
		Genre:  data["genre"],
		Status: data["status"],
		Notes:  data["notes"],
	}
	return id, nil
}

func funcFetchUpdateData(r *http.Request, entityID string) (map[string]string, error) {
	mu.Lock()
	defer mu.Unlock()

	b, ok := books[entityID]
	if !ok {
		return nil, fmt.Errorf("book %s not found", entityID)
	}
	return map[string]string{
		"title":  b.Title,
		"author": b.Author,
		"genre":  b.Genre,
		"status": b.Status,
		"notes":  b.Notes,
	}, nil
}

func funcUpdate(r *http.Request, entityID string, data map[string]string) error {
	mu.Lock()
	defer mu.Unlock()

	b, ok := books[entityID]
	if !ok {
		return fmt.Errorf("book %s not found", entityID)
	}
	b.Title = data["title"]
	b.Author = data["author"]
	b.Genre = data["genre"]
	b.Status = data["status"]
	b.Notes = data["notes"]
	return nil
}

func funcFetchReadData(r *http.Request, entityID string) ([]crud.KeyValue, error) {
	mu.Lock()
	defer mu.Unlock()

	b, ok := books[entityID]
	if !ok {
		return nil, fmt.Errorf("book %s not found", entityID)
	}

	genreLabel := b.Genre
	for _, o := range genreOptions {
		if o.Key == b.Genre {
			genreLabel = o.Value
		}
	}
	statusLabel := b.Status
	for _, o := range statusOptions {
		if o.Key == b.Status {
			statusLabel = o.Value
		}
	}

	return []crud.KeyValue{
		{Key: "ID", Value: b.ID},
		{Key: "Title", Value: b.Title},
		{Key: "Author", Value: b.Author},
		{Key: "Genre", Value: genreLabel},
		{Key: "Status", Value: statusLabel},
		{Key: "Notes", Value: b.Notes},
	}, nil
}

func funcTrash(r *http.Request, entityID string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := books[entityID]; !ok {
		return fmt.Errorf("book %s not found", entityID)
	}
	delete(books, entityID)
	return nil
}

// --- Fields -----------------------------------------------------------------

var createFields = []crud.FieldInterface{
	crud.NewStringField("title", "Title").WithRequired(),
	crud.NewStringField("author", "Author").WithRequired(),
	crud.NewSelectField("genre", "Genre", genreOptions),
	crud.NewSelectField("status", "Status", statusOptions),
	crud.NewTextAreaField("notes", "Notes").WithHelp("Optional notes about this book."),
}

var updateFields = []crud.FieldInterface{
	crud.NewStringField("title", "Title").WithRequired(),
	crud.NewStringField("author", "Author").WithRequired(),
	crud.NewSelectField("genre", "Genre", genreOptions),
	crud.NewSelectField("status", "Status", statusOptions),
	crud.NewTextAreaField("notes", "Notes").WithHelp("Optional notes about this book."),
}

// --- Main -------------------------------------------------------------------

func main() {
	c, err := crud.New(crud.Config{
		Endpoint:           "/books",
		HomeURL:            "/",
		EntityNameSingular: "Book",
		EntityNamePlural:   "Books",
		ColumnNames:        []string{"Title", "Author", "Genre", "Status"},

		CreateFields:        createFields,
		UpdateFields:        updateFields,
		FuncRows:            funcRows,
		FuncCreate:          funcCreate,
		FuncFetchUpdateData: funcFetchUpdateData,
		FuncUpdate:          funcUpdate,
		FuncFetchReadData:   funcFetchReadData,
		FuncTrash:           funcTrash,
	})
	if err != nil {
		log.Fatal("crud.New:", err)
	}

	http.HandleFunc("/books", c.Handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/books", http.StatusFound)
	})

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
