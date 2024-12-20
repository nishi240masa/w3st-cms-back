package models



type Site struct {
	ID        int  `gorm:"type:int;primary_key;auto_increment" json:"id"`
	UserId	UUID `gorm:"type:uuid;not null" json:"userId"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Url     string         `gorm:"type:text;not null" json:"url"`
	CreateAt  string `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt  string `gorm:"type:timestamp;not null" json:"updateAt"`
}