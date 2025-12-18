package repositories

import (
	"context"
	"fmt"

	"github.com/Pelfox/gidock/internal/models"
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
