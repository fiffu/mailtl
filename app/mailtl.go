package app

import (
	"context"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/processors"
	"github.com/fiffu/mailtl/app/storage"
	"github.com/flashmob/go-guerrilla"
	"go.uber.org/fx"
)

// MailTL is the main MailTL app.
type MailTL struct{ guerrilla.Daemon }

func (m MailTL) onStart(ctx context.Context) error { return m.Start() }
func (m MailTL) onStop(ctx context.Context) error  { m.Shutdown(); return nil }
func (m MailTL) registerBackends(log infra.LogFacade, backends ...infra.Backend) {
	for _, backend := range backends {
		m.AddProcessor(
			backend.Name(),
			infra.MakeProcessor(log, backend),
		)
	}
}

func NewMailTL(
	lc fx.Lifecycle,
	daemon guerrilla.Daemon,
	store storage.Storage,
	root infra.RootLogger,
	senderFilter *processors.FilterBySender,
	saveInstaremCharge *processors.SaveInstaremCharge,
	saveDBSDigibankCharge *processors.SaveDBSDigibankCharge,
) (MailTL, error) {
	m := MailTL{daemon}
	m.registerBackends(
		infra.NewLogger(root, "backends"),
		senderFilter,
		saveInstaremCharge,
		saveDBSDigibankCharge,
	)

	hooks := fx.Hook{
		OnStart: m.onStart,
		OnStop:  m.onStop,
	}
	lc.Append(hooks)

	return m, nil
}
