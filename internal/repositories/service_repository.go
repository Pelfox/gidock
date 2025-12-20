package repositories

import (
	"context"
	"errors"
	"fmt"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/Pelfox/gidock/internal/repositories/commands"
	"github.com/google/uuid"
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

// CreateService inserts a new service record and returns the created *models.Service.
func (r *ServiceRepository) CreateService(
	projectID uuid.UUID,
	name string,
	image string,
	environment map[string]string,
	mounts []models.ServiceMount,
	dependencies []models.ServiceDependency,
	networkAccess bool,
) (*models.Service, error) {
	query, args, err := sq.Insert("services").
		Columns("project_id", "name", "image", "environment", "mounts", "dependencies", "network_access").
		Values(projectID.String(), name, image, environment, mounts, dependencies, networkAccess).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateService: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, internal.ErrRelationNotFound
		}
		return nil, fmt.Errorf("CreateService: failed to execute query: %w", err)
	}
	defer rows.Close()

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		return nil, fmt.Errorf("CreateService: failed to map: %w", err)
	}

	return &service, nil
}

// GetServiceByID retrieves a single service by its id.
func (r *ServiceRepository) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	query, args, err := sq.Select("*").
		From("services").
		Where(s.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetServiceByID: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetServiceByID: failed to execute query: %w", err)
	}
	defer rows.Close()

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("GetServiceByID: failed to map: %w", err)
	}

	return &service, nil
}

// UpdateServiceByID updates an existing service in the database by its id. It
// applies only the fields that are set in the provided `commands.UpdateServiceCommand`.
func (r *ServiceRepository) UpdateServiceByID(
	id uuid.UUID,
	command commands.UpdateServiceCommand,
) (*models.Service, error) {
	queryBuilder := sq.Update("services").
		Where(s.Eq{"id": id}).
		Suffix("RETURNING *")

	if command.ContainerID != nil {
		queryBuilder = queryBuilder.Set("container_id", *command.ContainerID)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UpdateServiceByID: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("UpdateServiceByID: failed to execute query: %w", err)
	}
	defer rows.Close()

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("UpdateServiceByID: failed to map: %w", err)
	}

	return &service, nil
}

// ListServices retrieves all services from the database.
func (r *ServiceRepository) ListServices() ([]models.Service, error) {
	query, args, err := s.Select("*").From("services").ToSql()
	if err != nil {
		return nil, fmt.Errorf("ListServices: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListServices: failed to execute query: %w", err)
	}
	defer rows.Close()

	services, err := pgx.CollectRows[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, internal.ErrRecordNotFound
		}
		return nil, fmt.Errorf("ListServices: failed to map: %w", err)
	}

	return services, nil
}
