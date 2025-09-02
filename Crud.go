package crud

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/str"
	"github.com/samber/lo"
)

type Crud struct {
	columnNames         []string
	createFields        []FormField
	endpoint            string
	entityNamePlural    string
	entityNameSingular  string
	fileManagerURL      string
	funcCreate          func(data map[string]string) (userID string, err error)
	funcReadExtras      func(entityID string) []hb.TagInterface
	funcFetchReadData   func(entityID string) ([][2]string, error)
	funcFetchUpdateData func(entityID string) (map[string]string, error)
	funcLayout          func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string
	funcRows            func() (rows []Row, err error)
	funcTrash           func(entityID string) error
	funcUpdate          func(entityID string, data map[string]string) error
	homeURL             string
	readFields          []FormField
	updateFields        []FormField
}

func (crud Crud) Handler(w http.ResponseWriter, r *http.Request) {
	path := req.GetStringTrimmed(r, "path")

	if path == "" {
		path = "home"
	}

	ctx := context.WithValue(r.Context(), "", r.URL.Path)

	routeFunc := crud.getRoute(path)
	routeFunc(w, r.WithContext(ctx))
}

func (crud *Crud) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		"home": crud.pageEntityManager,
		// START: Custom Entities
		pathEntityCreateAjax: crud.pageEntityCreateAjax,
		pathEntityManager:    crud.pageEntityManager,
		pathEntityRead:       crud.pageEntityRead,
		pathEntityUpdate:     crud.pageEntityUpdate,
		pathEntityUpdateAjax: crud.pageEntityUpdateAjax,
		pathEntityTrashAjax:  crud.pageEntityTrashAjax,
		// END: Custom Entities

	}
	// log.Println(route)
	if val, ok := routes[route]; ok {
		return val
	}

	return routes["home"]
}

func (crud *Crud) pageEntityCreateAjax(w http.ResponseWriter, r *http.Request) {
	names := crud.listCreateNames()

	posts := map[string]string{}
	for _, name := range names {
		posts[name] = req.GetStringTrimmed(r, name)
	}

	// Check required fields
	for _, field := range crud.createFields {
		if !field.Required {
			continue
		}

		if _, exists := posts[field.Name]; !exists {
			api.Respond(w, r, api.Error(field.Label+" is required field"))
			return
		}

		if lo.IsEmpty(posts[field.Name]) {
			api.Respond(w, r, api.Error(field.Label+" is required field"))
			return
		}
	}

	entityID, err := crud.funcCreate(posts)

	if err != nil {
		api.Respond(w, r, api.Error("Save failed: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Saved successfully", map[string]interface{}{"entity_id": entityID}))
}

func (crud *Crud) pageEntityManager(w http.ResponseWriter, r *http.Request) {
	// header := cms.cmsHeader(endpoint)
	breadcrumbs := crud._breadcrumbs([]Breadcrumb{
		{
			Name: "Home",
			URL:  crud.urlHome(),
		},
		{
			Name: crud.entityNameSingular + " Manager",
			URL:  crud.UrlEntityManager(),
		},
	})

	buttonCreate := hb.Button().
		Class("btn btn-success float-end").
		Attr("v-on:click", "showEntityCreateModal").
		AddChild(hb.I().Class("bi-plus-circle").Style("margin-top:-4px;margin-right:8px;")).
		HTML("New " + crud.entityNameSingular)

	heading := hb.Heading1().
		HTML(crud.entityNameSingular + " Manager").
		Child(buttonCreate)

	rows, errRows := crud.funcRows()

	tableContent := lo.IfF(errRows != nil, func() hb.TagInterface {
		alert := hb.Div().
			Class("alert alert-danger").
			HTML("There was an error retrieving the data. Please try again later")

		return alert
	}).ElseF(func() hb.TagInterface {
		table := hb.Table().
			ID("TableEntities").
			Class("table table-responsive table-striped mt-3").
			Child(
				hb.Thead().
					Children([]hb.TagInterface{
						hb.TR().
							Children(lo.Map(crud.columnNames, func(columnName string, _ int) hb.TagInterface {
								columnName = strings.ReplaceAll(columnName, "{!!", "")
								columnName = strings.ReplaceAll(columnName, "!!}", "")
								return hb.TH().Text(columnName)
							})).
							Child(hb.TD().
								HTML("Actions").
								Style("width:120px;")),
					})).
			Child(
				hb.Tbody().
					Children(lo.Map(rows, func(row Row, _ int) hb.TagInterface {
						buttonView := hb.Hyperlink().
							Class("btn btn-sm btn-outline-info").
							Child(hb.I().Class("bi-eye").Style("margin-top:-4px;")).
							Attr("title", "Show").
							Href(crud.UrlEntityRead() + "&entity_id=" + row.ID).
							Style("margin-right:5px")

						buttonEdit := hb.Hyperlink().
							Class("btn btn-sm btn-outline-warning").
							Child(hb.I().Class("bi-pencil-square").Style("margin-top:-4px;")).
							Attr("title", "Edit").
							Attr("type", "button").
							Href(crud.UrlEntityUpdate() + "&entity_id=" + row.ID).
							Style("margin-right:5px")

						buttonTrash := hb.Button().
							Class("btn btn-sm btn-outline-danger").
							Child(hb.I().Class("bi-trash").Style("margin-top:-4px;")).
							Attr("title", "Trash").
							Attr("type", "button").
							Attr("v-on:click", "showEntityTrashModal('"+row.ID+"')")

						tr := hb.TR().
							Children(lo.Map(row.Data, func(cell string, index int) hb.TagInterface {
								name := crud.columnNames[index]
								isRaw := strings.HasPrefix(name, "{!!") && strings.HasSuffix(name, "!!}")
								cell = strings.ReplaceAll(cell, "{!!", "")
								cell = strings.ReplaceAll(cell, "!!}", "")
								cell = strings.TrimSpace(cell)
								return hb.TD().TextIf(!isRaw, cell).HTMLIf(isRaw, cell)
							})).
							Child(
								hb.TD().
									Style(`white-space:nowrap;`).
									ChildIf(crud.funcFetchReadData != nil, buttonView).
									ChildIf(crud.funcFetchUpdateData != nil, buttonEdit).
									ChildIf(crud.funcTrash != nil, buttonTrash),
							)
						return tr
					})))

		return table
	})

	container := hb.Div().
		ID("entity-manager").
		Class("container").
		Child(heading).
		Child(hb.Raw(breadcrumbs)).
		Child(crud.pageEntitiesEntityCreateModal()).
		Child(crud.pageEntitiesEntityTrashModal()).
		Child(tableContent)

	content := container.ToHTML()

	urlEntityCreateAjax, _ := json.Marshal(crud.UrlEntityCreateAjax())
	urlEntityTrashAjax, _ := json.Marshal(crud.UrlEntityTrashAjax())
	urlEntityUpdate, _ := json.Marshal(crud.UrlEntityUpdate())

	customAttrValues := map[string]string{}
	lo.ForEach(crud.createFields, func(field FormField, index int) {
		customAttrValues[field.Name] = field.Value
	})
	jsonCustomValues, _ := json.Marshal(customAttrValues)

	inlineScript := `
const entityCreateUrl = ` + string(urlEntityCreateAjax) + `;
const entityUpdateUrl = ` + string(urlEntityUpdate) + `;
const entityTrashUrl = ` + string(urlEntityTrashAjax) + `;
const customValues = ` + string(jsonCustomValues) + `;
const EntityManager = {
	data() {
		return {
		  entityModel:{
			...customValues
		  },
		  entityTrashModel:{
			entityId:null,
		  }
		}
	},
	created(){
		//setTimeout(() => {
		//	console.log("Init data table...");
			this.initDataTable();
		//}, 1000);
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableEntities').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
        showEntityCreateModal(){
			const modalEntityCreate = new bootstrap.Modal(document.getElementById('ModalEntityCreate'));
			modalEntityCreate.show();
		},
		showEntityTrashModal(entityId){
			this.entityTrashModel.entityId = entityId;
			const modalEntityDelete = new bootstrap.Modal(document.getElementById('ModalEntityTrash'));
			modalEntityDelete.show();
		},
		entityCreate(){
		    $.post(entityCreateUrl, this.entityModel).done((result)=>{
				if (result.status==="success"){
					const modalEntityCreate = new bootstrap.Modal(document.getElementById('ModalEntityCreate'));
			        modalEntityCreate.hide();
					return location.href = entityUpdateUrl+ "&entity_id=" + result.data.entity_id;
				}
				
				return Swal.fire({icon: 'error', title: 'Oops...', text: result.message});
			}).fail((result)=>{
				return Swal.fire({icon: 'error', title: 'Oops...', text: result});
			});
		},

		entityTrash(){
			const entityId = this.entityTrashModel.entityId;

			$.post(entityTrashUrl, {
				entity_id:entityId
			}).done((response)=>{
				if (response.status !== "success") {
					return Swal.fire({icon: 'error', title: 'Oops...', text: result.message});
				}

				setTimeout(()=>{return location.href = location.href;}, 3000)

				return Swal.fire({icon: 'success', title: 'Entity trashed'});
			}).fail((result)=>{
				console.log(result);
				return Swal.fire({icon: 'error', title: 'Oops...', text: result});
			});
		}
	}
};
Vue.createApp(EntityManager).mount('#entity-manager')
	`
	title := crud.entityNameSingular + " Manager"
	html := crud.layout(w, r, title, content, []string{
		cdn.JqueryDataTablesCss_1_13_4(),
	}, "html{width:100%;}", []string{
		cdn.JqueryDataTablesJs_1_13_4(),
	}, inlineScript)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (crud *Crud) pageEntityRead(w http.ResponseWriter, r *http.Request) {
	entityID := req.GetStringTrimmed(r, "entity_id")
	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	if crud.funcFetchReadData == nil {
		api.Respond(w, r, api.Error("FuncFetchReadData is required"))
		return
	}

	breadcrumbs := crud._breadcrumbs([]Breadcrumb{
		{
			Name: "Home",
			URL:  crud.urlHome(),
		},
		{
			Name: crud.entityNameSingular + " Manager",
			URL:  crud.UrlEntityManager(),
		},
		{
			Name: "View " + crud.entityNameSingular,
			URL:  crud.UrlEntityUpdate() + "&entity_id=" + entityID,
		},
	})

	buttonEdit := hb.Hyperlink().
		Class("btn btn-primary ml-2 float-end").
		Child(hb.I().Class("bi-pencil-square").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Edit").
		Href(crud.UrlEntityUpdate() + "&entity_id=" + entityID)

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ml-2 float-end").
		Child(hb.I().Class("bi-chevron-left").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Back").
		Href(crud.UrlEntityManager())

	heading := hb.Heading1().
		HTML("View " + crud.entityNameSingular).
		Child(buttonEdit).
		Child(buttonCancel)

	container := hb.Div().
		ID("entity-read").
		Class("container").
		Child(heading).
		Child(hb.Raw(breadcrumbs))

	data, err := crud.funcFetchReadData(entityID)

	table := lo.IfF(err != nil, func() hb.TagInterface {
		alert := hb.Div().
			Class("alert alert-danger").
			HTML("There was an error retrieving the data. Please try again later")

		return alert
	}).ElseF(func() hb.TagInterface {
		table := hb.Table().
			Class("table table-hover table-striped").
			Child(hb.Thead().Child(hb.TR())).
			Child(hb.Tbody().Children(lo.Map(data, func(row [2]string, _ int) hb.TagInterface {
				key := row[0]
				value := row[1]
				isRawKey := strings.HasPrefix(key, "{!!") && strings.HasSuffix(key, "!!}")
				isRawValue := strings.HasPrefix(value, "{!!") && strings.HasSuffix(value, "!!}")

				key = strings.ReplaceAll(key, "{!!", "")
				key = strings.ReplaceAll(key, "!!}", "")
				key = strings.TrimSpace(key)

				value = strings.ReplaceAll(value, "{!!", "")
				value = strings.ReplaceAll(value, "!!}", "")
				value = strings.TrimSpace(value)

				return hb.TR().Children([]hb.TagInterface{
					hb.TH().TextIf(!isRawKey, key).HTMLIf(isRawKey, key),
					hb.TD().TextIf(!isRawValue, value).HTMLIf(isRawValue, value),
				})
			})))

		return table
	})

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTML(crud.entityNameSingular + " Details").
					Style("margin-bottom:0;display:inline-block;")).
				Child(buttonEdit),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(table))

	container.Child(card)
	if crud.funcReadExtras != nil {
		container.Children(crud.funcReadExtras(entityID))
	}
	content := container.ToHTML()
	title := "View " + crud.entityNameSingular
	html := crud.layout(w, r, title, content, []string{}, "", []string{}, "")

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (crud *Crud) pageEntityUpdate(w http.ResponseWriter, r *http.Request) {
	entityID := req.GetStringTrimmed(r, "entity_id")
	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	breadcrumbs := crud._breadcrumbs([]Breadcrumb{
		{
			Name: "Home",
			URL:  crud.urlHome(),
		},
		{
			Name: crud.entityNameSingular + " Manager",
			URL:  crud.UrlEntityManager(),
		},
		{
			Name: "Edit " + crud.entityNameSingular,
			URL:  crud.UrlEntityUpdate() + "&entity_id=" + entityID,
		},
	})

	buttonSave := hb.Button().Class("btn btn-success float-end").Attr("v-on:click", "entitySave(true)").
		AddChild(hb.I().Class("bi-check-all").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Save")
	buttonApply := hb.Button().Class("btn btn-success float-end").Attr("v-on:click", "entitySave").
		Style("margin-right:10px;").
		AddChild(hb.I().Class("bi-check").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Apply")
	heading := hb.Heading1().Text("Edit " + crud.entityNameSingular).
		AddChild(buttonSave).
		AddChild(buttonApply)

	// container.AddChild(hb.HTML(header))
	container := hb.Div().Attr("class", "container").Attr("id", "entity-update").
		AddChild(heading).
		AddChild(hb.Raw(breadcrumbs))

	customAttrValues, errData := crud.funcFetchUpdateData(entityID)

	if errData != nil {
		api.Respond(w, r, api.Error("Fetch data failed"))
		return
	}

	container.AddChildren(crud.form(crud.updateFields))

	content := container.ToHTML()

	jsonCustomValues, _ := json.Marshal(customAttrValues)

	urlHome, _ := json.Marshal(crud.endpoint)
	urlEntityTrashAjax, _ := json.Marshal(crud.UrlEntityTrashAjax())
	urlEntityUpdateAjax, _ := json.Marshal(crud.UrlEntityUpdateAjax())

	inlineScript := `
	const entityManagerUrl = ` + string(urlHome) + `;
	const entityUpdateUrl = ` + string(urlEntityUpdateAjax) + `;
	const entityTrashUrl = ` + string(urlEntityTrashAjax) + `;
	const entityId = "` + entityID + `";
	const customValues = ` + string(jsonCustomValues) + `;
	const EntityUpdate = {
		data() {
			return {
				entityModel:{
					entityId,
					...customValues
			    },
				tmp:{},
				trumbowigConfig: {
					btns: [
						['undo', 'redo'], 
						['formatting'], 
						['strong', 'em', 'del', 'superscript', 'subscript'], 
						['link','justifyLeft','justifyRight','justifyCenter','justifyFull'], 
						['unorderedList', 'orderedList'], 
						['horizontalRule'], 
						['removeformat'], 
						['fullscreen']
					],	
					autogrow: true,
					removeformatPasted: true,
					tagsToRemove: ['script', 'link', 'embed', 'iframe', 'input'],
					tagsToKeep: ['hr', 'img', 'i'],
					autogrowOnEnter: true,
					linkTargets: ['_blank'],
				},
			}
		},
		methods: {
			entitySave(redirect){
				const entityId = this.entityModel.entityId;
				let data = JSON.parse(JSON.stringify(this.entityModel));
				data["entity_id"] = data["entityId"];
				delete data["entityId"];

				$.post(entityUpdateUrl, data).done((response)=>{
					if (response.status !== "success") {
						return Swal.fire({icon: 'error', title: 'Oops...', text: response.message});
					}

					if (redirect===true) {
						setTimeout(()=>{
							window.location.href=entityManagerUrl;
						}, 3000)
					}

					return Swal.fire({icon: 'success',title: 'Entity saved'});
				}).fail((result)=>{
					console.log(result);
					return Swal.fire({icon: 'error', title: 'Oops...', text: result});
				});
			},
			uploadImage(event, fieldName) {
				const self = this;
				if ( event.target.files && event.target.files[0] ) {
					var FR= new FileReader();
					FR.onload = function(e) {
						self.entityModel[fieldName] = e.target.result;
						event.target.value = "";
					};       
					FR.readAsDataURL( event.target.files[0] );
				}
			}
		}
	};
	Vue.createApp(EntityUpdate).use(ElementPlus).component('Trumbowyg', VueTrumbowyg.default).mount('#entity-update')
		`

	// webpage := crud.webpage("Edit "+crud.entityNameSingular, h)
	// webpage.AddScript(inlineScript)

	title := "Edit " + crud.entityNameSingular
	html := crud.layout(w, r, title, content, []string{
		cdn.JqueryDataTablesCss_1_13_4(),
		cdn.TrumbowygCss_2_27_3(),
	}, "", []string{
		cdn.JqueryDataTablesJs_1_13_4(),
		cdn.TrumbowygJs_2_27_3(),
		"https://cdn.jsdelivr.net/npm/vue-trumbowyg@4",
		"https://cdn.jsdelivr.net/npm/element-plus",
	}, inlineScript)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (crud *Crud) pageEntityUpdateAjax(w http.ResponseWriter, r *http.Request) {
	entityID := strings.Trim(req.GetStringTrimmed(r, "entity_id"), " ")

	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	names := crud.listUpdateNames()
	posts := map[string]string{}
	for _, name := range names {
		posts[name] = req.GetStringTrimmed(r, name)
	}

	// Check required fields
	for _, field := range crud.updateFields {
		if !field.Required {
			continue
		}

		if _, exists := posts[field.Name]; !exists {
			api.Respond(w, r, api.Error(field.Label+" is required field"))
			return
		}

		if lo.IsEmpty(posts[field.Name]) {
			api.Respond(w, r, api.Error(field.Label+" is required field"))
			return
		}
	}

	err := crud.funcUpdate(entityID, posts)

	if err != nil {
		api.Respond(w, r, api.Error("Save failed: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Saved successfully", map[string]interface{}{"entity_id": entityID}))
}

func (crud *Crud) pageEntityTrashAjax(w http.ResponseWriter, r *http.Request) {
	entityID := strings.Trim(req.GetStringTrimmed(r, "entity_id"), " ")

	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	err := crud.funcTrash(entityID)

	if err != nil {
		api.Respond(w, r, api.Error("Entity failed to be trashed: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Entity trashed successfully", map[string]interface{}{"entity_id": entityID}))
}

func (crud *Crud) pageEntitiesEntityTrashModal() hb.TagInterface {
	modal := hb.Div().ID("ModalEntityTrash").Class("modal fade")
	modalDialog := hb.Div().Attr("class", "modal-dialog")
	modalContent := hb.Div().Attr("class", "modal-content")
	modalHeader := hb.Div().Attr("class", "modal-header").AddChild(hb.Heading5().Text("Trash Entity"))
	modalBody := hb.Div().Attr("class", "modal-body")
	modalBody.AddChild(hb.Paragraph().Text("Are you sure you want to move this entity to trash bin?"))
	modalFooter := hb.Div().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.Button().Text("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.Button().Text("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "entityTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func (crud *Crud) pageEntitiesEntityCreateModal() hb.TagInterface {
	form := crud.form(crud.createFields)

	modalHeader := hb.Div().Class("modal-header").
		AddChild(hb.Heading5().Text("New " + crud.entityNameSingular))

	modalBody := hb.Div().Class("modal-body").AddChildren(form)

	modalFooter := hb.Div().Class("modal-footer").
		AddChild(hb.Button().Text("Close").Class("btn btn-secondary").Attr("data-bs-dismiss", "modal")).
		AddChild(hb.Button().Text("Create & Continue").Class("btn btn-primary").Attr("v-on:click", "entityCreate"))

	modal := hb.Div().ID("ModalEntityCreate").Class("modal fade").AddChildren([]hb.TagInterface{
		hb.Div().Class("modal-dialog").AddChildren([]hb.TagInterface{
			hb.Div().Class("modal-content").AddChildren([]hb.TagInterface{
				modalHeader,
				modalBody,
				modalFooter,
			}),
		}),
	})

	return modal
}

func (crud *Crud) urlHome() string {
	url := crud.homeURL
	return url
}

func (crud *Crud) UrlEntityManager() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	url := crud.endpoint + q + "path=" + pathEntityManager
	return url
}

func (crud *Crud) UrlEntityCreateAjax() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	url := crud.endpoint + q + "path=" + pathEntityCreateAjax
	return url
}

func (crud *Crud) UrlEntityTrashAjax() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	url := crud.endpoint + q + "path=" + pathEntityTrashAjax
	return url
}

func (crud *Crud) UrlEntityRead() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	url := crud.endpoint + q + "path=" + pathEntityRead
	return url
}

func (crud *Crud) UrlEntityUpdate() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	url := crud.endpoint + q + "path=" + pathEntityUpdate
	return url
}

func (crud *Crud) UrlEntityUpdateAjax() string {
	q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
	url := crud.endpoint + q + "path=" + pathEntityUpdateAjax
	return url
}

// Webpage returns the webpage template for the website
func (crud *Crud) webpage(title, content string) *hb.HtmlWebpage {
	faviconImgCms := `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`
	app := ""
	webpage := hb.Webpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(faviconImgCms)

	webpage.AddStyleURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/css/bootstrap.min.css",
	})
	webpage.AddScriptURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/js/bootstrap.bundle.min.js",
		"https://code.jquery.com/jquery-3.6.0.min.js",
		"https://unpkg.com/vue@3/dist/vue.global.js",
		"https://cdn.jsdelivr.net/npm/sweetalert2@9",
	})
	webpage.AddScripts([]string{
		app,
	})
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
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
	webpage.Child(hb.Raw(content))
	return webpage
}

func (crud *Crud) _breadcrumbs(breadcrumbs []Breadcrumb) string {
	nav := hb.Nav().Attr("aria-label", "breadcrumb")
	ol := hb.OL().Attr("class", "breadcrumb")

	for _, breadcrumb := range breadcrumbs {
		li := hb.LI().Attr("class", "breadcrumb-item")
		link := hb.Hyperlink().Text(breadcrumb.Name).Attr("href", breadcrumb.URL)

		li.AddChild(link)

		ol.AddChild(li)
	}

	nav.AddChild(ol)

	return nav.ToHTML()
}

// layout is a function that generates an HTML layout for a web page.
//
// Parameters:
// - w: an http.ResponseWriter object for writing the HTTP response.
// - r: a pointer to an http.Request object representing the HTTP request.
// - title: a string containing the title of the web page.
// - content: a string containing the content of the web page.
// - styleFiles: a slice of strings representing the URLs of the style files to be included in the web page.
// - style: a string containing the CSS style to be applied to the web page.
// - jsFiles: a slice of strings representing the URLs of the JavaScript files to be included in the web page.
// - js: a string containing the JavaScript code to be executed in the web page.
//
// Returns:
// - string - a string representing the generated HTML layout.
func (crud *Crud) layout(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string {
	html := ""

	if crud.funcLayout != nil {
		// jsFiles = append([]string{"//unpkg.com/naive-ui"}, jsFiles...)
		jsFiles = append([]string{"//cdn.jsdelivr.net/npm/element-plus"}, jsFiles...)
		jsFiles = append([]string{"https://unpkg.com/vue@3/dist/vue.global.js"}, jsFiles...)
		jsFiles = append([]string{"https://cdn.jsdelivr.net/npm/sweetalert2@9"}, jsFiles...)
		styleFiles = append([]string{"https://unpkg.com/element-plus/dist/index.css"}, styleFiles...)
		html = crud.funcLayout(w, r, title, content, styleFiles, style, jsFiles, js)
	} else {
		webpage := crud.webpage(title, content)
		webpage.AddStyleURLs(styleFiles)
		webpage.AddStyle(style)
		webpage.AddScriptURLs(jsFiles)
		webpage.AddScript(js)
		html = webpage.ToHTML()
	}

	return html
}

// form generates a form with entries for each form field.
//
// Parameters:
// - fields: a slice of FormField structs representing the fields in the form.
//
// Returns:
// - a slice of hb.Tags representing the form.
func (crud *Crud) form(fields []FormField) []hb.TagInterface {
	tags := []hb.TagInterface{}
	for _, field := range fields {
		fieldID := field.ID
		if fieldID == "" {
			fieldID = "id_" + str.RandomFromGamma(32, "abcdefghijklmnopqrstuvwxyz1234567890")
		}
		fieldName := field.Name
		fieldValue := field.Value
		fieldLabel := field.Label
		if fieldLabel == "" {
			fieldLabel = fieldName
		}

		formGroup := hb.Div().Class("form-group mt-3")

		formGroupLabel := hb.Label().
			Text(fieldLabel).
			Class("form-label").
			ChildIf(
				field.Required,
				hb.Sup().Text("*").Class("text-danger ml-1"),
			)

		formGroupInput := hb.Input().
			Class("form-control").
			Attr("v-model", "entityModel."+fieldName)

		if field.Type == FORM_FIELD_TYPE_IMAGE {
			formGroupInput = hb.Div().Children([]hb.TagInterface{
				hb.Image("").
					Attr(`v-bind:src`, `entityModel.`+fieldName+`||'https://www.freeiconspng.com/uploads/no-image-icon-11.PNG'`).
					Style(`width:200px;`),
				bs.InputGroup().Children([]hb.TagInterface{
					hb.Input().Type(hb.TYPE_URL).Class("form-control").Attr("v-model", "entityModel."+fieldName),
					hb.If(crud.fileManagerURL != "", bs.InputGroupText().Children([]hb.TagInterface{
						hb.Hyperlink().Text("Browse").Href(crud.fileManagerURL).Target("_blank"),
					})),
				}),
			})
		}

		if field.Type == FORM_FIELD_TYPE_IMAGE_INLINE {
			formGroupInput = hb.Div().
				Children([]hb.TagInterface{
					hb.Image("").
						Attr(`v-bind:src`, `entityModel.`+fieldName+`||'https://www.freeiconspng.com/uploads/no-image-icon-11.PNG'`).
						Style(`width:200px;`),
					hb.Input().
						Type(hb.TYPE_FILE).
						Attr("v-on:change", "uploadImage($event, '"+fieldName+"')").
						Attr("accept", "image/*"),
					hb.Button().
						HTML("See Image Data").
						Attr("v-on:click", "tmp.show_url_"+fieldName+" = !tmp.show_url_"+fieldName),
					hb.TextArea().
						Type(hb.TYPE_URL).
						Class("form-control").
						Attr("v-if", "tmp.show_url_"+fieldName).
						Attr("v-model", "entityModel."+fieldName),
				})
		}

		if field.Type == FORM_FIELD_TYPE_DATETIME {
			// formGroupInput = hb.Input().Type(hb.TYPE_DATETIME).Class("form-control").Attr("v-model", "entityModel."+fieldName)
			formGroupInput = hb.NewTag(`el-date-picker`).Attr("type", "datetime").Attr("v-model", "entityModel."+fieldName)
			// formGroupInput = hb.Tag(`n-date-picker`).Attr("type", "datetime").Class("form-control").Attr("v-model", "entityModel."+fieldName)
		}

		if field.Type == FORM_FIELD_TYPE_HTMLAREA {
			formGroupInput = hb.NewTag("trumbowyg").Attr("v-model", "entityModel."+fieldName).Attr(":config", "trumbowigConfig").Class("form-control")
		}

		if field.Type == FORM_FIELD_TYPE_NUMBER {
			formGroupInput.Type(hb.TYPE_NUMBER)
		}

		if field.Type == FORM_FIELD_TYPE_PASSWORD {
			formGroupInput.Type(hb.TYPE_PASSWORD)
		}

		if field.Type == FORM_FIELD_TYPE_SELECT {
			formGroupInput = hb.Select().Class("form-select").Attr("v-model", "entityModel."+fieldName)
			for _, opt := range field.Options {
				option := hb.Option().Value(opt.Key).Text(opt.Value)
				formGroupInput.AddChild(option)
			}
			if field.OptionsF != nil {
				for _, opt := range field.OptionsF() {
					option := hb.Option().Value(opt.Key).Text(opt.Value)
					formGroupInput.AddChild(option)
				}
			}
		}

		if field.Type == FORM_FIELD_TYPE_TEXTAREA {
			formGroupInput = hb.TextArea().Class("form-control").Attr("v-model", "entityModel."+fieldName)
		}

		if field.Type == FORM_FIELD_TYPE_BLOCKAREA {
			formGroupInput = hb.TextArea().Class("form-control").Attr("v-model", "entityModel."+fieldName)
		}

		if field.Type == FORM_FIELD_TYPE_RAW {
			formGroupInput = hb.Raw(fieldValue)
		}

		formGroupInput.ID(fieldID)
		if field.Type != FORM_FIELD_TYPE_RAW {
			formGroup.AddChild(formGroupLabel)
		}
		formGroup.AddChild(formGroupInput)

		// Add help
		if field.Help != "" {
			formGroupHelp := hb.Paragraph().Class("text-info").HTML(field.Help)
			formGroup.AddChild(formGroupHelp)
		}

		tags = append(tags, formGroup)

		if field.Type == FORM_FIELD_TYPE_BLOCKAREA {
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

// listCreateNames returns a list of names from the createFields
// slice in the Crud struct.
//
// Parameters:
//   - None
//
// Returns:
//   - []string - a list of field names
func (crud *Crud) listCreateNames() []string {
	names := []string{}

	for _, field := range crud.createFields {
		if field.Name == "" {
			continue
		}
		names = append(names, field.Name)
	}

	return names
}

// listUpdateNames returns a list of names from the updateFields
// slice in the Crud struct.
//
// Parameters:
//   - None
//
// Returns:
//   - []string - a list of field names
func (crud *Crud) listUpdateNames() []string {
	names := []string{}

	for _, field := range crud.updateFields {
		if field.Name == "" {
			continue
		}
		names = append(names, field.Name)
	}

	return names
}
