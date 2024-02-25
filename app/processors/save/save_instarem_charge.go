package save

import (
	"encoding/json"
	"strings"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/lib/extraction"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/oriser/regroup"
)

const name = "save_instarem_charge"

var instaremSender = mail.Address{
	User:        "donotreply",
	Host:        "instarem.com",
	DisplayName: "Instarem notification",
}

var (
	contentXPath = "//span[2]"

	contentPattern = regroup.MustCompile(
		`We're notifying you that` +
			` (?P<localCurrency>[A-Z]{2,3}) (?P<localAmount>[\d\.]+)` +
			` was debited from your card (?P<cardNumber>[\dx-]+)` +
			` on (?P<timestamp>\d{4}-\d{2}\-\d{2} \d{2}:\d{2}:\d{2} GMT)` +
			` towards your transaction of` +
			` (?P<chargeCurrency>[A-Z]{2,3}) (?P<chargeAmount>[\d\.]+)` +
			` at (?P<purpose>.+?)\.$`,
	)
)

type matchData struct {
	LocalCurrency  string  `regroup:"localCurrency"`
	LocalAmount    float64 `regroup:"localAmount"`
	CardNumber     string  `regroup:"cardNumber"`
	Timestamp      string  `regroup:"timestamp"`
	ChargeCurrency string  `regroup:"chargeCurrency"`
	ChargeAmount   float64 `regroup:"chargeAmount"`
	Purpose        string  `regroup:"purpose"`
}

type SaveInstaremCharge struct{ infra.LogFacade }

func NewSaveInstaremCharge(root infra.RootLogger) (*SaveInstaremCharge, error) {
	return &SaveInstaremCharge{infra.NewLogger(root, name)}, nil
}

func (d *SaveInstaremCharge) Name() string { return name }

func (d *SaveInstaremCharge) Initialize(_ backends.BackendConfig) error { return nil }

func (d *SaveInstaremCharge) Shutdown() error { return nil }

func (d *SaveInstaremCharge) SaveMail(e *mail.Envelope) (continueProcessing bool, err error) {
	if e.MailFrom.String() != instaremSender.String() {
		d.Infof("Ignoring; expected: %s, got: %s)", instaremSender.String(), e.MailFrom.String())
		return true, nil
	}

	data := d.extractMarkup(e)

	text, err := extraction.ExtractXPath(data, contentXPath)
	if err != nil {
		d.Infof("Ignoring: %v", err)
		return true, nil
	}

	match, err := d.matchesPattern(text)
	if err != nil {
		d.Errorf(err, "No match in mail body: '%s'", text)
	}

	go d.save(match)
	return true, nil
}

func (d *SaveInstaremCharge) extractMarkup(e *mail.Envelope) string {
	// replace email delimiters
	data := e.Data.String()
	return strings.ReplaceAll(data, "=\n", "")
}

func (d *SaveInstaremCharge) matchesPattern(text string) (match matchData, err error) {
	err = contentPattern.MatchToTarget(text, &match)
	return
}

func (d *SaveInstaremCharge) save(match matchData) error {
	formatted, _ := json.MarshalIndent(match, "", "  ")
	d.Infof(string(formatted))
	return nil
}
