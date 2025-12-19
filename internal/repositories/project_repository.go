package repositories

import (
	"context"
	"fmt"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) CreateProject(name string) (*models.Project, error) {
	query, args, err := sq.Insert("projects").
		Columns("name").
		Values(name).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateProject: failed to build query: %w", err)
	}

	var project models.Project
	if err := r.pool.QueryRow(context.Background(), query, args...).Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("CreateProject: failed to execute query: %w", err)
	}

	project.Name = name
	return &project, nil
}

func (r *ProjectRepository) ListProjects() ([]models.Project, error) {
	query, args, err := sq.Select("*").From("projects").ToSql()
	if err != nil {
		return nil, fmt.Errorf("ListProjects: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListProjects: failed to execute query: %w", err)
	}
	defer rows.Close()

	projects, err := pgx.CollectRows[models.Project](rows, pgx.RowToStructByName[models.Project])
	if err != nil {
		return nil, fmt.Errorf("ListProjects: failed to map: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepository) GetProjectById(id uuid.UUID, includeServices bool) (*models.ProjectWithServices, error) {
	query, args, err := sq.Select("*").
		From("projects").
		Where(s.Eq{"id": id.String()}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetProjectById: failed to build query: %w", err)
	}

	projectRows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetProjectById: failed to execute query: %w", err)
	}
	defer projectRows.Close()

	project, err := pgx.CollectOneRow[models.Project](projectRows, pgx.RowToStructByName[models.Project])
	if err != nil {
		return nil, fmt.Errorf("GetProjectById: failed to map: %w", err)
	}

	result := models.ProjectWithServices{Project: project}
	if !includeServices {
		return &result, nil
	}

	servicesQuery, servicesArgs, servicesErr := sq.Select("*").
		From("services").
		Where(s.Eq{"project_id": id.String()}).
		ToSql()
	if servicesErr != nil {
		return nil, fmt.Errorf("GetProjectById: failed to build query for services: %w", servicesErr)
	}

	servicesRows, err := r.pool.Query(context.Background(), servicesQuery, servicesArgs...)
	if err != nil {
		return nil, fmt.Errorf("GetProjectById: failed to execute services query: %w", err)
	}
	defer servicesRows.Close()

	services, err := pgx.CollectRows[models.Service](servicesRows, pgx.RowToStructByName[models.Service])
	if err != nil {
		return nil, fmt.Errorf("GetProjectById: failed to map services: %w", err)
	}

	result.Services = &services
	return &result, nil
}
