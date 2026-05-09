package repository

import (
	"context"
	"database/sql"
	"fmt"

	"sms-backend/internal/models"
)

type PageRepo struct{ db *sql.DB }

func NewPageRepo(db *sql.DB) *PageRepo { return &PageRepo{db: db} }

func (r *PageRepo) List(ctx context.Context) ([]models.Page, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, slug, created_at, updated_at FROM pages ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	pages := make([]models.Page, 0)
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

func (r *PageRepo) Get(ctx context.Context, id string) (*models.Page, error) {
	var p models.Page
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, slug, created_at, updated_at FROM pages WHERE id = $1`, id,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PageRepo) GetBySlug(ctx context.Context, slug string) (*models.Page, error) {
	var p models.Page
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, slug, created_at, updated_at FROM pages WHERE slug = $1`, slug,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PageRepo) Create(ctx context.Context, title, slug string) (*models.Page, error) {
	var p models.Page
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO pages (title, slug) VALUES ($1, $2)
		 RETURNING id, title, slug, created_at, updated_at`,
		title, slug,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &p, nil
}

func (r *PageRepo) Update(ctx context.Context, id, title, slug string) (*models.Page, error) {
	var p models.Page
	err := r.db.QueryRowContext(ctx,
		`UPDATE pages SET title = $2, slug = $3, updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, title, slug, created_at, updated_at`,
		id, title, slug,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PageRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pages WHERE id = $1`, id)
	return err
}
