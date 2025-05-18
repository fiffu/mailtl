package processors

import (
	"context"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/processors/extract"
	"github.com/fiffu/mailtl/app/storage"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
)

type SaveDBSDigibankCharge struct {
	infra.LogFacade
	storage.Storage
}

func NewSaveDBSDigibankCharge(root infra.RootLogger, storage storage.Storage) (*SaveDBSDigibankCharge, error) {
	return &SaveDBSDigibankCharge{
		infra.NewLogger(root, "save_dbs_digibank_charge"),
		storage,
	}, nil
}

func (p *SaveDBSDigibankCharge) Name() string { return "save_dbs_digibank_charge" }

func (p *SaveDBSDigibankCharge) Initialize(_ backends.BackendConfig) error { return nil }

func (p *SaveDBSDigibankCharge) Shutdown() error { return nil }

func (p *SaveDBSDigibankCharge) HandleTaskSaveMail(ctx context.Context, e *mail.Envelope) (stopProcessing bool, err error) {
	extracter := extract.DBSDigibank{Envelope: e}

	if !extracter.Match() {
		p.Debugf(ctx, "Ignoring due to mismatch")
		return
	}

	charge, err := extracter.Extract()
	if err != nil {
		p.Errorf(ctx, err, "Extraction error")
		err = nil
		return
	}

	err = charge.Save(ctx, p.Storage)
	return
}
