// example1 demonstrates the v2 CRUD package using a realistic book library.
// The Book entity uses real-world fields plus extras to showcase every field type.
// Run with: go run main.go
// Then open: http://localhost:8080/books
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	crud "github.com/dracory/crud/v2"
)

// Book is the core entity.
type Book struct {
	ID     string
	Title  string
	Author string
	Genre  string
	Status string
	Notes  string

	// Bibliographic
	ISBN        string
	Website     string
	Phone       string
	Price       string
	PublishedAt string
	UpdatedAt   string
	CoverURL    string
	CoverColor  string
	Description string
	Featured    string
	Condition   string
	Tags        string // repeater stored as JSON

	// Inventory — Element Plus components
	Rating      string
	Discount    string
	InStock     string
	Stock       string
	Supplier    string
	SpineColor  string
	RestockDate string
	RestockTime string
	RestockAt   string
}

// --- In-memory store --------------------------------------------------------

var (
	mu      sync.Mutex
	counter int
	books   = map[string]*Book{
		"1": {
			ID: "1", Title: "The Go Programming Language", Author: "Donovan & Kernighan",
			Genre: "technology", Status: "available", Notes: "Classic Go reference.",
			ISBN: "978-0134190440", Website: "https://gopl.io", Price: "49",
			PublishedAt: "2015-10-26", CoverColor: "#4a90d9", Featured: "1",
			Condition: "new", Rating: "5", Discount: "10", InStock: "true",
			Stock: "12", Supplier: "supplier1",
		},
		"2": {
			ID: "2", Title: "Clean Code", Author: "Robert C. Martin",
			Genre: "technology", Status: "borrowed", Notes: "Timeless software craftsmanship.",
			ISBN: "978-0132350884", Price: "35", CoverColor: "#e74c3c",
			Condition: "good", Rating: "4", Discount: "0", InStock: "false",
			Stock: "0", Supplier: "supplier2",
		},
		"3": {
			ID: "3", Title: "Dune", Author: "Frank Herbert",
			Genre: "sci-fi", Status: "available", Notes: "Epic sci-fi saga.",
			ISBN: "978-0441013593", Price: "18", CoverColor: "#f39c12",
			Condition: "good", Rating: "5", Discount: "15", InStock: "true",
			Stock: "5", Supplier: "supplier1",
		},
	}
)

func nextID() string {
	counter++
	return fmt.Sprintf("%d", 100+counter)
}

// --- Options ----------------------------------------------------------------

var genreOptions = []crud.FieldOption{
	{Key: "technology", Value: "Technology"},
	{Key: "sci-fi", Value: "Sci-Fi"},
	{Key: "fantasy", Value: "Fantasy"},
	{Key: "history", Value: "History"},
	{Key: "biography", Value: "Biography"},
	{Key: "other", Value: "Other"},
}

var statusOptions = []crud.FieldOption{
	{Key: "available", Value: "Available"},
	{Key: "borrowed", Value: "Borrowed"},
	{Key: "reserved", Value: "Reserved"},
}

var conditionOptions = []crud.FieldOption{
	{Key: "new", Value: "New"},
	{Key: "good", Value: "Good"},
	{Key: "fair", Value: "Fair"},
	{Key: "poor", Value: "Poor"},
}

var supplierOptions = []crud.FieldOption{
	{Key: "supplier1", Value: "Penguin Books"},
	{Key: "supplier2", Value: "O'Reilly Media"},
	{Key: "supplier3", Value: "Tor Books"},
}

// --- Helpers ----------------------------------------------------------------

func labelOf(opts []crud.FieldOption, key string) string {
	for _, o := range opts {
		if o.Key == key {
			return o.Value
		}
	}
	return key
}

func bookFromData(id string, d map[string]string) *Book {
	return &Book{
		ID: id, Title: d["title"], Author: d["author"],
		Genre: d["genre"], Status: d["status"], Notes: d["notes"],
		ISBN: d["isbn"], Website: d["website"], Phone: d["phone"],
		Price: d["price"], PublishedAt: d["published_at"], UpdatedAt: d["updated_at"],
		CoverURL: d["cover_url"], CoverColor: d["cover_color"],
		Description: d["description"], Featured: d["featured"],
		Condition: d["condition"], Tags: d["tags"],
		Rating: d["rating"], Discount: d["discount"], InStock: d["in_stock"],
		Stock: d["stock"], Supplier: d["supplier"], SpineColor: d["spine_color"],
		RestockDate: d["restock_date"], RestockTime: d["restock_time"], RestockAt: d["restock_at"],
	}
}

func bookToData(b *Book) map[string]string {
	return map[string]string{
		"title": b.Title, "author": b.Author, "genre": b.Genre,
		"status": b.Status, "notes": b.Notes, "isbn": b.ISBN,
		"website": b.Website, "phone": b.Phone, "price": b.Price,
		"published_at": b.PublishedAt, "updated_at": b.UpdatedAt,
		"cover_url": b.CoverURL, "cover_color": b.CoverColor,
		"description": b.Description, "featured": b.Featured,
		"condition": b.Condition, "tags": b.Tags,
		"rating": b.Rating, "discount": b.Discount, "in_stock": b.InStock,
		"stock": b.Stock, "supplier": b.Supplier, "spine_color": b.SpineColor,
		"restock_date": b.RestockDate, "restock_time": b.RestockTime, "restock_at": b.RestockAt,
	}
}

func tagsProcess(raw string) []map[string]string {
	var rows []map[string]string
	json.Unmarshal([]byte(raw), &rows)
	return rows
}

// --- CRUD callbacks ---------------------------------------------------------

func funcRows(r *http.Request) ([]crud.Row, error) {
	mu.Lock()
	defer mu.Unlock()
	var rows []crud.Row
	for _, b := range books {
		rows = append(rows, crud.Row{
			ID:   b.ID,
			Data: []string{b.Title, b.Author, labelOf(genreOptions, b.Genre), labelOf(statusOptions, b.Status)},
		})
	}
	return rows, nil
}

func funcCreate(r *http.Request, data map[string]string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	id := nextID()
	books[id] = bookFromData(id, data)
	return id, nil
}

func funcFetchUpdateData(r *http.Request, entityID string) (map[string]string, error) {
	mu.Lock()
	defer mu.Unlock()
	b, ok := books[entityID]
	if !ok {
		return nil, fmt.Errorf("book %s not found", entityID)
	}
	return bookToData(b), nil
}

func funcUpdate(r *http.Request, entityID string, data map[string]string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := books[entityID]; !ok {
		return fmt.Errorf("book %s not found", entityID)
	}
	books[entityID] = bookFromData(entityID, data)
	return nil
}

func funcFetchReadData(r *http.Request, entityID string) ([]crud.KeyValue, error) {
	mu.Lock()
	defer mu.Unlock()
	b, ok := books[entityID]
	if !ok {
		return nil, fmt.Errorf("book %s not found", entityID)
	}
	return []crud.KeyValue{
		{Key: "ID", Value: b.ID},
		{Key: "Title", Value: b.Title},
		{Key: "Author", Value: b.Author},
		{Key: "Genre", Value: labelOf(genreOptions, b.Genre)},
		{Key: "Status", Value: labelOf(statusOptions, b.Status)},
		{Key: "Condition", Value: labelOf(conditionOptions, b.Condition)},
		{Key: "ISBN", Value: b.ISBN},
		{Key: "Website", Value: b.Website},
		{Key: "Phone", Value: b.Phone},
		{Key: "Price ($)", Value: b.Price},
		{Key: "Published", Value: b.PublishedAt},
		{Key: "Cover Color", Value: b.CoverColor},
		{Key: "Featured", Value: b.Featured},
		{Key: "Rating", Value: b.Rating},
		{Key: "Discount (%)", Value: b.Discount},
		{Key: "In Stock", Value: b.InStock},
		{Key: "Stock Qty", Value: b.Stock},
		{Key: "Supplier", Value: labelOf(supplierOptions, b.Supplier)},
		{Key: "Notes", Value: b.Notes},
		{Key: "Tags (JSON)", Value: b.Tags},
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
	crud.NewRawField(`<h6 class="mt-2 text-muted">Core Details</h6>`),
	crud.NewStringField("title", "Title").WithRequired(),
	crud.NewStringField("author", "Author").WithRequired(),
	crud.NewSelectField("genre", "Genre", genreOptions),
	crud.NewSelectField("status", "Status", statusOptions),
	crud.NewTextAreaField("notes", "Notes").WithHelp("Optional notes about this book."),

	crud.NewRawField(`<h6 class="mt-3 text-muted">Bibliographic</h6>`),
	crud.NewStringField("isbn", "ISBN").WithHelp("e.g. 978-0134190440"),
	crud.NewURLField("website", "Website"),
	crud.NewTelField("phone", "Publisher Phone"),
	crud.NewNumberField("price", "Price ($)"),
	crud.NewDateField("published_at", "Published Date"),
	crud.NewRadioField("condition", "Condition", conditionOptions),
	crud.NewCheckboxField("featured", "Featured"),
	crud.NewColorField("cover_color", "Cover Color"),
	crud.NewImageField("cover_url", "Cover Image URL"),
}

var updateFields = []crud.FieldInterface{
	crud.NewRawField(`<h6 class="mt-2 text-muted">Core Details</h6>`),
	crud.NewStringField("title", "Title").WithRequired(),
	crud.NewStringField("author", "Author").WithRequired(),
	crud.NewSelectField("genre", "Genre", genreOptions),
	crud.NewSelectField("status", "Status", statusOptions),
	crud.NewTextAreaField("notes", "Notes").WithHelp("Optional notes about this book."),

	crud.NewRawField(`<h6 class="mt-3 text-muted">Bibliographic</h6>`),
	crud.NewStringField("isbn", "ISBN").WithHelp("e.g. 978-0134190440"),
	crud.NewURLField("website", "Website"),
	crud.NewTelField("phone", "Publisher Phone"),
	crud.NewNumberField("price", "Price ($)"),
	crud.NewDateField("published_at", "Published Date"),
	crud.NewDateTimeField("updated_at", "Last Updated"),
	crud.NewRadioField("condition", "Condition", conditionOptions),
	crud.NewCheckboxField("featured", "Featured"),
	crud.NewColorField("cover_color", "Cover Color"),
	crud.NewImageField("cover_url", "Cover Image URL"),

	crud.NewRawField(`<h6 class="mt-3 text-muted">Description</h6>`),
	crud.NewHtmlAreaField("description", "Description").WithHelp("Rich text — Trumbowyg editor"),

	crud.NewRawField(`<h6 class="mt-3 text-muted">Tags</h6>`),
	crud.NewRepeater(crud.RepeaterOptions{
		Name:  "tags",
		Label: "Tags",
		Help:  "Add searchable tags for this book",
		Fields: []crud.FieldInterface{
			crud.NewStringField("tag", "Tag"),
		},
		FuncValuesProcess: tagsProcess,
	}),

	crud.NewRawField(`<h6 class="mt-3 text-muted">Inventory — Element Plus</h6>`),
	crud.NewRateElField("rating", "Rating"),
	crud.NewSliderElField("discount", "Discount (%)"),
	crud.NewSwitchElField("in_stock", "In Stock"),
	crud.NewInputNumberElField("stock", "Stock Qty"),
	crud.NewSelectElField("supplier", "Supplier", supplierOptions),
	crud.NewColorElField("spine_color", "Spine Color"),
	crud.NewDateElField("restock_date", "Restock Date"),
	crud.NewTimeElField("restock_time", "Restock Time"),
	crud.NewDateTimeElField("restock_at", "Restock At"),
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
