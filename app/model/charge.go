package model

import (
	"context"
	"time"

	"github.com/fiffu/mailtl/app/storage"
	"github.com/fiffu/mailtl/lib"
)

type Charge struct {
	Fingerprint    string
	LocalCurrency  string
	LocalAmount    float64
	Platform       string
	Instrument     string
	Timestamp      time.Time
	ChargeCurrency string
	ChargeAmount   float64
	Purpose        string
	IngestedAt     time.Time
}

func (c *Charge) Save(ctx context.Context, storage storage.Storage) error {
	return lib.DropResult(storage.DB().ExecContext(
		ctx,
		`
			INSERT INTO charges (
				fingerprint,
				timestamp,
				platform,
				instrument,
				local_currency,
				local_amount,
				charge_currency,
				charge_amount,
				purpose,
				ingested_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		c.Fingerprint,
		c.Timestamp.Unix(),
		c.Platform,
		c.Instrument,
		c.LocalCurrency,
		c.LocalAmount,
		c.ChargeCurrency,
		c.ChargeAmount,
		c.Purpose,
		c.IngestedAt.Unix(),
	))
}
