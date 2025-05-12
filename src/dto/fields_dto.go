package dto

type Field struct {
	FieldID      string `json:"field_id"`
	ViewName     string `json:"view_name"`
	FieldType    string `json:"field_type"`
	IsRequired   bool   `json:"is_required"`
	DefaultValue string `json:"default_value"`
}
