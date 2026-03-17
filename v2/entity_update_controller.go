package crud

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/samber/lo"
)

type entityUpdateController struct {
	crud *Crud
}

func (crud *Crud) newEntityUpdateController() *entityUpdateController {
	return &entityUpdateController{
		crud: crud,
	}
}

func (controller *entityUpdateController) page(w http.ResponseWriter, r *http.Request) {
	entityID := req.GetStringTrimmed(r, "entity_id")
	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	if controller.crud.funcFetchUpdateData == nil {
		api.Respond(w, r, api.Error("Update functionality is not configured"))
		return
	}

	breadcrumbs := controller.crud.renderBreadcrumbs([]Breadcrumb{
		{
			Name: "Home",
			URL:  controller.crud.urlHome(),
		},
		{
			Name: controller.crud.entityNameSingular + " Manager",
			URL:  controller.crud.UrlEntityManager(),
		},
		{
			Name: "Edit " + controller.crud.entityNameSingular,
			URL:  controller.crud.UrlEntityUpdate() + "&entity_id=" + entityID,
		},
	})

	buttonSave := hb.Button().Class("btn btn-success float-end").Attr("v-on:click", "entitySave(true)").
		AddChild(hb.I().Class("bi-check-all").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Save")
	buttonApply := hb.Button().Class("btn btn-success float-end").Attr("v-on:click", "entitySave").
		Style("margin-right:10px;").
		AddChild(hb.I().Class("bi-check").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Apply")
	heading := hb.Heading1().Text("Edit " + controller.crud.entityNameSingular).
		AddChild(buttonSave).
		AddChild(buttonApply)

	// container.AddChild(hb.HTML(header))
	container := hb.Div().Attr("class", "container").Attr("id", "entity-update").
		AddChild(heading).
		AddChild(hb.Raw(breadcrumbs))

	customAttrValues, errData := controller.crud.funcFetchUpdateData(r, entityID)

	if errData != nil {
		api.Respond(w, r, api.Error("Fetch data failed"))
		return
	}

	formFields := controller.crud.form(controller.crud.updateFields)

	// Add bottom buttons container
	buttonCancel := hb.Button().Class("btn btn-secondary").Attr("onclick", "window.location.href='"+controller.crud.UrlEntityManager()+"'").
		AddChild(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel")

	bottomButtons := hb.Div().
		Class("row mt-4 mb-4").
		Child(hb.Div().
			Class("col-12").
			Child(hb.Div().
				Class("d-flex justify-content-between").
				Child(buttonCancel).
				Child(hb.Div().Class("d-flex gap-2").
					Child(buttonApply).
					Child(buttonSave))))

	container.AddChildren(formFields)
	container.AddChild(bottomButtons)

	content := container.ToHTML()

	jsonCustomValues, err := json.Marshal(customAttrValues)
	if err != nil {
		api.Respond(w, r, api.Error("Failed to encode form data"))
		return
	}

	saveRedirectURL := controller.crud.updateRedirectURL
	if saveRedirectURL == "" {
		saveRedirectURL = controller.crud.endpoint
	}
	urlHome, err := json.Marshal(saveRedirectURL)
	if err != nil {
		api.Respond(w, r, api.Error("Failed to encode URL"))
		return
	}
	urlEntityTrashAjax, err := json.Marshal(controller.crud.UrlEntityTrashAjax())
	if err != nil {
		api.Respond(w, r, api.Error("Failed to encode URL"))
		return
	}
	urlEntityUpdateAjax, err := json.Marshal(controller.crud.UrlEntityUpdateAjax())
	if err != nil {
		api.Respond(w, r, api.Error("Failed to encode URL"))
		return
	}

	entityIDJson, err := json.Marshal(entityID)
	if err != nil {
		api.Respond(w, r, api.Error("Failed to encode entity ID"))
		return
	}

	inlineScript := `
	const entityManagerUrl = ` + string(urlHome) + `;
	const entityUpdateUrl = ` + string(urlEntityUpdateAjax) + `;
	const entityTrashUrl = ` + string(urlEntityTrashAjax) + `;
	const entityId = ` + string(entityIDJson) + `;
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
						Swal.fire({
							icon:'success', 
							title:'Saved successfully',
							position:'top-left',
							showConfirmButton:false,
							timer:3000,
							timerProgressBar:true
						}).then(() => {
							window.location.href=entityManagerUrl;
						});
					} else {
						return Swal.fire({icon: 'success',title: 'Entity saved'});
					}
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

	title := "Edit " + controller.crud.entityNameSingular
	html := controller.crud.layout(w, r, title, content, []string{
		cdn.JqueryDataTablesCss_1_13_4(),
		cdn.TrumbowygCss_2_27_3(),
	}, "", []string{
		cdn.JqueryDataTablesJs_1_13_4(),
		cdn.TrumbowygJs_2_27_3(),
		"https://cdn.jsdelivr.net/npm/vue-trumbowyg@4",
		"https://cdn.jsdelivr.net/npm/element-plus",
	}, inlineScript)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(html))
}

func (controller *entityUpdateController) pageSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}

	entityID := req.GetStringTrimmed(r, "entity_id")

	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	names := controller.crud.listUpdateNames()
	posts := map[string]string{}
	for _, name := range names {
		posts[name] = req.GetStringTrimmed(r, name)
	}

	// Check required fields
	for _, field := range controller.crud.updateFields {
		if !field.GetRequired() {
			continue
		}

		if _, exists := posts[field.GetName()]; !exists {
			api.Respond(w, r, api.Error(field.GetLabel()+" is required field"))
			return
		}

		if lo.IsEmpty(posts[field.GetName()]) {
			api.Respond(w, r, api.Error(field.GetLabel()+" is required field"))
			return
		}
	}

	err := controller.crud.funcUpdate(r, entityID, posts)

	if err != nil {
		api.Respond(w, r, api.Error("Save failed: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Saved successfully", map[string]interface{}{"entity_id": entityID}))
}
