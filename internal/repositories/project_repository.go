package repositories

import (
	"context"
	"errors"
	"fmt"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProjectRepository provides data access methods for the `projects` table.
type ProjectRepository struct {
	pool *pgxpool.Pool
}

// NewProjectRepository creates a new ProjectRepository instance from the given `*pgxpool.Pool`.
func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

// CreateProject creates a new project in the database with the given name and
// returns the created *models.Project.
func (r *ProjectRepository) CreateProject(name string) (*models.Project, error) {
	query, args, err := sq.Insert("projects").
		Columns("name").
		Values(name).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateProject: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("CreateProject: failed to execute query: %w", err)
	}
	defer rows.Close()

	project, err := pgx.CollectOneRow[models.Project](rows, pgx.RowToStructByName[models.Project])
	if err != nil {
		return nil, fmt.Errorf("CreateProject: failed to map: %w", err)
	}

	return &project, nil
}

// GetProjectById retrieves a project by its id, optionally including its
// associated services.
func (r *ProjectRepository) GetProjectById(
	id uuid.UUID,
	includeServices bool,
) (*models.ProjectWithServices, error) {
	query, args, err := sq.Select("*").
		From("projects").
		Where(s.Eq{"id": id}).
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("GetProjectById: failed to map: %w", err)
	}

	result := models.ProjectWithServices{Project: project}
	if !includeServices {
		return &result, nil
	}

	servicesQuery, servicesArgs, servicesErr := sq.Select("*").
		From("services").
		Where(s.Eq{"project_id": id}).
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("GetProjectById: failed to map services: %w", err)
	}

	result.Services = &services
	return &result, nil
}

// ListProjects retrieves all projects from the database.
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("ListProjects: failed to map: %w", err)
	}

	return projects, nil
}
