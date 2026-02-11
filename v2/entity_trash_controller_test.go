package crud

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTrash_MethodNotAllowed_GET(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityTrashController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-trash-ajax&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.pageEntityTrashAjax(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Method not allowed") {
		t.Fatalf("expected 'Method not allowed' message, got: %v", resp["message"])
	}
}

func TestTrash_MissingEntityID(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityTrashController()

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-trash-ajax", nil)
	w := httptest.NewRecorder()

	ctrl.pageEntityTrashAjax(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Entity ID is required") {
		t.Fatalf("expected 'Entity ID is required' message, got: %v", resp["message"])
	}
}

func TestTrash_NilFuncTrash(t *testing.T) {
	crud := newTestCrud()
	crud.funcTrash = nil
	ctrl := crud.newEntityTrashController()

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-trash-ajax&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.pageEntityTrashAjax(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Trash functionality is not configured") {
		t.Fatalf("expected 'Trash functionality is not configured' message, got: %v", resp["message"])
	}
}

func TestTrash_FuncTrashReturnsError(t *testing.T) {
	crud := newTestCrud()
	crud.funcTrash = func(r *http.Request, entityID string) error {
		return errors.New("database error")
	}
	ctrl := crud.newEntityTrashController()

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-trash-ajax&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.pageEntityTrashAjax(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "database error") {
		t.Fatalf("expected 'database error' in message, got: %v", resp["message"])
	}
}

func TestTrash_Success(t *testing.T) {
	var trashedID string
	crud := newTestCrud()
	crud.funcTrash = func(r *http.Request, entityID string) error {
		trashedID = entityID
		return nil
	}
	ctrl := crud.newEntityTrashController()

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-trash-ajax&entity_id=abc-123", nil)
	w := httptest.NewRecorder()

	ctrl.pageEntityTrashAjax(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "success" {
		t.Fatalf("expected status 'success', got: %v", resp["status"])
	}
	if trashedID != "abc-123" {
		t.Fatalf("expected trashedID 'abc-123', got: %s", trashedID)
	}
}
