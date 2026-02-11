package crud

import (
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
	if !strings.Contains(body, "Method not allowed") {
		t.Fatalf("expected 'Method not allowed' in response, got: %s", body)
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
	if !strings.Contains(body, "Create functionality is not configured") {
		t.Fatalf("expected 'Create functionality is not configured' in response, got: %s", body)
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
	if !strings.Contains(body, "is required field") {
		t.Fatalf("expected 'is required field' in response, got: %s", body)
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
	if !strings.Contains(body, "Save failed") {
		t.Fatalf("expected 'Save failed' in response, got: %s", body)
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
	if !strings.Contains(body, "Saved successfully") {
		t.Fatalf("expected 'Saved successfully' in response, got: %s", body)
	}
	if receivedData["title"] != "Test Product" {
		t.Fatalf("expected title 'Test Product', got: %s", receivedData["title"])
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
