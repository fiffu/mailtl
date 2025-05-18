package processors

import (
	"context"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"golang.org/x/exp/slices"
)

type FilterBySender struct {
	infra.LogFacade
	allowedSenders []string
}

func NewFilterBySender(root infra.RootLogger, rootConfig infra.RootConfig) (*FilterBySender, error) {
	return &FilterBySender{
		LogFacade:      infra.NewLogger(root, "filter_by_sender"),
		allowedSenders: rootConfig.AllowedSenders,
	}, nil
}

func (sf *FilterBySender) Name() string { return "filter_by_sender" }

func (sf *FilterBySender) Initialize(_ backends.BackendConfig) error { return nil }

func (sf *FilterBySender) Shutdown() error { return nil }

func (sf *FilterBySender) logFields(ctx context.Context, e *mail.Envelope) {

}

func (sf *FilterBySender) isAllowedSender(mailFrom string) bool {
	return slices.Contains(sf.allowedSenders, mailFrom)
}

func (sf *FilterBySender) HandleTaskSaveMail(ctx context.Context, e *mail.Envelope) (stopProcessing bool, err error) {
	if sf.isAllowedSender(e.MailFrom.String()) {
		sf.Infof(ctx, "Mail from '%v': %s", e.MailFrom.String(), e.Subject)
	} else {
		sf.Debugf(ctx, "Ignoring mail from '%v': %s", e.MailFrom.String(), e.Subject)
		stopProcessing = true
	}
	return
}
