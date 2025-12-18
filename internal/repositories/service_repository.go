package repositories

import (
	"context"
	"fmt"

	"github.com/Pelfox/gidock/internal/models"
	"github.com/google/uuid"
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
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateService: failed to build query: %w", err)
	}

	var service models.Service
	if err := r.pool.QueryRow(context.Background(), query, args...).Scan(
		&service.ID,
		&service.CreatedAt,
		&service.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("CreateService: failed to execute query: %w", err)
	}

	// set other fields from request
	service.ProjectID = projectID
	service.Name = name
	service.Image = image
	service.Environment = environment
	service.Mounts = mounts
	service.Dependencies = dependencies
	service.NetworkAccess = networkAccess

	return &service, nil
}
