package crud

import (
	"net/http"

	"github.com/dracory/hb"
)

type Config struct {
	ColumnNames         []string
	CreateFields        []FieldInterface
	Endpoint            string
	EntityNamePlural    string
	EntityNameSingular  string
	FileManagerURL      string
	FuncCreate          func(r *http.Request, data map[string]string) (userID string, err error)
	FuncFetchReadData   func(r *http.Request, entityID string) ([]KeyValue, error)
	FuncFetchUpdateData func(r *http.Request, entityID string) (map[string]string, error)
	FuncLayout          func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string
	FuncRows            func(r *http.Request) (rows []Row, err error)
	FuncTrash           func(r *http.Request, entityID string) error
	FuncUpdate          func(r *http.Request, entityID string, data map[string]string) error
	HomeURL             string
	CreateRedirectURL   string
	UpdateRedirectURL   string
	UpdateFields        []FieldInterface
	FuncReadExtras      func(r *http.Request, entityID string) []hb.TagInterface
	PageSize            int
	FuncRowsCount       func(r *http.Request) (int64, error)
	FuncBeforeAction    func(w http.ResponseWriter, r *http.Request, action string) bool
	FuncAfterAction     func(w http.ResponseWriter, r *http.Request, action string)
	FuncValidateCSRF    func(r *http.Request) error
	FuncLog             func(level string, message string, attrs map[string]any)
}
