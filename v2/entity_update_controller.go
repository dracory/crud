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
	buttonApply := hb.Button().Class("btn btn-success float-end").Attr("v-on:click", "entitySave(false)").
		Style("margin-right:10px;").
		AddChild(hb.I().Class("bi-check").Style("margin-top:-4px;margin-right:8px;")).
		HTML("Apply")
	heading := hb.Heading1().Text("Edit " + controller.crud.entityNameSingular).
		AddChild(buttonSave).
		AddChild(buttonApply)

	// container.AddChild(hb.HTML(header))
	container := hb.Div().Attr("class", "container mt-3").Attr("id", "entity-update").
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

	// Build a map[string]any so repeater fields can contribute structured values
	// ([]map[string]string) while all other fields remain plain strings.
	processedValues := make(map[string]any, len(customAttrValues))
	for k, v := range customAttrValues {
		processedValues[k] = v
	}
	for _, field := range controller.crud.updateFields {
		fn := field.GetFuncValuesProcess()
		if fn == nil {
			continue
		}
		name := field.GetName()
		if name == "" {
			continue
		}
		raw, _ := customAttrValues[name]
		processedValues[name] = fn(raw)
	}

	jsonCustomValues, err := json.Marshal(processedValues)
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

				const formData = new URLSearchParams();
				for (const key in data) {
					const val = data[key];
					if (Array.isArray(val)) {
						val.forEach((row, i) => {
							if (row !== null && typeof row === 'object') {
								Object.keys(row).forEach(subKey => {
									formData.append(key + '[' + i + '][' + subKey + ']', row[subKey] ?? '');
								});
							} else {
								formData.append(key + '[' + i + ']', row ?? '');
							}
						});
					} else {
						formData.append(key, val ?? '');
					}
				}

				fetch(entityUpdateUrl, {
					method: 'POST',
					headers: {'Content-Type': 'application/x-www-form-urlencoded'},
					body: formData.toString()
				})
				.then(r => r.json())
				.then((response) => {
					if (response.status !== "success") {
						return Swal.fire({icon: 'error', title: 'Oops...', text: response.message});
					}

					if (redirect===true) {
						Swal.fire({
							icon:'success', 
							title:'Saved successfully',
							position:'top-end',
							showConfirmButton:false,
							timer:3000,
							timerProgressBar:true
						}).then(() => {
							window.location.href=entityManagerUrl;
						});
					} else {
						Swal.fire({
							icon:'success',
							title:'Saved successfully',
							position:'top-end',
							showConfirmButton:false,
							timer:2000,
							timerProgressBar:true
						});
					}
				})
				.catch((result)=>{
					console.log(result);
					return Swal.fire({icon: 'error', title: 'Oops...', text: result.toString()});
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
			},
			addRepeaterItem(fieldName, item){
				if (!this.entityModel[fieldName]) {
					this.entityModel[fieldName] = [];
				}
				this.entityModel[fieldName].push(item !== undefined ? item : {});
			},
			removeRepeaterItem(fieldName, index){
				if (this.entityModel[fieldName] && this.entityModel[fieldName].length > 0) {
					this.entityModel[fieldName].splice(index, 1);
				}
			},
			moveRepeaterItemUp(fieldName, index){
				if (index > 0 && this.entityModel[fieldName]) {
					const arr = this.entityModel[fieldName];
					const item = arr.splice(index, 1)[0];
					arr.splice(index - 1, 0, item);
				}
			},
			moveRepeaterItemDown(fieldName, index){
				if (this.entityModel[fieldName] && index < this.entityModel[fieldName].length - 1) {
					const arr = this.entityModel[fieldName];
					const item = arr.splice(index, 1)[0];
					arr.splice(index + 1, 0, item);
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
		cdn.TrumbowygCss_2_27_3(),
	}, "", []string{
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

	posts := map[string]string{}
	for _, field := range controller.crud.updateFields {
		name := field.GetName()
		if name == "" {
			continue
		}
		if field.GetType() == FORM_FIELD_TYPE_REPEATER {
			posts[name] = collectRepeaterField(r, name)
		} else {
			posts[name] = req.GetStringTrimmed(r, name)
		}
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
