package processors

import (
	"context"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/processors/extract"
	"github.com/fiffu/mailtl/app/storage"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
)

type SaveInstaremCharge struct {
	infra.LogFacade
	storage.Storage
}

func NewSaveInstaremCharge(root infra.RootLogger, storage storage.Storage) (*SaveInstaremCharge, error) {
	return &SaveInstaremCharge{
		infra.NewLogger(root, "save_instarem_charge"),
		storage,
	}, nil
}

func (p *SaveInstaremCharge) Name() string { return "save_instarem_charge" }

func (p *SaveInstaremCharge) Initialize(_ backends.BackendConfig) error { return nil }

func (p *SaveInstaremCharge) Shutdown() error { return nil }

func (p *SaveInstaremCharge) HandleTaskSaveMail(ctx context.Context, e *mail.Envelope) (stopProcessing bool, err error) {
	extracter := extract.Instarem{Envelope: e}

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
