package crud

import (
	"strings"
	"testing"
)

// buildWidget is a helper that calls buildFieldWidget with a dummy ID.
func buildWidget(crud *Crud, field FieldInterface) string {
	return crud.buildFieldWidget(field, "test-id").ToHTML()
}

// buildSubWidget is a helper that calls buildSubFieldWidget with a fixed vModel.
func buildSubWidget(crud *Crud, field FieldInterface) string {
	return crud.buildSubFieldWidget(field, "entityModel.name").ToHTML()
}

func newBuilderCrud() *Crud {
	c := newTestCrud()
	return &c
}

// == Plain HTML fields =======================================================

func TestBuildFieldWidget_String(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewStringField("name", "Name"))
	if !strings.Contains(html, `<input`) || !strings.Contains(html, `form-control`) {
		t.Fatalf("expected text input with form-control, got: %s", html)
	}
}

func TestBuildFieldWidget_Email(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewEmailField("email", "Email"))
	if !strings.Contains(html, `type="email"`) {
		t.Fatalf("expected email input, got: %s", html)
	}
}

func TestBuildFieldWidget_Number(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewNumberField("qty", "Qty"))
	if !strings.Contains(html, `type="number"`) {
		t.Fatalf("expected number input, got: %s", html)
	}
}

func TestBuildFieldWidget_Password(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewPasswordField("pass", "Pass"))
	if !strings.Contains(html, `type="password"`) {
		t.Fatalf("expected password input, got: %s", html)
	}
}

func TestBuildFieldWidget_Date(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewDateField("dob", "DOB"))
	if !strings.Contains(html, `type="date"`) {
		t.Fatalf("expected date input, got: %s", html)
	}
	if strings.Contains(html, "el-date-picker") {
		t.Fatalf("plain date should not use el-date-picker, got: %s", html)
	}
}

func TestBuildFieldWidget_DateTime(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewDateTimeField("ts", "Timestamp"))
	if !strings.Contains(html, `type="datetime-local"`) {
		t.Fatalf("expected datetime-local input, got: %s", html)
	}
	if strings.Contains(html, "el-date-picker") {
		t.Fatalf("plain datetime should not use el-date-picker, got: %s", html)
	}
}

func TestBuildFieldWidget_Color(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewColorField("bg", "BG"))
	if !strings.Contains(html, `type="color"`) {
		t.Fatalf("expected color input, got: %s", html)
	}
	if strings.Contains(html, "el-color-picker") {
		t.Fatalf("plain color should not use el-color-picker, got: %s", html)
	}
}

func TestBuildFieldWidget_Select(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewSelectField("status", "Status", []FieldOption{
		{Key: "active", Value: "Active"},
		{Key: "inactive", Value: "Inactive"},
	}))
	if !strings.Contains(html, "<select") {
		t.Fatalf("expected <select>, got: %s", html)
	}
	if !strings.Contains(html, "Active") || !strings.Contains(html, "Inactive") {
		t.Fatalf("expected options in select, got: %s", html)
	}
}

func TestBuildFieldWidget_Radio(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewRadioField("size", "Size", []FieldOption{
		{Key: "s", Value: "Small"},
		{Key: "l", Value: "Large"},
	}))
	if !strings.Contains(html, `type="radio"`) {
		t.Fatalf("expected radio inputs, got: %s", html)
	}
}

func TestBuildFieldWidget_Checkbox(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewCheckboxField("active", "Active"))
	if !strings.Contains(html, `type="checkbox"`) {
		t.Fatalf("expected checkbox input, got: %s", html)
	}
}

func TestBuildFieldWidget_Textarea(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewTextAreaField("bio", "Bio"))
	if !strings.Contains(html, "<textarea") {
		t.Fatalf("expected textarea, got: %s", html)
	}
}

func TestBuildFieldWidget_Hidden(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewHiddenField("token", "abc123"))
	if !strings.Contains(html, `type="hidden"`) {
		t.Fatalf("expected hidden input, got: %s", html)
	}
}

func TestBuildFieldWidget_Tel(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewTelField("phone", "Phone"))
	if !strings.Contains(html, `type="tel"`) {
		t.Fatalf("expected tel input, got: %s", html)
	}
}

func TestBuildFieldWidget_URL(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewURLField("site", "Site"))
	if !strings.Contains(html, `type="url"`) {
		t.Fatalf("expected url input, got: %s", html)
	}
}

// == Element Plus fields =====================================================

func TestBuildFieldWidget_ColorEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewColorElField("bg", "BG"))
	if !strings.Contains(html, "el-color-picker") {
		t.Fatalf("expected el-color-picker, got: %s", html)
	}
}

func TestBuildFieldWidget_DateEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewDateElField("dob", "DOB"))
	if !strings.Contains(html, "el-date-picker") {
		t.Fatalf("expected el-date-picker, got: %s", html)
	}
	if !strings.Contains(html, `type="date"`) {
		t.Fatalf("expected type=date on el-date-picker, got: %s", html)
	}
}

func TestBuildFieldWidget_DateTimeEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewDateTimeElField("ts", "Timestamp"))
	if !strings.Contains(html, "el-date-picker") {
		t.Fatalf("expected el-date-picker, got: %s", html)
	}
	if !strings.Contains(html, `type="datetime"`) {
		t.Fatalf("expected type=datetime on el-date-picker, got: %s", html)
	}
}

func TestBuildFieldWidget_InputNumberEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewInputNumberElField("qty", "Qty"))
	if !strings.Contains(html, "el-input-number") {
		t.Fatalf("expected el-input-number, got: %s", html)
	}
}

func TestBuildFieldWidget_RateEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewRateElField("score", "Score"))
	if !strings.Contains(html, "el-rate") {
		t.Fatalf("expected el-rate, got: %s", html)
	}
}

func TestBuildFieldWidget_SelectEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewSelectElField("status", "Status", []FieldOption{
		{Key: "active", Value: "Active"},
		{Key: "inactive", Value: "Inactive"},
	}))
	if !strings.Contains(html, "el-select") {
		t.Fatalf("expected el-select, got: %s", html)
	}
	if !strings.Contains(html, "el-option") {
		t.Fatalf("expected el-option children, got: %s", html)
	}
	if !strings.Contains(html, "filterable") {
		t.Fatalf("expected filterable attr, got: %s", html)
	}
}

func TestBuildFieldWidget_SliderEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewSliderElField("vol", "Volume"))
	if !strings.Contains(html, "el-slider") {
		t.Fatalf("expected el-slider, got: %s", html)
	}
}

func TestBuildFieldWidget_SwitchEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewSwitchElField("enabled", "Enabled"))
	if !strings.Contains(html, "el-switch") {
		t.Fatalf("expected el-switch, got: %s", html)
	}
}

func TestBuildFieldWidget_TimeEl(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewTimeElField("start", "Start Time"))
	if !strings.Contains(html, "el-time-picker") {
		t.Fatalf("expected el-time-picker, got: %s", html)
	}
}

// == Sub-field widgets (inside repeaters) ====================================

func TestBuildSubFieldWidget_Date_IsPlainHTML(t *testing.T) {
	html := buildSubWidget(newBuilderCrud(), NewDateField("dob", "DOB"))
	if !strings.Contains(html, `type="date"`) {
		t.Fatalf("expected plain date input, got: %s", html)
	}
	if strings.Contains(html, "el-date-picker") {
		t.Fatalf("plain date sub-field should not use el-date-picker, got: %s", html)
	}
}

func TestBuildSubFieldWidget_DateTime_IsPlainHTML(t *testing.T) {
	html := buildSubWidget(newBuilderCrud(), NewDateTimeField("ts", "Timestamp"))
	if !strings.Contains(html, `type="datetime-local"`) {
		t.Fatalf("expected plain datetime-local input, got: %s", html)
	}
	if strings.Contains(html, "el-date-picker") {
		t.Fatalf("plain datetime sub-field should not use el-date-picker, got: %s", html)
	}
}

func TestBuildSubFieldWidget_DateEl(t *testing.T) {
	html := buildSubWidget(newBuilderCrud(), NewDateElField("dob", "DOB"))
	if !strings.Contains(html, "el-date-picker") {
		t.Fatalf("expected el-date-picker in sub-field, got: %s", html)
	}
}

func TestBuildSubFieldWidget_DateTimeEl(t *testing.T) {
	html := buildSubWidget(newBuilderCrud(), NewDateTimeElField("ts", "Timestamp"))
	if !strings.Contains(html, "el-date-picker") {
		t.Fatalf("expected el-date-picker in sub-field, got: %s", html)
	}
}

func TestBuildSubFieldWidget_SelectEl(t *testing.T) {
	html := buildSubWidget(newBuilderCrud(), NewSelectElField("status", "Status", []FieldOption{
		{Key: "y", Value: "Yes"},
	}))
	if !strings.Contains(html, "el-select") {
		t.Fatalf("expected el-select in sub-field, got: %s", html)
	}
	if !strings.Contains(html, "el-option") {
		t.Fatalf("expected el-option in sub-field, got: %s", html)
	}
}

// == vModel binding ==========================================================

func TestBuildFieldWidget_BindsVModel(t *testing.T) {
	html := buildWidget(newBuilderCrud(), NewStringField("title", "Title"))
	if !strings.Contains(html, `v-model="entityModel.title"`) {
		t.Fatalf("expected v-model binding, got: %s", html)
	}
}

func TestBuildSubFieldWidget_BindsVModel(t *testing.T) {
	html := buildSubWidget(newBuilderCrud(), NewStringField("title", "Title"))
	if !strings.Contains(html, `v-model="entityModel.name"`) {
		t.Fatalf("expected v-model binding, got: %s", html)
	}
}
