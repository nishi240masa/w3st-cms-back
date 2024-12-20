package models

type Tag struct {
	ID        int  `gorm:"type:int;primary_key;auto_increment" json:"id"`
	Word 	string         `gorm:"type:varchar(50);not null" json:"word"`
}