package crud

import (
	"net/http"
	"testing"

	
)

func TestNew_MissingFuncRows(t *testing.T) {
	_, err := New(Config{})
	if err == nil {
		t.Fatal("expected error when FuncRows is nil")
	}
	if err.Error() != "FuncRows function is required" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestNew_MissingUpdateFields(t *testing.T) {
	_, err := New(Config{
		FuncRows: func(r *http.Request) ([]Row, error) { return nil, nil },
	})
	if err == nil {
		t.Fatal("expected error when UpdateFields is nil")
	}
	if err.Error() != "UpdateFields is required" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestNew_UpdateFieldsWithoutFuncUpdate(t *testing.T) {
	_, err := New(Config{
		FuncRows: func(r *http.Request) ([]Row, error) { return nil, nil },
		UpdateFields: []FieldInterface{
			NewField(FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
		},
	})
	if err == nil {
		t.Fatal("expected error when UpdateFields provided without FuncUpdate")
	}
	if err.Error() != "FuncUpdate function is required when UpdateFields are provided" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestNew_UpdateFieldsWithoutFuncFetchUpdateData(t *testing.T) {
	_, err := New(Config{
		FuncRows: func(r *http.Request) ([]Row, error) { return nil, nil },
		UpdateFields: []FieldInterface{
			NewField(FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
		},
		FuncUpdate: func(r *http.Request, entityID string, data map[string]string) error { return nil },
	})
	if err == nil {
		t.Fatal("expected error when UpdateFields provided without FuncFetchUpdateData")
	}
	if err.Error() != "FuncFetchUpdateData function is required when UpdateFields are provided" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestNew_EmptyUpdateFieldsNoFuncUpdateRequired(t *testing.T) {
	_, err := New(Config{
		FuncRows:     func(r *http.Request) ([]Row, error) { return nil, nil },
		UpdateFields: []FieldInterface{},
	})
	if err != nil {
		t.Fatalf("expected no error with empty UpdateFields, got: %s", err.Error())
	}
}

func TestNew_ValidFullConfig(t *testing.T) {
	crud, err := New(Config{
		Endpoint:           "/admin",
		EntityNameSingular: "Product",
		EntityNamePlural:   "Products",
		HomeURL:            "/",
		ColumnNames:        []string{"ID", "Name"},
		FuncRows:           func(r *http.Request) ([]Row, error) { return nil, nil },
		FuncUpdate:         func(r *http.Request, entityID string, data map[string]string) error { return nil },
		FuncFetchUpdateData: func(r *http.Request, entityID string) (map[string]string, error) {
			return map[string]string{}, nil
		},
		UpdateFields: []FieldInterface{
			NewField(FieldOptions{Name: "title", Type: FORM_FIELD_TYPE_STRING}),
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}
	if crud.endpoint != "/admin" {
		t.Fatalf("expected endpoint '/admin', got: %s", crud.endpoint)
	}
	if crud.entityNameSingular != "Product" {
		t.Fatalf("expected entityNameSingular 'Product', got: %s", crud.entityNameSingular)
	}
	if crud.entityNamePlural != "Products" {
		t.Fatalf("expected entityNamePlural 'Products', got: %s", crud.entityNamePlural)
	}
	if crud.homeURL != "/" {
		t.Fatalf("expected homeURL '/', got: %s", crud.homeURL)
	}
	if len(crud.columnNames) != 2 {
		t.Fatalf("expected 2 column names, got: %d", len(crud.columnNames))
	}
	if len(crud.updateFields) != 1 {
		t.Fatalf("expected 1 update field, got: %d", len(crud.updateFields))
	}
}

func TestNew_ConfigFieldsAreCopied(t *testing.T) {
	createFields := []FieldInterface{
		NewField(FieldOptions{Name: "name", Type: FORM_FIELD_TYPE_STRING}),
	}
	updateFields := []FieldInterface{
		NewField(FieldOptions{Name: "name", Type: FORM_FIELD_TYPE_STRING}),
	}

	crud, err := New(Config{
		FuncRows:            func(r *http.Request) ([]Row, error) { return nil, nil },
		FuncUpdate:          func(r *http.Request, entityID string, data map[string]string) error { return nil },
		FuncFetchUpdateData: func(r *http.Request, entityID string) (map[string]string, error) { return nil, nil },
		CreateFields:        createFields,
		UpdateFields:        updateFields,
		FileManagerURL:      "/files",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if len(crud.createFields) != 1 {
		t.Fatalf("expected 1 create field, got: %d", len(crud.createFields))
	}
	if crud.fileManagerURL != "/files" {
		t.Fatalf("expected fileManagerURL '/files', got: %s", crud.fileManagerURL)
	}
}
