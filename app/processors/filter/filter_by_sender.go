package filter

import (
	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/lib"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"golang.org/x/exp/slices"
)

const name = "filter_by_sender"

type FilterBySender struct {
	infra.LogFacade
	allowedSenders []string
}

func NewFilterBySender(root infra.RootLogger, rootConfig infra.RootConfig) (*FilterBySender, error) {
	return &FilterBySender{
		LogFacade:      infra.NewLogger(root, name),
		allowedSenders: rootConfig.AllowedSenders,
	}, nil
}

func (sf *FilterBySender) Name() string { return name }

func (sf *FilterBySender) Initialize(_ backends.BackendConfig) error { return nil }

func (sf *FilterBySender) Shutdown() error { return nil }

func (sf *FilterBySender) logFields(e *mail.Envelope) {
	sf.Infof("* MailFrom: %v", e.MailFrom.String())
	sf.Infof("* RemoteIP: %v", e.RemoteIP)
	sf.Infof("* RcptTo:   %v", lib.StringsOf(lib.IndirectsOf(e.RcptTo)))
	sf.Infof("* Subject:  %s", e.Subject)
}

func (sf *FilterBySender) isAllowedSender(mailFrom string) bool {
	return slices.Contains(sf.allowedSenders, mailFrom)
}

func (sf *FilterBySender) SaveMail(e *mail.Envelope) (continueProcessing bool, err error) {
	sf.logFields(e)

	if !sf.isAllowedSender(e.MailFrom.String()) {
		return false, nil
	}

	return true, nil
}
