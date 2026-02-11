package crud

import "errors"

func New(config Config) (crud Crud, err error) {
	if config.FuncRows == nil {
		return Crud{}, errors.New("FuncRows function is required")
	}

	if config.UpdateFields == nil {
		return Crud{}, errors.New("UpdateFields is required")
	}

	if len(config.UpdateFields) > 0 {
		if config.FuncUpdate == nil {
			return Crud{}, errors.New("FuncUpdate function is required when UpdateFields are provided")
		}
		if config.FuncFetchUpdateData == nil {
			return Crud{}, errors.New("FuncFetchUpdateData function is required when UpdateFields are provided")
		}
	}

	crud = Crud{}
	crud.columnNames = config.ColumnNames
	crud.createFields = config.CreateFields
	crud.endpoint = config.Endpoint
	crud.entityNamePlural = config.EntityNamePlural
	crud.entityNameSingular = config.EntityNameSingular
	crud.fileManagerURL = config.FileManagerURL
	crud.funcCreate = config.FuncCreate
	crud.funcReadExtras = config.FuncReadExtras
	crud.funcFetchReadData = config.FuncFetchReadData
	crud.funcFetchUpdateData = config.FuncFetchUpdateData
	crud.funcLayout = config.FuncLayout
	crud.funcRows = config.FuncRows
	crud.funcTrash = config.FuncTrash
	crud.funcUpdate = config.FuncUpdate
	crud.homeURL = config.HomeURL
	crud.createRedirectURL = config.CreateRedirectURL
	crud.updateRedirectURL = config.UpdateRedirectURL
	crud.updateFields = config.UpdateFields
	crud.funcBeforeAction = config.FuncBeforeAction
	crud.funcAfterAction = config.FuncAfterAction
	crud.funcValidateCSRF = config.FuncValidateCSRF

	return crud, err
}
