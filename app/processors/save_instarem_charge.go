package processors

import (
	"context"
	"strings"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/model"
	"github.com/fiffu/mailtl/app/storage"
	"github.com/fiffu/mailtl/lib/extraction"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/oriser/regroup"
)

var instaremSender = mail.Address{
	User:        "donotreply",
	Host:        "instarem.com",
	DisplayName: "Instarem notification",
}

var (
	instaremContentXPath = "//span[2]"

	instaremContentPattern = regroup.MustCompile(
		`We're notifying you that` +
			` (?P<localCurrency>[A-Z]{2,3}) (?P<localAmount>[\d\.]+)` +
			` was debited from your card (?P<cardNumber>[\dx-]+)` +
			` on (?P<timestamp>\d{4}-\d{2}\-\d{2} \d{2}:\d{2}:\d{2} GMT)` +
			` towards your transaction of` +
			` (?P<chargeCurrency>[A-Z]{2,3}) (?P<chargeAmount>[\d\.]+)` +
			` at (?P<purpose>.+?)\.$`,
	)
)

type instaremMatchData struct {
	LocalCurrency  string  `regroup:"localCurrency"`
	LocalAmount    float64 `regroup:"localAmount"`
	CardNumber     string  `regroup:"cardNumber"`
	Timestamp      string  `regroup:"timestamp"`
	ChargeCurrency string  `regroup:"chargeCurrency"`
	ChargeAmount   float64 `regroup:"chargeAmount"`
	Purpose        string  `regroup:"purpose"`
}

func (i *instaremMatchData) toModel() *model.Charge {
	return &model.Charge{
		LocalCurrency:  i.LocalCurrency,
		LocalAmount:    i.LocalAmount,
		CardNumber:     i.CardNumber,
		Timestamp:      i.Timestamp,
		ChargeCurrency: i.ChargeCurrency,
		ChargeAmount:   i.ChargeAmount,
		Purpose:        i.Purpose,
	}
}

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
	if e.MailFrom.String() != instaremSender.String() {
		p.Debugf(ctx, "Ignoring: expected: %s, got: %s)", instaremSender.String(), e.MailFrom.String())
		return
	}

	data := p.extractMarkup(e)

	text, err := extraction.ExtractXPath(data, instaremContentXPath)
	if err != nil {
		p.Debugf(ctx, "Ignoring: %v", err)
		return
	}

	match, err := p.matchesPattern(text)
	if err != nil {
		p.Debugf(ctx, "Ignoring: no match in email body: '%s', err: %v", text, err)
		return
	}

	err = match.toModel().Save(ctx, p.Storage)
	return
}

func (p *SaveInstaremCharge) extractMarkup(e *mail.Envelope) string {
	// replace email delimiters
	data := e.Data.String()
	return strings.ReplaceAll(data, "=\n", "")
}

func (p *SaveInstaremCharge) matchesPattern(text string) (match instaremMatchData, err error) {
	err = instaremContentPattern.MatchToTarget(text, &match)
	return
}
