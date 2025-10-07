package dto

type CreateApiKeyRequest struct {
	Name          string `json:"name" binding:"required"`
	CollectionIds []int  `json:"collection_ids" binding:"required"`
}
