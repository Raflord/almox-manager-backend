package database

import (
	"almox-manager-backend/internal/types"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func (s *service) QueryLatest() ([]types.Load, error) {
	sqlQuery := `
	SELECT * FROM loads
	ORDER BY created_at DESC
	LIMIT 10
	`
	rows, err := s.db.Query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	loads, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Load])
	if err != nil {
		return nil, err
	}

	return loads, nil
}

func (s *service) QueryDay() ([]types.LoadSummary, error) {
	now := time.Now().Format(time.DateOnly)

	sqlQuery := `
	SELECT material, SUM(average_weight) AS total_weight
	FROM loads
	WHERE created_at::date=$1
	GROUP BY material
	`
	rows, err := s.db.Query(context.Background(), sqlQuery, now)
	if err != nil {
		return nil, err
	}

	loads, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.LoadSummary])
	if err != nil {
		return nil, err
	}

	return loads, nil
}

func (s *service) QueryFiltered(inputData types.LoadFiltered) ([]types.Load, error) {
	sqlQuery := `
	SELECT *
	FROM loads
	WHERE (material = $1 OR NULLIF($1, '') IS NULL)
	  AND (
	    (
	      NULLIF($2, '') IS NOT NULL AND
	      NULLIF($3, '') IS NOT NULL AND
	      created_at::date BETWEEN $2::date AND $3::date
	    )
		OR (
		  NULLIF($2, '') IS NOT NULL AND
		  NULLIF($3, '') IS NULL AND
	      created_at::date = $2::date
		)
	    OR (
	      NULLIF($2, '') IS NULL AND
	      NULLIF($3, '') IS NULL
	    )
	  )
	ORDER BY created_at ASC
	`

	rows, err := s.db.Query(context.Background(), sqlQuery, inputData.Material, inputData.FirstDate, inputData.SeccondDate)
	if err != nil {
		return nil, err
	}

	loads, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Load])
	if err != nil {
		return nil, err
	}

	return loads, nil
}

func (s *service) CreateLoad(inputData types.Load) error {
	sqlQuery := `
	INSERT INTO loads (id, material, average_weight, unit, created_at, timezone, operator, shift)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.db.Exec(context.Background(), sqlQuery, ulid.Make().String(), inputData.Material, inputData.AverageWeight, inputData.Unit, inputData.CreatedAt, inputData.Timezone, inputData.Operator, inputData.Shift)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateLoad(inputData types.Load) error {
	sqlQuery := `
	UPDATE loads
	SET material=$1,
		created_at=$2,
		operator=$3,
		shift=$4
	WHERE id=$5
	`
	_, err := s.db.Exec(context.Background(), sqlQuery, inputData.Material, inputData.CreatedAt, inputData.Operator, inputData.Shift, inputData.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteLoad(inputData string) error {
	sqlQuery := `
	DELETE FROM loads
	WHERE id=$1
	`
	_, err := s.db.Exec(context.Background(), sqlQuery, inputData)
	if err != nil {
		return err
	}
	return nil
}
