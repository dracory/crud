package crud

// FieldOption represents a key-value pair for select, radio, and other option-based fields.
type FieldOption struct {
	Key   string
	Value string
}

// FieldInterface defines the contract for fields used in crud create/update forms.
type FieldInterface interface {
	GetID() string
	GetLabel() string
	GetName() string
	GetHelp() string
	GetOptions() []FieldOption
	GetOptionsF() func() []FieldOption
	GetRequired() bool
	GetType() string
	GetValue() string
	// GetFields returns sub-fields for repeater fields; nil for all other types.
	GetFields() []FieldInterface
	// GetFuncValuesProcess returns a function that converts a raw string value
	// (e.g. JSON from FuncFetchUpdateData) into []map[string]string rows for Vue binding.
	// Returns nil for field types that do not need pre-processing.
	GetFuncValuesProcess() func(raw string) []map[string]string
}

// FieldOptions configures a new Field instance.
type FieldOptions struct {
	ID                string
	Type              string
	Name              string
	Label             string
	Help              string
	Options           []FieldOption
	OptionsF          func() []FieldOption
	Value             string
	Required          bool
	Fields            []FieldInterface
	FuncValuesProcess func(raw string) []map[string]string
}

// Field is the default FieldInterface implementation.
type Field struct {
	id                string
	fieldType         string
	name              string
	label             string
	help              string
	options           []FieldOption
	optionsF          func() []FieldOption
	value             string
	required          bool
	fields            []FieldInterface
	funcValuesProcess func(raw string) []map[string]string
}

// NewField creates a new Field with the given options.
func NewField(opts FieldOptions) *Field {
	return &Field{
		id:                opts.ID,
		fieldType:         opts.Type,
		name:              opts.Name,
		label:             opts.Label,
		help:              opts.Help,
		options:           opts.Options,
		optionsF:          opts.OptionsF,
		value:             opts.Value,
		required:          opts.Required,
		fields:            opts.Fields,
		funcValuesProcess: opts.FuncValuesProcess,
	}
}

func (f *Field) GetID() string                                        { return f.id }
func (f *Field) GetType() string                                      { return f.fieldType }
func (f *Field) GetName() string                                      { return f.name }
func (f *Field) GetLabel() string                                     { return f.label }
func (f *Field) GetHelp() string                                      { return f.help }
func (f *Field) GetOptions() []FieldOption                            { return f.options }
func (f *Field) GetOptionsF() func() []FieldOption                    { return f.optionsF }
func (f *Field) GetValue() string                                     { return f.value }
func (f *Field) GetRequired() bool                                    { return f.required }
func (f *Field) GetFields() []FieldInterface                          { return f.fields }
func (f *Field) GetFuncValuesProcess() func(raw string) []map[string]string {
	return f.funcValuesProcess
}

// RepeaterOptions configures a new repeater field instance.
type RepeaterOptions struct {
	Label             string
	Name              string
	Help              string
	Fields            []FieldInterface
	// FuncValuesProcess converts a raw JSON string (from FuncFetchUpdateData)
	// into []map[string]string rows for Vue binding. Optional.
	FuncValuesProcess func(raw string) []map[string]string
}

// NewRepeater creates a repeater FieldInterface for use in crud create/update forms.
func NewRepeater(opts RepeaterOptions) *Field {
	return &Field{
		fieldType:         FORM_FIELD_TYPE_REPEATER,
		name:              opts.Name,
		label:             opts.Label,
		help:              opts.Help,
		fields:            opts.Fields,
		funcValuesProcess: opts.FuncValuesProcess,
	}
}
