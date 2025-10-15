package usecase

import (
	"context"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

type ProjectUsecase interface {
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	CreateProject(ctx context.Context, name, description string, rateLimit int) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	GetRateLimitByProjectID(ctx context.Context, projectID int) (int, error)
}

type projectUsecase struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectUsecase(projectRepo repositories.ProjectRepository) ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

func (u *projectUsecase) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	project, err := u.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, myerrors.WrapDomainError("projectUsecase.GetProjectByID", err)
	}

	return project, nil
}

func (u *projectUsecase) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	projects, err := u.projectRepo.FindAll(ctx)
	if err != nil {
		return nil, myerrors.WrapDomainError("projectUsecase.GetAllProjects", err)
	}

	return projects, nil
}

func (u *projectUsecase) CreateProject(ctx context.Context, name, description string, rateLimit int) (*models.Project, error) {
	project := &models.Project{
		Name:             name,
		Description:      description,
		RateLimitPerHour: rateLimit,
	}

	createdProject, err := u.projectRepo.Create(ctx, project)
	if err != nil {
		return nil, myerrors.WrapDomainError("projectUsecase.CreateProject", err)
	}

	return createdProject, nil
}

func (u *projectUsecase) UpdateProject(ctx context.Context, project *models.Project) error {
	err := u.projectRepo.Update(ctx, project)
	if err != nil {
		return myerrors.WrapDomainError("projectUsecase.UpdateProject", err)
	}

	return nil
}

func (u *projectUsecase) GetRateLimitByProjectID(ctx context.Context, projectID int) (int, error) {
	project, err := u.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return 0, myerrors.WrapDomainError("projectUsecase.GetRateLimitByProjectID", err)
	}

	return project.RateLimitPerHour, nil
}
