package app

import (
	"context"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/processors/filter"
	"github.com/fiffu/mailtl/app/processors/save"
	"github.com/flashmob/go-guerrilla"
	"go.uber.org/fx"
)

// MailTL is the main MailTL app.
type MailTL struct{ guerrilla.Daemon }

func (m MailTL) onStart(ctx context.Context) error { return m.Start() }
func (m MailTL) onStop(ctx context.Context) error  { m.Shutdown(); return nil }
func (m MailTL) registerBackends(backends ...infra.Backend) {
	for _, backend := range backends {
		m.AddProcessor(
			backend.Name(),
			infra.MakeProcessor(backend),
		)
	}
}

func NewMailTL(
	lc fx.Lifecycle,
	daemon guerrilla.Daemon,
	senderFilter *filter.FilterBySender,
	saveInstaremCharge *save.SaveInstaremCharge,
) (MailTL, error) {
	m := MailTL{daemon}
	m.registerBackends(
		senderFilter,
		saveInstaremCharge,
	)

	hooks := fx.Hook{
		OnStart: m.onStart,
		OnStop:  m.onStop,
	}
	lc.Append(hooks)

	return m, nil
}
