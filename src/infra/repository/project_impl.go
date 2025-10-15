package infrastructure

import (
	"context"
	"w3st/domain/models"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (r *ProjectRepository) FindByID(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	result := r.db.WithContext(ctx).First(&project, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, result.Error)
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return &project, nil
}

func (r *ProjectRepository) FindAll(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.WithContext(ctx).Find(&projects)

	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return projects, nil
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	result := r.db.WithContext(ctx).Create(project)

	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return project, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	result := r.db.WithContext(ctx).Save(project)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return nil
}