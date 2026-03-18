package crud

import (
	"strings"
	"testing"

	
)

func TestRepeaterFieldRendering(t *testing.T) {
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:     "features",
				Label:    "Product Features",
				Type:     FORM_FIELD_TYPE_REPEATER,
				Required: true,
				Help:     "Add multiple features for this product",
			}),
		},
	}

	// Test that repeater field renders correctly
	tags := crud.form(crud.createFields)
	if len(tags) != 1 {
		t.Errorf("Expected 1 form group, got %d", len(tags))
	}

	// Check that the repeater container is present
	formGroup := tags[0]
	if formGroup == nil {
		t.Error("Form group should not be nil")
	}

	// Check that the repeater container has the correct class
	formGroupHTML := formGroup.ToHTML()
	if !strings.Contains(formGroupHTML, "repeater-container") {
		t.Error("Should contain repeater container")
	}
	if !strings.Contains(formGroupHTML, "repeater-items") {
		t.Error("Should contain repeater items container")
	}
	if !strings.Contains(formGroupHTML, "addRepeaterItem") {
		t.Error("Should contain add repeater item function call")
	}
	if !strings.Contains(formGroupHTML, "removeRepeaterItem") {
		t.Error("Should contain remove repeater item function call")
	}
	if !strings.Contains(formGroupHTML, "entityModel.features") {
		t.Error("Should bind to entityModel.features")
	}
}

func TestRepeaterFieldValidation(t *testing.T) {
	crud := &Crud{}

	// Test required repeater field validation
	requiredRepeaterField := NewField(FieldOptions{
		Name:     "features",
		Label:    "Product Features",
		Type:     FORM_FIELD_TYPE_REPEATER,
		Required: true,
	})

	// Test with empty data
	data := map[string]string{
		"features": "",
	}
	err := crud.validateRepeaterFields(data, []FieldInterface{requiredRepeaterField})
	if err == nil {
		t.Error("Should return error for empty repeater field")
	}
	if !strings.Contains(err.Error(), "Product Features is required field") {
		t.Error("Should contain field label in error")
	}

	// Test with empty array
	data["features"] = "[]"
	err = crud.validateRepeaterFields(data, []FieldInterface{requiredRepeaterField})
	if err == nil {
		t.Error("Should return error for empty array repeater field")
	}
	if !strings.Contains(err.Error(), "Product Features is required field") {
		t.Error("Should contain field label in error")
	}

	// Test with valid data
	data["features"] = `[{"name": "Feature 1", "description": "Description 1"}]`
	err = crud.validateRepeaterFields(data, []FieldInterface{requiredRepeaterField})
	if err != nil {
		t.Errorf("Should not return error for valid repeater field data: %v", err)
	}

	// Test with non-required repeater field
	nonRequiredRepeaterField := NewField(FieldOptions{
		Name:     "optional_features",
		Label:    "Optional Features",
		Type:     FORM_FIELD_TYPE_REPEATER,
		Required: false,
	})

	data["optional_features"] = ""
	err = crud.validateRepeaterFields(data, []FieldInterface{nonRequiredRepeaterField})
	if err != nil {
		t.Errorf("Should not return error for empty non-required repeater field: %v", err)
	}
}

func TestRepeaterFieldWithMultipleFields(t *testing.T) {
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:  "name",
				Label: "Product Name",
				Type:  FORM_FIELD_TYPE_STRING,
			}),
			NewField(FieldOptions{
				Name:     "features",
				Label:    "Product Features",
				Type:     FORM_FIELD_TYPE_REPEATER,
				Required: true,
			}),
			NewField(FieldOptions{
				Name:  "description",
				Label: "Description",
				Type:  FORM_FIELD_TYPE_TEXTAREA,
			}),
		},
	}

	// Test that all fields render correctly including repeater
	tags := crud.form(crud.createFields)
	if len(tags) != 3 {
		t.Errorf("Expected 3 form groups, got %d", len(tags))
	}

	// Check that the repeater field is in the correct position
	repeaterFormGroup := tags[1]
	repeaterHTML := repeaterFormGroup.ToHTML()
	if !strings.Contains(repeaterHTML, "repeater-container") {
		t.Error("Second form group should be repeater")
	}
	if !strings.Contains(repeaterHTML, "entityModel.features") {
		t.Error("Should bind to features field")
	}
}

func TestRepeaterFieldRenderingWithHelpText(t *testing.T) {
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:     "features",
				Label:    "Product Features",
				Type:     FORM_FIELD_TYPE_REPEATER,
				Required: true,
				Help:     "Add multiple features for this product. Each feature can have a name and description.",
			}),
		},
	}

	tags := crud.form(crud.createFields)
	formGroupHTML := tags[0].ToHTML()

	// Check that help text is included
	if !strings.Contains(formGroupHTML, "text-info") {
		t.Error("Should contain help text with text-info class")
	}
	if !strings.Contains(formGroupHTML, "Add multiple features for this product") {
		t.Error("Should contain help text content")
	}
}

func TestRepeaterFieldRenderingWithRequiredIndicator(t *testing.T) {
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:     "features",
				Label:    "Product Features",
				Type:     FORM_FIELD_TYPE_REPEATER,
				Required: true,
			}),
		},
	}

	tags := crud.form(crud.createFields)
	formGroupHTML := tags[0].ToHTML()

	// Check that required indicator is present
	if !strings.Contains(formGroupHTML, "text-danger") {
		t.Error("Should contain required indicator")
	}
	if !strings.Contains(formGroupHTML, "*") {
		t.Error("Should contain asterisk for required field")
	}
}

func TestRepeaterFieldVueMethodsIntegration(t *testing.T) {
	// Test that the Vue.js methods are properly integrated
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:     "features",
				Label:    "Product Features",
				Type:     FORM_FIELD_TYPE_REPEATER,
				Required: true,
			}),
		},
	}

	// Test that the repeater field renders with correct Vue.js bindings
	tags := crud.form(crud.createFields)
	formGroupHTML := tags[0].ToHTML()

	// Check for Vue.js directives specific to repeaters
	if !strings.Contains(formGroupHTML, "v-for") {
		t.Error("Should contain v-for directive for repeater items")
	}
	if !strings.Contains(formGroupHTML, "v-model") {
		t.Error("Should contain v-model binding for repeater data")
	}
	if !strings.Contains(formGroupHTML, ":key") {
		t.Error("Should contain key binding for Vue.js list rendering")
	}
}

func TestRepeaterFieldUpDownButtons(t *testing.T) {
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:  "items",
				Label: "Items",
				Type:  FORM_FIELD_TYPE_REPEATER,
			}),
		},
	}

	tags := crud.form(crud.createFields)
	html := tags[0].ToHTML()

	if !strings.Contains(html, "moveRepeaterItemUp") {
		t.Error("Should contain moveRepeaterItemUp call on up button")
	}
	if !strings.Contains(html, "moveRepeaterItemDown") {
		t.Error("Should contain moveRepeaterItemDown call on down button")
	}
	if !strings.Contains(html, "bi-arrow-up") {
		t.Error("Should contain up arrow icon")
	}
	if !strings.Contains(html, "bi-arrow-down") {
		t.Error("Should contain down arrow icon")
	}
	// Up button disabled when first item
	if !strings.Contains(html, "index === 0") {
		t.Error("Should disable up button for first item")
	}
	// Down button disabled when last item
	if !strings.Contains(html, "index === entityModel.items.length - 1") {
		t.Error("Should disable down button for last item")
	}
}

func TestRepeaterFieldAddItemPassesEmptyObject(t *testing.T) {
	crud := &Crud{
		createFields: []FieldInterface{
			NewField(FieldOptions{
				Name:  "links",
				Label: "Links",
				Type:  FORM_FIELD_TYPE_REPEATER,
			}),
		},
	}

	tags := crud.form(crud.createFields)
	html := tags[0].ToHTML()

	// Generic repeater (no sub-fields) passes "" as the item seed.
	// HTML attributes are escaped, so " becomes &#34; and ' becomes &#39;.
	if !strings.Contains(html, `addRepeaterItem(&#39;links&#39;, &#34;&#34;)`) {
		t.Errorf("Generic repeater should pass empty string seed to addRepeaterItem, got: %s", html)
	}
}

func TestRepeaterFieldDataStructure(t *testing.T) {
	// Test the expected data structure for repeater fields
	crud := &Crud{}

	// Simulate form data that would be submitted
	formData := map[string]string{
		"name": "Test Product",
		"features": `[
			{"name": "Feature 1", "description": "Description 1", "enabled": true},
			{"name": "Feature 2", "description": "Description 2", "enabled": false}
		]`,
		"description": "Test product description",
	}

	// Test that the data structure is valid JSON
	err := crud.validateRepeaterFields(formData, []FieldInterface{
		NewField(FieldOptions{
			Name:     "features",
			Label:    "Product Features",
			Type:     FORM_FIELD_TYPE_REPEATER,
			Required: true,
		}),
	})

	// The validation should pass for valid JSON structure
	if err != nil {
		t.Errorf("Should validate valid repeater field data: %v", err)
	}
}

func TestRepeaterFieldEdgeCases(t *testing.T) {
	crud := &Crud{}

	requiredRepeaterField := NewField(FieldOptions{
		Name:     "features",
		Label:    "Product Features",
		Type:     FORM_FIELD_TYPE_REPEATER,
		Required: true,
	})

	// Test various edge cases for repeater field validation
	testCases := []struct {
		name     string
		data     string
		shouldErr bool
	}{
		{"Empty string", "", true},
		{"Empty array", "[]", true},
		{"Whitespace only", "   ", true},
		{"Single item", `[{"name": "Feature 1"}]`, false},
		{"Multiple items", `[{"name": "Feature 1"}, {"name": "Feature 2"}]`, false},
		{"Nested objects", `[{"name": "Feature 1", "options": {"color": "red"}}]`, false},
		{"Invalid JSON", `invalid json`, true},
		{"Null value", "null", true},
		{"Number", "123", true},
		{"Boolean", "true", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]string{
				"features": tc.data,
			}

			err := crud.validateRepeaterFields(data, []FieldInterface{requiredRepeaterField})

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Should return error for: %s", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Should not return error for: %s, got: %v", tc.name, err)
				}
			}
		})
	}
}
