package database

import (
	"time"

	"github.com/google/uuid"
)

type RecordInput struct {
	Material      string
	AverageWeight int
	Unit          string
	CreatedAt     *time.Time
	Operator      string
	Shift         string
}

type FilteredInput struct {
	Material    string
	FirstDate   *time.Time
	SeccondDate *time.Time
}

type LoadRecord struct {
	Id            string    `db:"id" json:"id"`
	Material      string    `db:"material" json:"material"`
	AverageWeight int       `db:"average_weight" json:"average_weight"`
	Unit          string    `db:"unit" json:"unit"`
	CreatedAt     time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `db:"updatedAt" json:"updatedAt"`
	Operator      string    `db:"operator" json:"operator"`
	Shift         string    `db:"shift" json:"shift"`
}

func (s *service) QueryLatestRecords() ([]LoadRecord, error) {
	records := make([]LoadRecord, 10)
	sqlQuery := `
	SELECT * FROM load_record
	ORDER BY createdAt DESC
	LIMIT 10
	`
	err := s.db.Select(&records, sqlQuery)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *service) QueryDayRecords() ([]LoadRecord, error) {
	records := []LoadRecord{}

	// Work aournd to query UTC dates from DB
	// TODO: improve logic
	year := time.Now().Local().Year()
	month := time.Now().Local().Month()
	day := time.Now().Local().Day()

	firstDate := time.Date(year, month, day, 3, 0, 0, 0, time.UTC)
	seccondDate := time.Date(year, month, day+1, 2, 59, 59, 999, time.UTC)

	sqlQuery := `
	SELECT * FROM load_record
	WHERE createdAt BETWEEN ? AND ?
	`
	err := s.db.Select(&records, sqlQuery, firstDate, seccondDate)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *service) QueryFilteredRecords(inputData FilteredInput) ([]LoadRecord, error) {
	records := []LoadRecord{}
	sqlQuery := `
	SELECT * 
	FROM load_record 
	WHERE (material = ? OR NULLIF(?, '') IS NULL)
	AND (createdAt BETWEEN ? AND ? OR (? IS NULL AND ? IS NULL))
	ORDER BY createdAt ASC
	`
	err := s.db.Select(&records, sqlQuery,
		inputData.Material, inputData.Material,
		inputData.FirstDate, inputData.SeccondDate,
		inputData.FirstDate, inputData.SeccondDate)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *service) CreateNewRecord(inputData RecordInput) error {
	sqlQuery := `
	INSERT INTO load_record (id, material, average_weight, unit, createdAt, operator, shift)
	VALUES(?, ?, ?, ?, COALESCE(?, NOW()),?, ?)
	`
	_, err := s.db.Exec(sqlQuery, uuid.New(), inputData.Material, inputData.AverageWeight, inputData.Unit, inputData.CreatedAt, inputData.Operator, inputData.Shift)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteRecord(inputData string) error {
	sqlQuery := `
	DELETE FROM load_record 
	WHERE id=?
	`
	_, err := s.db.Exec(sqlQuery, inputData)
	if err != nil {
		return err
	}
	return nil
}
