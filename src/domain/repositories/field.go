package repositories

import (
	"w3st/domain/models"

	"github.com/google/uuid"
)

type FieldRepository interface {
	CreateField(newField *models.FieldData) error
	UpdateField(newField *models.FieldData) error
	DeleteFieldById(userId uuid.UUID, fieldId uuid.UUID) error
}
