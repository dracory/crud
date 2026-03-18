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

func TestCreate_ModalSave_MethodNotAllowed_GET(t *testing.T) {
	crud := newTestCrud()
	ctrl := crud.newEntityCreateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-create-ajax", nil)
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

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

func TestCreate_ModalSave_NilFuncCreate(t *testing.T) {
	crud := newTestCrud()
	crud.funcCreate = nil
	ctrl := crud.newEntityCreateController()

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", nil)
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Create functionality is not configured") {
		t.Fatalf("expected 'Create functionality is not configured' message, got: %v", resp["message"])
	}
}

func TestCreate_ModalSave_RequiredFieldEmpty(t *testing.T) {
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:     "title",
			Label:    "Title",
			Type:     FORM_FIELD_TYPE_STRING,
			Required: true,
		}),
	}
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "new-id", nil
	}
	ctrl := crud.newEntityCreateController()

	formData := url.Values{}
	formData.Set("title", "")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "is required field") {
		t.Fatalf("expected 'is required field' message, got: %v", resp["message"])
	}
}

func TestCreate_ModalSave_FuncCreateError(t *testing.T) {
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
		}),
	}
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "", errors.New("save error")
	}
	ctrl := crud.newEntityCreateController()

	formData := url.Values{}
	formData.Set("title", "Test Product")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "error" {
		t.Fatalf("expected status 'error', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Save failed") {
		t.Fatalf("expected 'Save failed' message, got: %v", resp["message"])
	}
}

func TestCreate_ModalSave_Success(t *testing.T) {
	var receivedData map[string]string
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
		}),
	}
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		receivedData = data
		return "new-entity-id", nil
	}
	ctrl := crud.newEntityCreateController()

	formData := url.Values{}
	formData.Set("title", "Test Product")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "success" {
		t.Fatalf("expected status 'success', got: %v", resp["status"])
	}
	if !strings.Contains(resp["message"].(string), "Saved successfully") {
		t.Fatalf("expected 'Saved successfully' message, got: %v", resp["message"])
	}
	data := resp["data"].(map[string]interface{})
	if data["entity_id"] != "new-entity-id" {
		t.Fatalf("expected entity_id 'new-entity-id', got: %v", data["entity_id"])
	}
	if data["redirect_url"] == nil || data["redirect_url"] == "" {
		t.Fatal("expected redirect_url in response data")
	}
	if receivedData["title"] != "Test Product" {
		t.Fatalf("expected title 'Test Product', got: %s", receivedData["title"])
	}
}

func TestCreate_ModalSave_DefaultRedirectURL(t *testing.T) {
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
	}
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "abc-123", nil
	}
	ctrl := crud.newEntityCreateController()

	formData := url.Values{}
	formData.Set("title", "Test")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	redirectURL := data["redirect_url"].(string)
	if !strings.Contains(redirectURL, "entity-update") || !strings.Contains(redirectURL, "abc-123") {
		t.Fatalf("expected default redirect to update page with entity ID, got: %s", redirectURL)
	}
}

func TestCreate_ModalSave_CustomRedirectURL(t *testing.T) {
	crud := newTestCrud()
	crud.createRedirectURL = "/custom/redirect"
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
	}
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "abc-123", nil
	}
	ctrl := crud.newEntityCreateController()

	formData := url.Values{}
	formData.Set("title", "Test")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	redirectURL := data["redirect_url"].(string)
	if redirectURL != "/custom/redirect" {
		t.Fatalf("expected custom redirect URL '/custom/redirect', got: %s", redirectURL)
	}
}

func TestCreate_ModalShow_ReturnsHTML(t *testing.T) {
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Name:  "title",
			Label: "Title",
			Type:  FORM_FIELD_TYPE_STRING,
		}),
	}
	ctrl := crud.newEntityCreateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-create-modal", nil)
	w := httptest.NewRecorder()

	ctrl.modalShow(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "ModalEntityCreate") {
		t.Fatal("expected 'ModalEntityCreate' in modal HTML output")
	}
	if !strings.Contains(body, "New Product") {
		t.Fatal("expected 'New Product' in modal HTML output")
	}
}

func TestCreate_ModalSave_RepeaterFieldCollected(t *testing.T) {
	var savedData map[string]string
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
		form.NewField(form.FieldOptions{Name: "tags", Type: FORM_FIELD_TYPE_REPEATER}),
	}
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		savedData = data
		return "new-id", nil
	}
	ctrl := crud.newEntityCreateController()

	formData := url.Values{}
	formData.Set("title", "My Product")
	formData.Add("tags[0][name]", "go")
	formData.Add("tags[1][name]", "web")
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ctrl.modalSave(w, r)

	body := w.Body.String()
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("expected JSON response, got: %s", body)
	}
	if resp["status"] != "success" {
		t.Fatalf("expected status 'success', got: %v — body: %s", resp["status"], body)
	}
	if savedData["title"] != "My Product" {
		t.Fatalf("expected title 'My Product', got: %s", savedData["title"])
	}
	if !strings.HasPrefix(savedData["tags"], "[") {
		t.Fatalf("expected tags to be JSON array, got: %s", savedData["tags"])
	}
	if !strings.Contains(savedData["tags"], `"name"`) {
		t.Fatalf("expected 'name' key in tags JSON, got: %s", savedData["tags"])
	}
}

func TestCreate_Modal_ContainsMoveRepeaterMethods(t *testing.T) {
	crud := newTestCrud()
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
	}
	ctrl := crud.newEntityCreateController()

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-create-modal", nil)
	w := httptest.NewRecorder()

	ctrl.modalShow(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "moveRepeaterItemUp") {
		t.Fatal("expected 'moveRepeaterItemUp' method in modal script")
	}
	if !strings.Contains(body, "moveRepeaterItemDown") {
		t.Fatal("expected 'moveRepeaterItemDown' method in modal script")
	}
}
