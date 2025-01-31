package models


type ListOptions struct {
	ID        int     `gorm:"type:serial;primary_key" json:"id"`
	ApiSchemaID    int    `gorm:"type:int;not null" json:"apiSchemaId"`
	Value  string  `gorm:"type:varchar(50);not null" json:"value"`
	CreateAt  string  `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt  string  `gorm:"type:timestamp;not null" json:"updateAt"`
}