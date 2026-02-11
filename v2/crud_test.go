package crud

import (
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
		funcRows:           func() ([]Row, error) { return []Row{}, nil },
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
