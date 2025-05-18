package processors

import (
	"context"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/lib"
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
	sf.Infof(ctx, "* MailFrom: %v", e.MailFrom.String())
	sf.Infof(ctx, "* RemoteIP: %v", e.RemoteIP)
	sf.Infof(ctx, "* RcptTo:   %v", lib.StringsOf(lib.IndirectsOf(e.RcptTo)))
	sf.Infof(ctx, "* Subject:  %s", e.Subject)
}

func (sf *FilterBySender) isAllowedSender(mailFrom string) bool {
	return slices.Contains(sf.allowedSenders, mailFrom)
}

func (sf *FilterBySender) Handle(ctx context.Context, e *mail.Envelope) (continueProcessing bool, err error) {
	sf.logFields(ctx, e)

	if !sf.isAllowedSender(e.MailFrom.String()) {
		return false, nil
	}

	return true, nil
}
