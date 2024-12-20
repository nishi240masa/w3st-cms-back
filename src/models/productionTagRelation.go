package models

type ProductionTagRelation struct {
	ID        int  `gorm:"type:int;primary_key;auto_increment" json:"id"`
	ProductionId	int `gorm:"type:int;not null" json:"productionId"`
	TagId	int `gorm:"type:int;not null" json:"tagId"`
}