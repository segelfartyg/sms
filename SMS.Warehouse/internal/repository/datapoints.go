package repository

import (
	"context"
	"database/sql"
	"fmt"

	"sms-warehouse/internal/models"
)

type DatapointRepo struct{ db *sql.DB }

func NewDatapointRepo(db *sql.DB) *DatapointRepo { return &DatapointRepo{db: db} }

func (r *DatapointRepo) List(ctx context.Context, datasourceID string) ([]models.Datapoint, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, datasource_id, tag, content, description, created_at, updated_at
		 FROM datapoints WHERE datasource_id = $1 ORDER BY created_at ASC`,
		datasourceID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	out := make([]models.Datapoint, 0)
	for rows.Next() {
		var dp models.Datapoint
		if err := rows.Scan(&dp.ID, &dp.DatasourceID, &dp.Tag, &dp.Content, &dp.Description, &dp.CreatedAt, &dp.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, dp)
	}
	return out, rows.Err()
}

func (r *DatapointRepo) Get(ctx context.Context, datasourceID, id string) (*models.Datapoint, error) {
	var dp models.Datapoint
	err := r.db.QueryRowContext(ctx,
		`SELECT id, datasource_id, tag, content, description, created_at, updated_at
		 FROM datapoints WHERE id = $1 AND datasource_id = $2`,
		id, datasourceID,
	).Scan(&dp.ID, &dp.DatasourceID, &dp.Tag, &dp.Content, &dp.Description, &dp.CreatedAt, &dp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (r *DatapointRepo) Create(ctx context.Context, datasourceID, tag, content, description string) (*models.Datapoint, error) {
	var dp models.Datapoint
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO datapoints (datasource_id, tag, content, description)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, datasource_id, tag, content, description, created_at, updated_at`,
		datasourceID, tag, content, description,
	).Scan(&dp.ID, &dp.DatasourceID, &dp.Tag, &dp.Content, &dp.Description, &dp.CreatedAt, &dp.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &dp, nil
}

func (r *DatapointRepo) Update(ctx context.Context, datasourceID, id, tag, content, description string) (*models.Datapoint, error) {
	var dp models.Datapoint
	err := r.db.QueryRowContext(ctx,
		`UPDATE datapoints SET tag = $3, content = $4, description = $5, updated_at = NOW()
		 WHERE id = $1 AND datasource_id = $2
		 RETURNING id, datasource_id, tag, content, description, created_at, updated_at`,
		id, datasourceID, tag, content, description,
	).Scan(&dp.ID, &dp.DatasourceID, &dp.Tag, &dp.Content, &dp.Description, &dp.CreatedAt, &dp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (r *DatapointRepo) Delete(ctx context.Context, datasourceID, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM datapoints WHERE id = $1 AND datasource_id = $2`, id, datasourceID)
	return err
}
