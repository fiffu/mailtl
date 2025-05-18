package model

import (
	"context"

	"github.com/fiffu/mailtl/app/storage"
	"github.com/fiffu/mailtl/lib"
)

type Charge struct {
	LocalCurrency  string
	LocalAmount    float64
	CardNumber     string
	Timestamp      string
	ChargeCurrency string
	ChargeAmount   float64
	Purpose        string
}

func (c *Charge) Save(ctx context.Context, storage storage.Storage) error {
	return lib.Just(storage.DB().ExecContext(
		ctx,
		`
			INSERT INTO charges (
				local_currency,
				local_amount,
				card_number,
				timestamp,
				charge_currency,
				charge_amount,
				purpose
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		c.LocalCurrency,
		c.LocalAmount,
		c.CardNumber,
		c.Timestamp,
		c.ChargeCurrency,
		c.ChargeAmount,
		c.Purpose,
	))
}
