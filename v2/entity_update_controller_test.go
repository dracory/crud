package crud

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/form"
)

func TestUpdate_Page_MissingEntityID(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityUpdateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-update", nil)
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

func TestUpdate_Page_NilFuncFetchUpdateData(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchUpdateData = nil
	ctrl := crud.newEntityUpdateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-update&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "Update functionality is not configured") {
		t.Fatalf("expected 'Update functionality is not configured', got: %v", resp["message"])
	}
}

func TestUpdate_Page_FetchDataError(t *testing.T) {
	crud := newTestCrud()
	crud.funcFetchUpdateData = func(entityID string) (map[string]string, error) {
		return nil, errors.New("fetch failed")
	}
	ctrl := crud.newEntityUpdateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-update&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "Fetch data failed") {
		t.Fatalf("expected 'Fetch data failed', got: %v", resp["message"])
	}
}

func TestUpdate_Page_Success(t *testing.T) {
	crud := newTestCrud()
	crud.updateFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
		}),
	}
	crud.funcFetchUpdateData = func(entityID string) (map[string]string, error) {
		return map[string]string{"title": "Existing Product"}, nil
	}
	ctrl := crud.newEntityUpdateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-update&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.page(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Edit Product") {
		t.Fatal("expected 'Edit Product' in page output")
	}
}

func TestUpdate_PageSave_MethodNotAllowed_GET(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityUpdateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-update-ajax&entity_id=123", nil)
	w := httptest.NewRecorder()

	ctrl.pageSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "Method not allowed") {
		t.Fatalf("expected 'Method not allowed', got: %v", resp["message"])
	}
}

func TestUpdate_PageSave_MissingEntityID(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityUpdateController()

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-update-ajax", nil)
	w := httptest.NewRecorder()

	ctrl.pageSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "Entity ID is required") {
		t.Fatalf("expected 'Entity ID is required', got: %v", resp["message"])
	}
}

func TestUpdate_PageSave_RequiredFieldEmpty(t *testing.T) {
	crud := newTestCrud()
	crud.updateFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:     "title",
			Label:    "Title",
			Type:     FORM_FIELD_TYPE_STRING,
			Required: true,
		}),
	}
	crud.funcUpdate = func(entityID string, data map[string]string) error { return nil }
	ctrl := crud.newEntityUpdateController()

	formData := url.Values{}
	formData.Set("entity_id", "123")
	formData.Set("title", "")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-update-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.pageSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "is required field") {
		t.Fatalf("expected 'is required field', got: %v", resp["message"])
	}
}

func TestUpdate_PageSave_FuncUpdateError(t *testing.T) {
	crud := newTestCrud()
	crud.updateFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
		}),
	}
	crud.funcUpdate = func(entityID string, data map[string]string) error {
		return errors.New("update failed")
	}
	ctrl := crud.newEntityUpdateController()

	formData := url.Values{}
	formData.Set("entity_id", "123")
	formData.Set("title", "Updated")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-update-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.pageSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if !strings.Contains(resp["message"].(string), "Save failed") {
		t.Fatalf("expected 'Save failed', got: %v", resp["message"])
	}
}

func TestUpdate_PageSave_Success(t *testing.T) {
	var savedID string
	var savedData map[string]string
	crud := newTestCrud()
	crud.updateFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
		}),
	}
	crud.funcUpdate = func(entityID string, data map[string]string) error {
		savedID = entityID
		savedData = data
		return nil
	}
	ctrl := crud.newEntityUpdateController()

	formData := url.Values{}
	formData.Set("entity_id", "abc-456")
	formData.Set("title", "Updated Product")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-update-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.pageSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "success" {
		t.Fatalf("expected status 'success', got: %v", resp["status"])
	}
	if savedID != "abc-456" {
		t.Fatalf("expected savedID 'abc-456', got: %s", savedID)
	}
	if savedData["title"] != "Updated Product" {
		t.Fatalf("expected title 'Updated Product', got: %s", savedData["title"])
	}
}
