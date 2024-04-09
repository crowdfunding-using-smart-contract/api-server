package usecase

import (
	"errors"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/apperrors"
	"fund-o/api-server/pkg/pagination"
	"fund-o/api-server/pkg/uploader"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ProjectUseCase interface {
	ListProjects(params entity.ProjectListParams) pagination.PaginateResult[entity.ProjectDto]
	CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error)
	GetProjectByID(projectID string) (*entity.ProjectDto, apperrors.Error)
	GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error)
	GetRecommendationProjects() ([]entity.ProjectDto, error)
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

func (uc *projectUseCase) ListProjects(params entity.ProjectListParams) pagination.PaginateResult[entity.ProjectDto] {
	result := pagination.MakePaginateResult(pagination.MakePaginateContextParameters[entity.ProjectDto]{
		PaginateOptions: params.PaginateOptions,
		CountDocuments: func() int64 {
			return uc.projectRepository.Count()
		},
		FindDocuments: func(findOptions pagination.PaginateFindOptions) []entity.ProjectDto {
			parsedCategoryId, err := uuid.Parse(params.CategoryID)
			if err != nil {
				parsedCategoryId = uuid.Nil
			}

			parsedSubCategoryId, err := uuid.Parse(params.SubCategoryID)
			if err != nil {
				parsedSubCategoryId = uuid.Nil
			}

			documents := uc.projectRepository.FindAll(findOptions, entity.ProjectListOptions{
				Query:         params.Query,
				CategoryID:    parsedCategoryId,
				SubCategoryID: parsedSubCategoryId,
			})

			projectDtos := make([]entity.ProjectDto, 0, len(documents))
			for _, document := range documents {
				projectDtos = append(projectDtos, *document.ToProjectDto())
			}

			return projectDtos
		},
	})

	return result
}

func (uc *projectUseCase) CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error) {
	ownerID := uuid.MustParse(project.OwnerID)

	endDate, err := time.Parse(time.RFC3339, project.EndDate)
	if err != nil {
		return nil, err
	}

	image, err := uc.imageUploader.Upload(uploader.ProjectImageFolder, project.Image)
	if err != nil {
		return nil, err
	}

	payload := &entity.Project{
		ProjectContractID: project.ProjectContractID,
		Title:             project.Title,
		SubTitle:          project.SubTitle,
		Description:       project.Description,
		CategoryID:        uuid.MustParse(project.CategoryID),
		SubCategoryID:     uuid.MustParse(project.SubCategoryID),
		Location:          project.Location,
		Image:             image,
		StartDate:         time.Now(),
		EndDate:           endDate,
		OwnerID:           ownerID,
	}
	newProject, err := uc.projectRepository.Create(payload)
	if err != nil {
		return nil, err
	}

	return newProject.ToProjectDto(), nil
}

func (uc *projectUseCase) GetProjectByID(projectID string) (*entity.ProjectDto, apperrors.Error) {
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, apperrors.New(http.StatusBadRequest, "Invalid project ID")
	}

	project, err := uc.projectRepository.FindByID(projectUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(http.StatusNotFound, "Project not found")
		}

		return nil, apperrors.New(http.StatusInternalServerError, "Failed to get project")
	}

	return project.ToProjectDto(), nil
}

func (uc *projectUseCase) GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error) {
	ownerID, err := uuid.Parse(requestOwnerID)
	if err != nil {
		return nil, apperrors.ErrInvalidUserID
	}

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

func (uc *projectUseCase) GetRecommendationProjects() ([]entity.ProjectDto, error) {
	projects, err := uc.projectRepository.FindRecommendation(3)
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
