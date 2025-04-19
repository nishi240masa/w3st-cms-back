package models

type FieldData struct {
	ID          int    `gorm:"type:serial;primary_key" json:"id"`
	ApiSchemaID int    `gorm:"type:int;not null" json:"apiSchemaId"`
	FieldType   string `gorm:"type:varchar(50);not null" json:"fieldType"`
	FieldValue  string `gorm:"type:jsonb" json:"fieldValue"`
	CreateAt    string `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt    string `gorm:"type:timestamp;not null" json:"updateAt"`
}
