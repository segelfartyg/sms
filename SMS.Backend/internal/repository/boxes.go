package repository

import (
	"context"
	"database/sql"
	"fmt"

	"sms-backend/internal/models"
)

type BoxRepo struct{ db *sql.DB }

func NewBoxRepo(db *sql.DB) *BoxRepo { return &BoxRepo{db: db} }

func (r *BoxRepo) List(ctx context.Context, pageID string) ([]models.Box, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, page_id, position, datasource_id, created_at, updated_at
		 FROM boxes WHERE page_id = $1 ORDER BY position ASC`,
		pageID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	boxes := make([]models.Box, 0)
	for rows.Next() {
		b, err := scanBox(rows)
		if err != nil {
			return nil, err
		}
		boxes = append(boxes, b)
	}
	return boxes, rows.Err()
}

func (r *BoxRepo) Get(ctx context.Context, pageID, id string) (*models.Box, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, page_id, position, datasource_id, created_at, updated_at
		 FROM boxes WHERE id = $1 AND page_id = $2`,
		id, pageID)
	b, err := scanBox(row)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BoxRepo) Create(ctx context.Context, pageID string, position int, datasourceID *string) (*models.Box, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO boxes (page_id, position, datasource_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, page_id, position, datasource_id, created_at, updated_at`,
		pageID, position, datasourceID)
	b, err := scanBox(row)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &b, nil
}

func (r *BoxRepo) Update(ctx context.Context, pageID, id string, position int, datasourceID *string) (*models.Box, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE boxes SET position = $3, datasource_id = $4, updated_at = NOW()
		 WHERE id = $1 AND page_id = $2
		 RETURNING id, page_id, position, datasource_id, created_at, updated_at`,
		id, pageID, position, datasourceID)
	b, err := scanBox(row)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BoxRepo) Delete(ctx context.Context, pageID, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM boxes WHERE id = $1 AND page_id = $2`, id, pageID)
	return err
}

type scanner interface{ Scan(dest ...any) error }

func scanBox(s scanner) (models.Box, error) {
	var b models.Box
	var dsID sql.NullString
	if err := s.Scan(&b.ID, &b.PageID, &b.Position, &dsID, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return b, fmt.Errorf("scan: %w", err)
	}
	if dsID.Valid {
		b.DatasourceID = &dsID.String
	}
	return b, nil
}
