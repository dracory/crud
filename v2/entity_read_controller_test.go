package crud

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRead_Page_MissingEntityID(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityReadController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-read", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Entity ID is required") {
		t.Fatalf("expected 'Entity ID is required', got: %v", resp["message"])
	}
}

func TestRead_Page_NilFuncFetchReadData(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchReadData = nil
	ctrl := crud.newEntityReadController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-read&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "FuncFetchReadData is required") {
		t.Fatalf("expected 'FuncFetchReadData is required', got: %v", resp["message"])
	}
}

func TestRead_Page_FetchDataError(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchReadData = func(r *http.Request, entityID string) ([][2]string, error) {
		return nil, errors.New("read error")
	}
	ctrl := crud.newEntityReadController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-read&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "error retrieving the data") {
		t.Fatal("expected error alert in page output when fetch fails")
	}
}

func TestRead_Page_Success(t *testing.T) {
	var fetchedID string
	crud := newTestCrud()
	crud.funcFetchReadData = func(r *http.Request, entityID string) ([][2]string, error) {
		fetchedID = entityID
		return [][2]string{
			{"Name", "Test Product"},
			{"Status", "Active"},
		}, nil
	}
	ctrl := crud.newEntityReadController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-read&entity_id=xyz-789", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "View Product") {
		t.Fatal("expected 'View Product' in page output")
	}
	if !strings.Contains(body, "Test Product") {
		t.Fatal("expected 'Test Product' in page output")
	}
	if !strings.Contains(body, "Active") {
		t.Fatal("expected 'Active' in page output")
	}
	if fetchedID != "xyz-789" {
		t.Fatalf("expected fetchedID 'xyz-789', got: %s", fetchedID)
	}
}

func TestRead_Page_RawKeyValueRendering(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchReadData = func(r *http.Request, entityID string) ([][2]string, error) {
		return [][2]string{
			{"{!!Name!!}", "{!!<b>Bold</b>!!}"},
		}, nil
	}
	ctrl := crud.newEntityReadController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-read&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "<b>Bold</b>") {
		t.Fatal("expected raw HTML rendering for {!! !!} wrapped values")
	}
}
