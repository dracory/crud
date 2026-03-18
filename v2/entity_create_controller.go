package crud

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/samber/lo"
)

type entityCreateController struct {
	crud *Crud
}

func (crud *Crud) newEntityCreateController() *entityCreateController {
	return &entityCreateController{
		crud: crud,
	}
}

func (controller *entityCreateController) modalShow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(controller.modal().ToHTML()))
}

func (controller *entityCreateController) modalSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}

	if controller.crud.funcCreate == nil {
		api.Respond(w, r, api.Error("Create functionality is not configured"))
		return
	}

	posts := map[string]string{}
	for _, field := range controller.crud.createFields {
		name := field.GetName()
		if name == "" {
			continue
		}
		if field.GetType() == FORM_FIELD_TYPE_REPEATER {
			posts[name] = collectRepeaterField(r, name)
		} else {
			posts[name] = req.GetString(r, name)
		}
	}

	// Check required fields
	for _, field := range controller.crud.createFields {
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

	entityID, err := controller.crud.funcCreate(r, posts)

	if err != nil {
		api.Respond(w, r, api.Error("Save failed: "+err.Error()))
		return
	}

	redirectURL := controller.crud.createRedirectURL
	if redirectURL == "" {
		redirectURL = controller.crud.UrlEntityUpdate() + "&entity_id=" + entityID
	}
	api.Respond(w, r, api.SuccessWithData("Saved successfully", map[string]interface{}{"entity_id": entityID, "redirect_url": redirectURL}))
}

func (controller *entityCreateController) modal() hb.TagInterface {
	formFields := controller.crud.form(controller.crud.createFields)

	submitUrl := controller.crud.UrlEntityCreateAjax()

	modalID := "ModalEntityCreate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := "closeModal" + modalID + "();"

	modalHeading := hb.Heading5().
		Text("New " + controller.crud.entityNameSingular).
		Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := "function closeModal" + modalID + "() {document.getElementById('" + modalID + "Wrapper').remove();[...document.getElementsByClassName('" + modalBackdropClass + "')].forEach(el => el.remove());}"

	jsVueApp := `
const EntityCreate = {
	data() {
		return { entityModel: {} }
	},
	methods: {
		submitModalModalEntityCreate() {
			const data = JSON.parse(JSON.stringify(this.entityModel));
			fetch('` + submitUrl + `', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify(data)
			})
			.then(r => r.json())
			.then(result => {
				if (result.status === 'success') {
					if (result.data && result.data.redirect_url) {
						Swal.fire({icon:'success', title:'Saved successfully', position:'top-left', showConfirmButton:false, timer:3000, timerProgressBar:true})
							.then(() => { window.location.href = result.data.redirect_url; });
					} else {
						Swal.fire({icon:'success', title:'Saved successfully'});
					}
				} else {
					Swal.fire({icon:'error', title:'Oops...', text: result.message});
				}
			})
			.catch(err => { Swal.fire({icon:'error', title:'Oops...', text: err.toString()}); });
		},
		addRepeaterItem(fieldName, item){
			if (!this.entityModel[fieldName]) { this.entityModel[fieldName] = []; }
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
Vue.createApp(EntityCreate).mount('#ModalEntityCreateWrapper');
`

	buttonSubmit := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Create & Edit").
		Class("btn btn-primary float-end").
		Attr("v-on:click", "submitModalModalEntityCreate()")

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Close").
		Class("btn btn-secondary float-start").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	formContainer := hb.Div().Children(formFields)

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(formContainer)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					Child(buttonSubmit)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Div().
		ID(modalID + "Wrapper").
		Children([]hb.TagInterface{
			hb.Script(jsCloseFn),
			modal,
			backdrop,
			hb.Script(jsVueApp),
		})
}
