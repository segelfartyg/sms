package repository

import (
	"context"
	"database/sql"
	"fmt"

	"sms-warehouse/internal/models"
)

type DatasourceRepo struct{ db *sql.DB }

func NewDatasourceRepo(db *sql.DB) *DatasourceRepo { return &DatasourceRepo{db: db} }

func (r *DatasourceRepo) List(ctx context.Context) ([]models.Datasource, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, type, description, content, created_at, updated_at
		 FROM datasources ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	out := make([]models.Datasource, 0)
	for rows.Next() {
		var d models.Datasource
		if err := rows.Scan(&d.ID, &d.Type, &d.Description, &d.Content, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DatasourceRepo) Get(ctx context.Context, id string) (*models.Datasource, error) {
	var d models.Datasource
	err := r.db.QueryRowContext(ctx,
		`SELECT id, type, description, content, created_at, updated_at
		 FROM datasources WHERE id = $1`, id,
	).Scan(&d.ID, &d.Type, &d.Description, &d.Content, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DatasourceRepo) Create(ctx context.Context, dsType, description, content string) (*models.Datasource, error) {
	var d models.Datasource
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO datasources (type, description, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, type, description, content, created_at, updated_at`,
		dsType, description, content,
	).Scan(&d.ID, &d.Type, &d.Description, &d.Content, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &d, nil
}

func (r *DatasourceRepo) Update(ctx context.Context, id, dsType, description, content string) (*models.Datasource, error) {
	var d models.Datasource
	err := r.db.QueryRowContext(ctx,
		`UPDATE datasources SET type = $2, description = $3, content = $4, updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, type, description, content, created_at, updated_at`,
		id, dsType, description, content,
	).Scan(&d.ID, &d.Type, &d.Description, &d.Content, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DatasourceRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM datasources WHERE id = $1`, id)
	return err
}
