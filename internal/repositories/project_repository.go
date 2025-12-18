package repositories

import (
	"context"
	"fmt"

	s "github.com/Masterminds/squirrel"
	"github.com/Pelfox/gidock/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
	psql s.StatementBuilderType
}

func NewProjectRepository(psql s.StatementBuilderType, pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{psql: psql, pool: pool}
}

func (r *ProjectRepository) CreateProject(name string) (*models.Project, error) {
	query, args, err := r.psql.Insert("projects").
		Columns("name").
		Values(name).
		Suffix("RETURNING id, name, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateProject: failed to build query: %w", err)
	}

	var project models.Project
	if err := r.pool.QueryRow(context.Background(), query, args...).Scan(
		&project.ID,
		&project.Name,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("CreateProject: failed to execute query: %w", err)
	}

	return &project, nil
}
