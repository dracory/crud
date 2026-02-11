package crud

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/form"
)

func TestManager_Page_Success(t *testing.T) {
	crud := newTestCrud()
	crud.funcRows = func(r *http.Request) ([]Row, error) {
		return []Row{
			{ID: "1", Data: []string{"1", "Product A"}},
			{ID: "2", Data: []string{"2", "Product B"}},
		}, nil
	}
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Product Manager") {
		t.Fatal("expected 'Product Manager' in page output")
	}
	if !strings.Contains(body, "Product A") {
		t.Fatal("expected 'Product A' in table output")
	}
	if !strings.Contains(body, "Product B") {
		t.Fatal("expected 'Product B' in table output")
	}
}

func TestManager_Page_FuncRowsError(t *testing.T) {
	crud := newTestCrud()
	crud.funcRows = func(r *http.Request) ([]Row, error) {
		return nil, errors.New("database error")
	}
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "error retrieving the data") {
		t.Fatal("expected error alert when funcRows fails")
	}
}

func TestManager_Page_EmptyRows(t *testing.T) {
	crud := newTestCrud()
	crud.funcRows = func(r *http.Request) ([]Row, error) {
		return []Row{}, nil
	}
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Product Manager") {
		t.Fatal("expected 'Product Manager' heading even with empty rows")
	}
	if !strings.Contains(body, "TableEntities") {
		t.Fatal("expected 'TableEntities' table ID in output")
	}
}

func TestManager_Page_WithCreateFields(t *testing.T) {
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
			Value: "default",
		}),
	}
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "New Product") {
		t.Fatal("expected 'New Product' button in page output")
	}
}

func TestManager_Page_ContentTypeHeader(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Fatalf("expected Content-Type 'text/html', got: %s", contentType)
	}
}

func TestManager_Page_ShowsActionButtons(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchReadData = func(r *http.Request, entityID string) ([][2]string, error) {
		return nil, nil
	}
	crud.funcFetchUpdateData = func(r *http.Request, entityID string) (map[string]string, error) {
		return nil, nil
	}
	crud.funcTrash = func(r *http.Request, entityID string) error {
		return nil
	}
	crud.funcRows = func(r *http.Request) ([]Row, error) {
		return []Row{
			{ID: "1", Data: []string{"1", "Product A"}},
		}, nil
	}
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "bi-eye") {
		t.Fatal("expected view button (bi-eye icon) when funcFetchReadData is set")
	}
	if !strings.Contains(body, "bi-pencil-square") {
		t.Fatal("expected edit button (bi-pencil-square icon) when funcFetchUpdateData is set")
	}
	if !strings.Contains(body, "bi-trash") {
		t.Fatal("expected trash button (bi-trash icon) when funcTrash is set")
	}
}

func TestManager_Page_HidesActionButtonsWhenFuncsNil(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchReadData = nil
	crud.funcFetchUpdateData = nil
	crud.funcTrash = nil
	crud.funcRows = func(r *http.Request) ([]Row, error) {
		return []Row{
			{ID: "1", Data: []string{"1", "Product A"}},
		}, nil
	}
	ctrl := crud.newEntityManagerController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	if strings.Contains(body, "bi-eye") {
		t.Fatal("expected no view button when funcFetchReadData is nil")
	}
	if strings.Contains(body, "bi-pencil-square") {
		t.Fatal("expected no edit button when funcFetchUpdateData is nil")
	}
	if strings.Contains(body, "bi-trash") {
		t.Fatal("expected no trash button when funcTrash is nil")
	}
}
