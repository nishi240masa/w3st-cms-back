package models

type ApiKeyCollections struct {
	ApiKeyID     int `gorm:"type:int;primary_key" json:"api_key_Id"`
	CollectionID int `gorm:"type:int;primary_key" json:"collection_id"`
}
