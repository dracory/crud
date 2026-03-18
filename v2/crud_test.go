package crud

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/form"
)

// ==========================================================================
// URL generation tests
// ==========================================================================

func TestUrlHome(t *testing.T) {
	crud := Crud{homeURL: "/dashboard"}
	if crud.urlHome() != "/dashboard" {
		t.Fatalf("expected '/dashboard', got: %s", crud.urlHome())
	}
}

func TestUrlEntityManager_WithoutQueryString(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityManager()
	expected := "/admin?path=entity-manager"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityManager_WithQueryString(t *testing.T) {
	crud := Crud{endpoint: "/admin?module=cms"}
	got := crud.UrlEntityManager()
	expected := "/admin?module=cms&path=entity-manager"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityCreateModal(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityCreateModal()
	expected := "/admin?path=entity-create-modal"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityCreateAjax(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityCreateAjax()
	expected := "/admin?path=entity-create-ajax"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityTrashAjax(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityTrashAjax()
	expected := "/admin?path=entity-trash-ajax"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityRead(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityRead()
	expected := "/admin?path=entity-read"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityUpdate(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityUpdate()
	expected := "/admin?path=entity-update"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

func TestUrlEntityUpdateAjax(t *testing.T) {
	crud := Crud{endpoint: "/admin"}
	got := crud.UrlEntityUpdateAjax()
	expected := "/admin?path=entity-update-ajax"
	if got != expected {
		t.Fatalf("expected %q, got: %q", expected, got)
	}
}

// ==========================================================================
// Breadcrumb rendering tests
// ==========================================================================

func TestRenderBreadcrumbs_Empty(t *testing.T) {
	crud := Crud{}
	html := crud.renderBreadcrumbs([]Breadcrumb{})
	if !strings.Contains(html, "breadcrumb") {
		t.Fatal("expected breadcrumb class in output")
	}
	if strings.Contains(html, "breadcrumb-item") {
		t.Fatal("expected no breadcrumb items for empty input")
	}
}

func TestRenderBreadcrumbs_SingleItem(t *testing.T) {
	crud := Crud{}
	html := crud.renderBreadcrumbs([]Breadcrumb{
		{Name: "Home", URL: "/"},
	})
	if !strings.Contains(html, "Home") {
		t.Fatal("expected 'Home' text in output")
	}
	if !strings.Contains(html, `href="/"`) {
		t.Fatal("expected href='/' in output")
	}
	if !strings.Contains(html, "breadcrumb-item") {
		t.Fatal("expected breadcrumb-item class")
	}
}

func TestRenderBreadcrumbs_MultipleItems(t *testing.T) {
	crud := Crud{}
	html := crud.renderBreadcrumbs([]Breadcrumb{
		{Name: "Home", URL: "/"},
		{Name: "Products", URL: "/products"},
		{Name: "Edit", URL: "/products/edit?id=1"},
	})
	if !strings.Contains(html, "Home") {
		t.Fatal("expected 'Home' in output")
	}
	if !strings.Contains(html, "Products") {
		t.Fatal("expected 'Products' in output")
	}
	if !strings.Contains(html, "Edit") {
		t.Fatal("expected 'Edit' in output")
	}
	if strings.Count(html, "breadcrumb-item") != 3 {
		t.Fatalf("expected 3 breadcrumb items, got: %d", strings.Count(html, "breadcrumb-item"))
	}
}

// ==========================================================================
// listCreateNames / listUpdateNames tests
// ==========================================================================

func TestListCreateNames_Empty(t *testing.T) {
	crud := Crud{}
	names := crud.listCreateNames()
	if len(names) != 0 {
		t.Fatalf("expected 0 names, got: %d", len(names))
	}
}

func TestListCreateNames_SkipsEmptyNames(t *testing.T) {
	crud := Crud{
		createFields: []form.FieldInterface{
			form.NewField(form.FieldOptions{Name: "title"}),
			form.NewField(form.FieldOptions{Name: ""}),
			form.NewField(form.FieldOptions{Name: "description"}),
		},
	}
	names := crud.listCreateNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got: %d", len(names))
	}
	if names[0] != "title" {
		t.Fatalf("expected 'title', got: %s", names[0])
	}
	if names[1] != "description" {
		t.Fatalf("expected 'description', got: %s", names[1])
	}
}

func TestListUpdateNames_Empty(t *testing.T) {
	crud := Crud{}
	names := crud.listUpdateNames()
	if len(names) != 0 {
		t.Fatalf("expected 0 names, got: %d", len(names))
	}
}

func TestListUpdateNames_SkipsEmptyNames(t *testing.T) {
	crud := Crud{
		updateFields: []form.FieldInterface{
			form.NewField(form.FieldOptions{Name: "name"}),
			form.NewField(form.FieldOptions{Name: ""}),
			form.NewField(form.FieldOptions{Name: "status"}),
		},
	}
	names := crud.listUpdateNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got: %d", len(names))
	}
	if names[0] != "name" {
		t.Fatalf("expected 'name', got: %s", names[0])
	}
	if names[1] != "status" {
		t.Fatalf("expected 'status', got: %s", names[1])
	}
}

// ==========================================================================
// Handler routing tests
// ==========================================================================

func newTestCrud() Crud {
	return Crud{
		endpoint:           "/admin",
		entityNameSingular: "Product",
		entityNamePlural:   "Products",
		homeURL:            "/",
		columnNames:        []string{"ID", "Name"},
		funcRows:           func(r *http.Request) ([]Row, error) { return []Row{}, nil },
		updateFields:       []form.FieldInterface{},
	}
}

func TestHandler_DefaultRouteIsEntityManager(t *testing.T) {
	crud := newTestCrud()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Product Manager") {
		t.Fatal("expected 'Product Manager' in default route response")
	}
}

func TestHandler_EntityManagerRoute(t *testing.T) {
	crud := newTestCrud()
	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Product Manager") {
		t.Fatal("expected 'Product Manager' in entity-manager route response")
	}
}

func TestHandler_UnknownRouteFallsBackToHome(t *testing.T) {
	crud := newTestCrud()
	r := httptest.NewRequest(http.MethodGet, "/admin?path=unknown-route", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Product Manager") {
		t.Fatal("expected fallback to entity manager page")
	}
}

// ==========================================================================
// Middleware hooks tests
// ==========================================================================

func TestHandler_FuncBeforeAction_Called(t *testing.T) {
	var calledAction string
	crud := newTestCrud()
	crud.funcBeforeAction = func(w http.ResponseWriter, r *http.Request, action string) bool {
		calledAction = action
		return true
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if calledAction != "entity-manager" {
		t.Fatalf("expected action 'entity-manager', got: %s", calledAction)
	}
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", w.Result().StatusCode)
	}
}

func TestHandler_FuncBeforeAction_Aborts(t *testing.T) {
	crud := newTestCrud()
	crud.funcBeforeAction = func(w http.ResponseWriter, r *http.Request, action string) bool {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Access denied"))
		return false
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got: %d", w.Result().StatusCode)
	}
	body := w.Body.String()
	if body != "Access denied" {
		t.Fatalf("expected 'Access denied', got: %s", body)
	}
}

func TestHandler_FuncAfterAction_Called(t *testing.T) {
	var calledAction string
	crud := newTestCrud()
	crud.funcAfterAction = func(w http.ResponseWriter, r *http.Request, action string) {
		calledAction = action
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if calledAction != "entity-manager" {
		t.Fatalf("expected action 'entity-manager', got: %s", calledAction)
	}
}

func TestHandler_FuncAfterAction_NotCalledWhenAborted(t *testing.T) {
	afterCalled := false
	crud := newTestCrud()
	crud.funcBeforeAction = func(w http.ResponseWriter, r *http.Request, action string) bool {
		return false
	}
	crud.funcAfterAction = func(w http.ResponseWriter, r *http.Request, action string) {
		afterCalled = true
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if afterCalled {
		t.Fatal("expected FuncAfterAction NOT to be called when FuncBeforeAction aborts")
	}
}

func TestHandler_FuncBeforeAction_ReceivesDefaultPath(t *testing.T) {
	var calledAction string
	crud := newTestCrud()
	crud.funcBeforeAction = func(w http.ResponseWriter, r *http.Request, action string) bool {
		calledAction = action
		return true
	}

	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if calledAction != "home" {
		t.Fatalf("expected action 'home' for default path, got: %s", calledAction)
	}
}

// ==========================================================================
// CSRF validation tests
// ==========================================================================

func TestHandler_CSRF_BlocksPostWhenValidationFails(t *testing.T) {
	crud := newTestCrud()
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "1", nil
	}
	crud.funcValidateCSRF = func(r *http.Request) error {
		return errors.New("invalid token")
	}

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "CSRF validation failed") {
		t.Fatalf("expected CSRF error in response, got: %s", body)
	}
	if !strings.Contains(body, "invalid token") {
		t.Fatalf("expected 'invalid token' in response, got: %s", body)
	}
}

func TestHandler_CSRF_AllowsPostWhenValidationPasses(t *testing.T) {
	crud := newTestCrud()
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "1", nil
	}
	crud.funcValidateCSRF = func(r *http.Request) error {
		return nil
	}
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
	}

	formData := "title=Test"
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	body := w.Body.String()
	if strings.Contains(body, "CSRF validation failed") {
		t.Fatalf("expected no CSRF error, got: %s", body)
	}
	if !strings.Contains(body, "Saved successfully") {
		t.Fatalf("expected 'Saved successfully', got: %s", body)
	}
}

func TestHandler_CSRF_SkippedOnGetRequests(t *testing.T) {
	csrfCalled := false
	crud := newTestCrud()
	crud.funcValidateCSRF = func(r *http.Request) error {
		csrfCalled = true
		return errors.New("should not be called")
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if csrfCalled {
		t.Fatal("expected CSRF validation NOT to be called on GET requests")
	}
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", w.Result().StatusCode)
	}
}

func TestHandler_CSRF_SkippedWhenNil(t *testing.T) {
	crud := newTestCrud()
	crud.funcValidateCSRF = nil
	crud.funcCreate = func(r *http.Request, data map[string]string) (string, error) {
		return "1", nil
	}
	crud.createFields = []form.FieldInterface{
		form.NewField(form.FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
	}

	formData := "title=Test"
	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", strings.NewReader(formData))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "Saved successfully") {
		t.Fatalf("expected 'Saved successfully' when CSRF is nil, got: %s", body)
	}
}

// ==========================================================================
// Structured logging tests
// ==========================================================================

func TestHandler_FuncLog_CalledOnRequest(t *testing.T) {
	var loggedLevel, loggedMessage string
	var loggedAttrs map[string]any
	crud := newTestCrud()
	crud.funcLog = func(level string, message string, attrs map[string]any) {
		if message == "handling request" {
			loggedLevel = level
			loggedMessage = message
			loggedAttrs = attrs
		}
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if loggedLevel != LogLevelInfo {
		t.Fatalf("expected log level 'info', got: %s", loggedLevel)
	}
	if loggedMessage != "handling request" {
		t.Fatalf("expected message 'handling request', got: %s", loggedMessage)
	}
	if loggedAttrs["action"] != "entity-manager" {
		t.Fatalf("expected action 'entity-manager', got: %v", loggedAttrs["action"])
	}
	if loggedAttrs["method"] != "GET" {
		t.Fatalf("expected method 'GET', got: %v", loggedAttrs["method"])
	}
}

func TestHandler_FuncLog_NotCalledWhenNil(t *testing.T) {
	crud := newTestCrud()
	crud.funcLog = nil

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	// Should not panic
	crud.Handler(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", w.Result().StatusCode)
	}
}

func TestHandler_FuncLog_LogsBeforeActionAbort(t *testing.T) {
	var loggedAbort bool
	crud := newTestCrud()
	crud.funcBeforeAction = func(w http.ResponseWriter, r *http.Request, action string) bool {
		return false
	}
	crud.funcLog = func(level string, message string, attrs map[string]any) {
		if message == "request aborted by before-action hook" {
			loggedAbort = true
		}
	}

	r := httptest.NewRequest(http.MethodGet, "/admin?path=entity-manager", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if !loggedAbort {
		t.Fatal("expected abort to be logged")
	}
}

func TestHandler_FuncLog_LogsCSRFFailure(t *testing.T) {
	var loggedCSRF bool
	var loggedError string
	crud := newTestCrud()
	crud.funcValidateCSRF = func(r *http.Request) error {
		return errors.New("bad token")
	}
	crud.funcLog = func(level string, message string, attrs map[string]any) {
		if message == "CSRF validation failed" {
			loggedCSRF = true
			loggedError = attrs["error"].(string)
		}
	}

	r := httptest.NewRequest(http.MethodPost, "/admin?path=entity-create-ajax", nil)
	w := httptest.NewRecorder()

	crud.Handler(w, r)

	if !loggedCSRF {
		t.Fatal("expected CSRF failure to be logged")
	}
	if loggedError != "bad token" {
		t.Fatalf("expected error 'bad token', got: %s", loggedError)
	}
}

// ==========================================================================
// collectRepeaterField tests
// ==========================================================================

func TestCollectRepeaterField_NoKeys(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("title=hello"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	result := collectRepeaterField(r, "image_urls")
	if result != "[]" {
		t.Fatalf("expected '[]' when no keys match, got: %s", result)
	}
}

func TestCollectRepeaterField_SingleItemSingleKey(t *testing.T) {
	body := "image_urls%5B0%5D%5Bimage%5D=http%3A%2F%2Fexample.com%2Fa.jpg"
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	result := collectRepeaterField(r, "image_urls")
	if !strings.Contains(result, `"image"`) {
		t.Fatalf("expected key 'image' in result, got: %s", result)
	}
	if !strings.Contains(result, `http://example.com/a.jpg`) {
		t.Fatalf("expected URL value in result, got: %s", result)
	}
}

func TestCollectRepeaterField_MultipleItemsMultipleKeys(t *testing.T) {
	body := "image_urls%5B0%5D%5Bimage%5D=1&image_urls%5B1%5D%5Bimage%5D=2"
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	result := collectRepeaterField(r, "image_urls")
	// Should be a JSON array with 2 items
	if !strings.HasPrefix(result, "[") || !strings.HasSuffix(result, "]") {
		t.Fatalf("expected JSON array, got: %s", result)
	}
	if !strings.Contains(result, `"image":"1"`) && !strings.Contains(result, `"image": "1"`) {
		t.Fatalf("expected first item with image=1, got: %s", result)
	}
	if !strings.Contains(result, `"image":"2"`) && !strings.Contains(result, `"image": "2"`) {
		t.Fatalf("expected second item with image=2, got: %s", result)
	}
}

func TestCollectRepeaterField_PreservesOrder(t *testing.T) {
	body := "tags%5B0%5D%5Bname%5D=first&tags%5B1%5D%5Bname%5D=second&tags%5B2%5D%5Bname%5D=third"
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	result := collectRepeaterField(r, "tags")
	firstIdx := strings.Index(result, "first")
	secondIdx := strings.Index(result, "second")
	thirdIdx := strings.Index(result, "third")
	if firstIdx < 0 || secondIdx < 0 || thirdIdx < 0 {
		t.Fatalf("expected all three values in result, got: %s", result)
	}
	if !(firstIdx < secondIdx && secondIdx < thirdIdx) {
		t.Fatalf("expected order first < second < third, got: %s", result)
	}
}

func TestCollectRepeaterField_MultipleSubKeys(t *testing.T) {
	body := "items%5B0%5D%5Btitle%5D=Hello&items%5B0%5D%5Burl%5D=http%3A%2F%2Fexample.com"
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ParseForm()

	result := collectRepeaterField(r, "items")
	if !strings.Contains(result, `"title"`) {
		t.Fatalf("expected 'title' key in result, got: %s", result)
	}
	if !strings.Contains(result, `"url"`) {
		t.Fatalf("expected 'url' key in result, got: %s", result)
	}
	if !strings.Contains(result, "Hello") {
		t.Fatalf("expected 'Hello' value in result, got: %s", result)
	}
}

// ==========================================================================
// Layout tests
// ==========================================================================

func TestLayout_WithoutFuncLayout(t *testing.T) {
	crud := newTestCrud()
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	html := crud.layout(w, r, "Test Title", "<p>Content</p>", []string{}, "", []string{}, "")
	if !strings.Contains(html, "Test Title") {
		t.Fatal("expected title in layout output")
	}
	if !strings.Contains(html, "<p>Content</p>") {
		t.Fatal("expected content in layout output")
	}
}

func TestLayout_WithFuncLayout(t *testing.T) {
	crud := newTestCrud()
	crud.funcLayout = func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string {
		return "<html><head><title>" + title + "</title></head><body>" + content + "</body></html>"
	}
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	html := crud.layout(w, r, "Custom Title", "<p>Custom</p>", []string{}, "", []string{}, "")
	if !strings.Contains(html, "Custom Title") {
		t.Fatal("expected title in custom layout output")
	}
	if !strings.Contains(html, "<p>Custom</p>") {
		t.Fatal("expected content in custom layout output")
	}
}
