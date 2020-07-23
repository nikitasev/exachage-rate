package entity

import (
	_ "github.com/go-sql-driver/mysql"
)

type Ticks struct {
	Timestamp int64   `db:"timestamp"`
	Symbol    string  `db:"symbol"`
	Bid       float64 `db:"bid"`
	Ask       float64 `db:"ask"`
}

func (t *Ticks) TableName() string {
	return "ticks"
}

func (t *Ticks) PK() string {
	return "id"
}

