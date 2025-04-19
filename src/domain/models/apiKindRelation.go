package models

type ApiKindRelation struct {
	ID           int    `gorm:"type:serial;primary_key" json:"id"`
	ApiSchemaID  int    `gorm:"type:int;not null" json:"apiSchemaId"`
	RelatedID    int    `gorm:"type:int;not null" json:"relatedId"`
	RelationType string `gorm:"type:varchar(50);not null" json:"relationType"`
	CreateAt     string `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt     string `gorm:"type:timestamp;not null" json:"updateAt"`
}
