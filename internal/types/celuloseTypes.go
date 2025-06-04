package types

import "time"

type LoadFiltered struct {
	Material    string `json:"material"`
	FirstDate   string `json:"first_date"`
	SeccondDate string `json:"seccond_date"`
}

type LoadSummary struct {
	Material    string `json:"material"`
	TotalWeight int    `json:"total_weight"`
}

type Load struct {
	Id            string    `db:"id" json:"id"`
	Material      string    `db:"material" json:"material"`
	AverageWeight int       `db:"average_weight" json:"average_weight"`
	Unit          string    `db:"unit" json:"unit"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Timezone      string    `db:"timezone" json:"timezone"`
	Operator      string    `db:"operator" json:"operator"`
	Shift         string    `db:"shift" json:"shift"`
}
