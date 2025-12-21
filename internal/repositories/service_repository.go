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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ServiceRepository provides data access methods for the `services` table.
type ServiceRepository struct {
	pool *pgxpool.Pool
}

// NewServiceRepository creates a new ServiceRepository instance from the given `*pgxpool.Pool`.
func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{pool: pool}
}

// Create creates a new service with the given command and returns it.
func (r *ServiceRepository) Create(
	ctx context.Context,
	command commands.CreateServiceCommand,
) (*models.Service, error) {
	query, args, err := sq.Insert("services").
		Columns("project_id", "name", "image", "environment", "mounts", "dependencies", "network_access").
		Values(
			command.ProjectID,
			command.Name,
			command.Image,
			command.Environment,
			command.Mounts,
			command.Dependencies,
			command.NetworkAccess,
		).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Create: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, internal.ErrRelationNotFound
		}
		return nil, fmt.Errorf("Create: failed to execute query: %w", err)
	}
	defer rows.Close()

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		return nil, fmt.Errorf("Create: failed to map: %w", err)
	}

	return &service, nil
}

// Get retrieves a service with given command.
func (r *ServiceRepository) Get(
	ctx context.Context,
	command commands.GetServiceCommand,
) (*models.Service, error) {
	query, args, err := sq.Select("*").
		From("services").
		Where(s.Eq{"id": command.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Get: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Get: failed to execute query: %w", err)
	}
	defer rows.Close()

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Get: failed to map: %w", err)
	}

	return &service, nil
}

// Update updates an existing service in the database by its id. It
// applies only the fields that are set in the provided `commands.UpdateServiceCommand`.
func (r *ServiceRepository) Update(
	ctx context.Context,
	command commands.UpdateServiceCommand,
) (*models.Service, error) {
	queryBuilder := sq.Update("services")

	// updating all selected (non-nil) fields
	if command.ContainerID != nil {
		queryBuilder = queryBuilder.Set("container_id", *command.ContainerID)
	}

	// if update fields are empty, return an error
	if queryBuilder == sq.Update("services") {
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

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Update: failed to map: %w", err)
	}

	return &service, nil
}

// Delete removes a service from the database with given command.
func (r *ServiceRepository) Delete(
	ctx context.Context,
	command commands.DeleteServiceCommand,
) error {
	query, args, err := sq.Delete("services").
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

// ListAll retrieves all services from the database.
func (r *ServiceRepository) ListAll(ctx context.Context) ([]models.Service, error) {
	query, args, err := s.Select("*").From("services").ToSql()
	if err != nil {
		return nil, fmt.Errorf("ListAll: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListAll: failed to execute query: %w", err)
	}
	defer rows.Close()

	services, err := pgx.CollectRows[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		return nil, fmt.Errorf("ListAll: failed to map: %w", err)
	}

	return services, nil
}
