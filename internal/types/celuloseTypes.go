package types

import "time"

type LoadFiltered struct {
	Material   string `json:"material"`
	FirstDate  string `json:"firstDate"`
	SecondDate string `json:"secondDate"`
}

type LoadSummary struct {
	Material    string `json:"material"`
	TotalWeight int    `json:"totalWeight"`
}

type Load struct {
	Id            string    `db:"id" json:"id"`
	Material      string    `db:"material" json:"material"`
	AverageWeight int       `db:"average_weight" json:"averageWeight"`
	Unit          string    `db:"unit" json:"unit"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	Timezone      string    `db:"timezone" json:"timezone"`
	Operator      string    `db:"operator" json:"operator"`
	Shift         string    `db:"shift" json:"shift"`
}
