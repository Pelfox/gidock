package repositories

import (
	"context"
	"errors"
	"fmt"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories/commands"
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

// Create creates a new project in the database with the given command and
// returns it.
func (r *ProjectRepository) Create(
	ctx context.Context,
	command commands.CreateProjectCommand,
) (*models.Project, error) {
	query, args, err := sq.Insert("projects").
		Columns("name").
		Values(command.Name).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Create: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Create: failed to execute query: %w", err)
	}
	defer rows.Close()

	project, err := pgx.CollectOneRow[models.Project](rows, pgx.RowToStructByName[models.Project])
	if err != nil {
		return nil, fmt.Errorf("Create: failed to map: %w", err)
	}

	return &project, nil
}

// Get retrieves a project with given command. If `IncludeServices` is true,
// associated services are also retrieved.
func (r *ProjectRepository) Get(
	ctx context.Context,
	command commands.GetProjectCommand,
) (*models.ProjectWithServices, error) {
	query, args, err := sq.Select("*").
		From("projects").
		Where(s.Eq{"id": command.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Get: failed to build query: %w", err)
	}

	projectRows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Get: failed to execute query: %w", err)
	}
	defer projectRows.Close()

	project, err := pgx.CollectOneRow[models.Project](projectRows, pgx.RowToStructByName[models.Project])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Get: failed to map: %w", err)
	}

	result := models.ProjectWithServices{Project: project}
	if !command.IncludeServices {
		return &result, nil
	}

	servicesQuery, servicesArgs, servicesErr := sq.Select("*").
		From("services").
		Where(s.Eq{"project_id": command.ID}).
		ToSql()
	if servicesErr != nil {
		return nil, fmt.Errorf("Get: failed to build query for services: %w", servicesErr)
	}

	servicesRows, err := r.pool.Query(ctx, servicesQuery, servicesArgs...)
	if err != nil {
		return nil, fmt.Errorf("Get: failed to execute services query: %w", err)
	}
	defer servicesRows.Close()

	services, err := pgx.CollectRows[models.Service](servicesRows, pgx.RowToStructByName[models.Service])
	if err != nil {
		return nil, fmt.Errorf("Get: failed to map services: %w", err)
	}

	result.Services = &services
	return &result, nil
}

// Update performs a partial update on a project with given command and
// returns the updated project.
func (r *ProjectRepository) Update(
	ctx context.Context,
	command commands.UpdateProjectCommand,
) (*models.Project, error) {
	queryBuilder := sq.Update("projects")

	// updating all selected (non-nil) fields
	if command.Name != nil {
		queryBuilder = queryBuilder.Set("name", *command.Name)
	}

	// if update fields are empty, return an error
	if queryBuilder == sq.Update("projects") {
		return nil, internal.ErrNoFields
	}

	query, args, err := queryBuilder.Where(s.Eq{"id": command.ID}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Update: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Update: failed to execute query: %w", err)
	}
	defer rows.Close()

	updatedProject, err := pgx.CollectOneRow[models.Project](rows, pgx.RowToStructByName[models.Project])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Update: failed to map: %w", err)
	}
	return &updatedProject, nil
}

// Delete removes a project from the database with given command.
func (r *ProjectRepository) Delete(
	ctx context.Context,
	command commands.DeleteProjectCommand,
) error {
	query, args, err := sq.Delete("projects").
		Where(s.Eq{"id": command.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("Delete: failed to build query: %w", err)
	}

	cmdTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Delete: failed to execute query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return internal.ErrRecordNotFound
	}

	return nil
}

// ListAll retrieves all projects from the database.
func (r *ProjectRepository) ListAll(ctx context.Context) ([]models.Project, error) {
	query, args, err := sq.Select("*").From("projects").ToSql()
	if err != nil {
		return nil, fmt.Errorf("ListAll: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListAll: failed to execute query: %w", err)
	}
	defer rows.Close()

	projects, err := pgx.CollectRows[models.Project](rows, pgx.RowToStructByName[models.Project])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("ListAll: failed to map: %w", err)
	}

	return projects, nil
}
