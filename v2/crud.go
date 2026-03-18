package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/str"
	"github.com/samber/lo"
)

type Crud struct {
	columnNames         []string
	createFields        []FieldInterface
	endpoint            string
	entityNamePlural    string
	entityNameSingular  string
	fileManagerURL      string
	funcCreate          func(r *http.Request, data map[string]string) (userID string, err error)
	funcReadExtras      func(r *http.Request, entityID string) []hb.TagInterface
	funcFetchReadData   func(r *http.Request, entityID string) ([]KeyValue, error)
	funcFetchUpdateData func(r *http.Request, entityID string) (map[string]string, error)
	funcLayout          func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string
	funcRows            func(r *http.Request) (rows []Row, err error)
	funcTrash           func(r *http.Request, entityID string) error
	funcUpdate          func(r *http.Request, entityID string, data map[string]string) error
	homeURL             string
	createRedirectURL   string
	updateRedirectURL   string
	updateFields        []FieldInterface
	pageSize            int
	funcRowsCount       func(r *http.Request) (int64, error)
	funcBeforeAction    func(w http.ResponseWriter, r *http.Request, action string) bool
	funcAfterAction     func(w http.ResponseWriter, r *http.Request, action string)
	funcValidateCSRF    func(r *http.Request) error
	funcLog             func(level string, message string, attrs map[string]any)
}

func (crud Crud) Handler(w http.ResponseWriter, r *http.Request) {
	path := req.GetStringTrimmed(r, "path")

	if path == "" {
		path = pathHome
	}

	crud.log(LogLevelInfo, "handling request", map[string]any{
		"action": path,
		"method": r.Method,
		"url":    r.URL.String(),
	})

	if crud.funcBeforeAction != nil {
		if !crud.funcBeforeAction(w, r, path) {
			crud.log(LogLevelWarn, "request aborted by before-action hook", map[string]any{
				"action": path,
			})
			return
		}
	}

	if crud.funcValidateCSRF != nil && r.Method == http.MethodPost {
		if err := crud.funcValidateCSRF(r); err != nil {
			crud.log(LogLevelWarn, "CSRF validation failed", map[string]any{
				"action": path,
				"error":  err.Error(),
			})
			api.Respond(w, r, api.Error("CSRF validation failed: "+err.Error()))
			return
		}
	}

	type contextKey string
	const ctxKeyPath contextKey = "path"
	ctx := context.WithValue(r.Context(), ctxKeyPath, r.URL.Path)

	routeFunc := crud.getRoute(path)
	routeFunc(w, r.WithContext(ctx))

	if crud.funcAfterAction != nil {
		crud.funcAfterAction(w, r, path)
	}
}

func (crud *Crud) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		pathHome:              crud.newEntityManagerController().page,
		pathEntityCreateAjax:  crud.newEntityCreateController().modalSave,
		pathEntityCreateModal: crud.newEntityCreateController().modalShow,
		pathEntityManager:     crud.newEntityManagerController().page,
		pathEntityRead:        crud.newEntityReadController().page,
		pathEntityUpdate:      crud.newEntityUpdateController().page,
		pathEntityUpdateAjax:  crud.newEntityUpdateController().pageSave,
		pathEntityTrashAjax:   crud.newEntityTrashController().pageEntityTrashAjax,
	}
	if val, ok := routes[route]; ok {
		return val
	}
	return routes[pathHome]
}

func (crud *Crud) urlHome() string {
	return crud.homeURL
}

func (crud *Crud) UrlEntityManager() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityManager
}

func (crud *Crud) UrlEntityCreateModal() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityCreateModal
}

func (crud *Crud) UrlEntityCreateAjax() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityCreateAjax
}

func (crud *Crud) UrlEntityTrashAjax() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityTrashAjax
}

func (crud *Crud) UrlEntityRead() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityRead
}

func (crud *Crud) UrlEntityUpdate() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityUpdate
}

func (crud *Crud) UrlEntityUpdateAjax() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	return crud.endpoint + q + "path=" + pathEntityUpdateAjax
}

// webpage returns the webpage template for the website
func (crud *Crud) webpage(title, content string) *hb.HtmlWebpage {
	faviconImgCms := `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`
	webpage := hb.Webpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(faviconImgCms)
	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_3_3(),
		cdn.BootstrapIconsCss_1_13_1(),
		cdn.VueElementPlusCss_2_3_8(),
	})
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_3_3(),
		cdn.Jquery_3_7_1(),
		cdn.VueJs_3(),
		cdn.Sweetalert2_11(),
	})
	webpage.AddScripts([]string{""})
	webpage.AddStyle(`html,body{height:100%;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}
	.form-select {
		display: block;
		width: 100%;
		padding: .375rem 2.25rem .375rem .75rem;
		font-size: 1rem;
		font-weight: 400;
		line-height: 1.5;
		color: #212529;
		background-color: #fff;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='none' stroke='%23343a40' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M2 5l6 6 6-6'/%3e%3c/svg%3e");
		background-repeat: no-repeat;
		background-position: right .75rem center;
		background-size: 16px 12px;
		border: 1px solid #ced4da;
		border-radius: .25rem;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
	}`)
	webpage.AddChild(hb.Raw(content))
	return webpage
}

func (crud *Crud) renderBreadcrumbs(breadcrumbs []Breadcrumb) string {
	nav := hb.Nav().Attr("aria-label", "breadcrumb")
	ol := hb.OL().Attr("class", "breadcrumb")
	for _, breadcrumb := range breadcrumbs {
		li := hb.LI().Attr("class", "breadcrumb-item")
		li.AddChild(hb.Hyperlink().Text(breadcrumb.Name).Attr("href", breadcrumb.URL))
		ol.AddChild(li)
	}
	nav.AddChild(ol)
	return nav.ToHTML()
}

func (crud *Crud) layout(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string {
	if crud.funcLayout != nil {
		jsFiles = append([]string{cdn.VueElementPlusJs_2_3_8()}, jsFiles...)
		jsFiles = append([]string{cdn.VueJs_3()}, jsFiles...)
		jsFiles = append([]string{cdn.Sweetalert2_11()}, jsFiles...)
		jsFiles = append([]string{cdn.Htmx_2_0_0()}, jsFiles...)
		styleFiles = append([]string{cdn.VueElementPlusCss_2_3_8()}, styleFiles...)
		return crud.funcLayout(w, r, title, content, styleFiles, style, jsFiles, js)
	}
	webpage := crud.webpage(title, content)
	webpage.AddStyleURLs(styleFiles)
	webpage.AddStyle(style)
	webpage.AddScriptURLs(jsFiles)
	webpage.AddScript(js)
	return webpage.ToHTML()
}

// form generates a form group for each field, delegating widget creation to buildFieldWidget.
func (crud *Crud) form(fields []FieldInterface) []hb.TagInterface {
	tags := []hb.TagInterface{}
	for _, field := range fields {
		fieldID := field.GetID()
		if fieldID == "" {
			id, err := str.RandomFromGamma(32, "abcdefghijklmnopqrstuvwxyz1234567890")
			if err != nil {
				id = "id_" + str.Random(32)
			}
			fieldID = id
		}

		fieldLabel := field.GetLabel()
		if fieldLabel == "" {
			fieldLabel = field.GetName()
		}

		formGroupInput := crud.buildFieldWidget(field, fieldID)

		formGroup := hb.Div().Class("form-group mb-3")
		if field.GetType() != FORM_FIELD_TYPE_RAW {
			formGroup.AddChild(hb.Label().
				Text(fieldLabel).
				Class("form-label").
				ChildIf(field.GetRequired(), hb.Sup().Text("*").Class("text-danger ml-1")))
		}
		formGroup.AddChild(formGroupInput)

		if field.GetHelp() != "" {
			formGroup.AddChild(hb.Paragraph().Class("text-info").HTML(field.GetHelp()))
		}

		tags = append(tags, formGroup)

		if field.GetType() == FORM_FIELD_TYPE_BLOCKAREA {
			script := hb.NewTag(`component`).
				Attr(`:is`, `'script'`).
				HTML(`setTimeout(() => {
				const blockArea = new BlockArea('` + fieldID + `');
				blockArea.registerBlock(BlockAreaHeading);
				blockArea.registerBlock(BlockAreaText);
				blockArea.registerBlock(BlockAreaImage);
				blockArea.registerBlock(BlockAreaCode);
				blockArea.registerBlock(BlockAreaRawHtml);
				blockArea.init();
			}, 2000)`)
			tags = append(tags, script)
		}
	}
	return tags
}

// listCreateNames returns the field names from createFields, including repeater fields.
func (crud *Crud) listCreateNames() []string {
	return listFieldNames(crud.createFields)
}

// listUpdateNames returns the field names from updateFields, including repeater fields.
func (crud *Crud) listUpdateNames() []string {
	return listFieldNames(crud.updateFields)
}

// listFieldNames collects all field names from a slice of fields.
// Repeater fields are included; their value arrives as a JSON-encoded string.
func listFieldNames(fields []FieldInterface) []string {
	names := []string{}
	for _, field := range fields {
		if field.GetName() == "" {
			continue
		}
		names = append(names, field.GetName())
	}
	return names
}

// collectRepeaterFields scans r.Form for keys matching fieldName[N][subKey]
// and returns a JSON-encoded array string, e.g. [{"image":"1"},{"image":"2"}].
// If no indexed keys are found it returns an empty JSON array "[]".
func collectRepeaterField(r *http.Request, fieldName string) string {
	// gather all sub-keys per index: rows[index][subKey] = value
	rows := map[int]map[string]string{}
	maxIndex := -1

	prefix := fieldName + "["
	for key, vals := range r.Form {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		// key looks like: fieldName[N][subKey]
		rest := key[len(prefix):]
		bracketEnd := strings.Index(rest, "]")
		if bracketEnd < 0 {
			continue
		}
		indexStr := rest[:bracketEnd]
		index := 0
		fmt.Sscanf(indexStr, "%d", &index)

		subRest := rest[bracketEnd+1:]
		if !strings.HasPrefix(subRest, "[") || !strings.HasSuffix(subRest, "]") {
			continue
		}
		subKey := subRest[1 : len(subRest)-1]

		if rows[index] == nil {
			rows[index] = map[string]string{}
		}
		if len(vals) > 0 {
			rows[index][subKey] = vals[0]
		}
		if index > maxIndex {
			maxIndex = index
		}
	}

	if maxIndex < 0 {
		return "[]"
	}

	// build ordered slice
	result := make([]map[string]string, maxIndex+1)
	for i := 0; i <= maxIndex; i++ {
		if rows[i] != nil {
			result[i] = rows[i]
		} else {
			result[i] = map[string]string{}
		}
	}

	b, err := json.Marshal(result)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// validateRepeaterFields validates repeater fields in the form data.
func (crud *Crud) validateRepeaterFields(data map[string]string, fields []FieldInterface) error {
	for _, field := range fields {
		if field.GetType() != FORM_FIELD_TYPE_REPEATER || !field.GetRequired() {
			continue
		}
		fieldName := field.GetName()
		if fieldName == "" {
			continue
		}
		fieldValue := strings.TrimSpace(data[fieldName])
		if fieldValue == "" || fieldValue == "[]" {
			return fmt.Errorf("%s is required field", field.GetLabel())
		}
		if strings.HasPrefix(fieldValue, "[") && strings.HasSuffix(fieldValue, "]") {
			continue
		}
		return fmt.Errorf("%s must be a valid JSON array", field.GetLabel())
	}
	return nil
}
