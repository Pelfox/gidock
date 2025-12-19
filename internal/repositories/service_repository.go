package repositories

import (
	"context"
	"fmt"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal/dto"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	pool *pgxpool.Pool
}

func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{pool: pool}
}

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

	// TODO: validate `projectID`

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("CreateService: failed to execute query: %w", err)
	}
	defer rows.Close()

	service, err := pgx.CollectOneRow[models.Service](rows, pgx.RowToStructByName[models.Service])
	if err != nil {
		return nil, fmt.Errorf("CreateService: failed to map: %w", err)
	}

	return &service, nil
}

func (r *ServiceRepository) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	query, args, err := sq.Select("*").
		From("services").
		Where(s.Eq{"id": id.String()}).
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
		return nil, fmt.Errorf("GetServiceByID: failed to map: %w", err)
	}

	return &service, nil
}

func (r *ServiceRepository) UpdateServiceByID(id uuid.UUID, request dto.UpdateServiceFields) (*models.Service, error) {
	queryBuilder := sq.Update("services").
		Where(s.Eq{"id": id.String()}).
		Suffix("RETURNING *")

	if request.ContainerID != nil {
		queryBuilder = queryBuilder.Set("container_id", *request.ContainerID)
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
		return nil, fmt.Errorf("UpdateServiceByID: failed to map: %w", err)
	}

	return &service, nil
}

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
		return nil, fmt.Errorf("ListServices: failed to map: %w", err)
	}

	return services, nil
}
