package usecase

import (
	"errors"
	"fmt"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/apperrors"
	"fund-o/api-server/pkg/uploader"
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProjectUseCase interface {
	CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error)
	GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error)
	CreateProjectRating(rating *entity.ProjectRatingCreatePayload) error
	IsRatedProject(userID string, projectID string) (bool, error)
}

type projectUseCase struct {
	projectRepository repository.ProjectRepository
	imageUploader     uploader.ImageUploader
}

type ProjectUseCaseOptions struct {
	repository.ProjectRepository
	uploader.ImageUploader
}

func NewProjectUseCase(options *ProjectUseCaseOptions) ProjectUseCase {
	return &projectUseCase{
		projectRepository: options.ProjectRepository,
		imageUploader:     options.ImageUploader,
	}
}

func (uc *projectUseCase) CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error) {
	ownerID := uuid.MustParse(project.OwnerID)

	endDate, err := time.Parse(time.RFC3339, project.EndDate)
	if err != nil {
		return nil, err
	}

	var launchDate = endDate
	if project.LaunchDate != "" {
		if launchDate, err = time.Parse(time.RFC3339, project.LaunchDate); err != nil {
			fmt.Println("Error while parsing launch date: ", err)
			return nil, err
		}
	}

	targetFunding, err := decimal.NewFromString(project.TargetFunding)
	if err != nil {
		return nil, err
	}

	image, err := uc.imageUploader.Upload(uploader.ProjectImageFolder, project.Image)
	if err != nil {
		return nil, err
	}

	payload := &entity.Project{
		Title:          project.Title,
		SubTitle:       project.SubTitle,
		CategoryID:     uuid.MustParse(project.CategoryID),
		SubCategoryID:  uuid.MustParse(project.SubCategoryID),
		Location:       project.Location,
		Image:          image,
		Description:    project.Description,
		TargetFunding:  targetFunding,
		MonetaryUnit:   project.MonetaryUnit,
		StartDate:      time.Now(),
		EndDate:        endDate,
		LaunchDate:     launchDate,
		CurrentFunding: decimal.NewFromInt(0),
		OwnerID:        ownerID,
	}
	newProject, err := uc.projectRepository.Create(payload)
	if err != nil {
		return nil, err
	}

	return newProject.ToProjectDto(), nil
}

func (uc *projectUseCase) GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error) {
	ownerID := uuid.MustParse(requestOwnerID)

	projects, err := uc.projectRepository.FindAllByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	projectDtos := make([]entity.ProjectDto, 0, len(projects))
	for _, project := range projects {
		projectDtos = append(projectDtos, *project.ToProjectDto())
	}

	return projectDtos, nil
}
func (uc *projectUseCase) CreateProjectRating(rating *entity.ProjectRatingCreatePayload) error {
	projectID, err := uuid.Parse(rating.ProjectID)
	if err != nil {
		return apperrors.ErrInvalidProjectID
	}

	userID, err := uuid.Parse(rating.UserID)
	if err != nil {
		return apperrors.ErrInvalidUserID
	}

	pr, prErr := uc.projectRepository.FindProjectRating(userID, projectID)
	if prErr != nil && !errors.Is(prErr, gorm.ErrRecordNotFound) {
		return prErr
	}

	if pr.ID != uuid.Nil {
		return apperrors.ErrAlreadyRatedProject
	}

	_, err = uc.projectRepository.CreateProjectRating(&entity.ProjectRating{
		Rating:    rating.Rating,
		ProjectID: projectID,
		UserID:    userID,
	})
	return err
}
func (uc *projectUseCase) IsRatedProject(userID string, projectID string) (bool, error) {
	userIDParsed := uuid.MustParse(userID)
	projectIDParsed := uuid.MustParse(projectID)
	pr, err := uc.projectRepository.FindProjectRating(userIDParsed, projectIDParsed)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if pr.ID != uuid.Nil {
		return true, nil
	}

	return false, nil
}
