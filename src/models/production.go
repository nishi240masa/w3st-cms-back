package models


type Production struct {
	ID        int  `gorm:"type:int;primary_key;auto_increment" json:"id"`
	SiteId	int `gorm:"type:int;not null" json:"siteId"`
	Title	  string         `gorm:"type:varchar(100);not null" json:"title"`
	Description     string         `gorm:"type:text;not null" json:"description"`
	CreateAt  string `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt  string `gorm:"type:timestamp;not null" json:"updateAt"`
}