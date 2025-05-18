package extract

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/fiffu/mailtl/app/model"
	"github.com/fiffu/mailtl/lib"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/oriser/regroup"
)

var (
	errUnknownParserVersion = errors.New("unknown parser version")

	instaremSender = mail.Address{
		User:        "donotreply",
		Host:        "instarem.com",
		DisplayName: "Instarem notification",
	}
	instaremV1Pattern = regroup.MustCompile(
		`We're notifying you that` +
			` (?P<localCurrency>[A-Z]{2,3}) (?P<localAmount>[\d\.]+)` +
			` was debited from your card (?P<instrument>[\dx-]+)` +
			` on (?P<timestamp>\d{4}-\d{2}\-\d{2} \d{2}:\d{2}:\d{2} [A-Z]{3})` +
			` towards your transaction of` +
			` (?P<chargeCurrency>[A-Z]{2,3}) (?P<chargeAmount>[\d\.]+)` +
			` at (?P<purpose>.+?)\.$`,
	)
	instaremV2Pattern = regroup.MustCompile(
		`TransactionAmount____((?P<chargeAmount>[\d\.]+) (?P<chargeCurrency>[A-Z]{2,3}))?` +
			`____AmountPaid____(?P<localAmount>[\d\.]+) (?P<localCurrency>[A-Z]{2,3})` +
			`____Merchant____(?P<purpose>.+)` +
			`____DateAndTime____(?P<timestamp>\d{1,2}(st|nd|rd|th) \w{3}, \d{4} .+[A-Z]{3})` +
			`____PaymentSource____(?P<instrument>.+)`)
	dateOrdinalPrefixPattern = regexp.MustCompile(`^(\d+)(st|nd|rd|th)`)
)

type instaremMatch struct {
	LocalCurrency  string  `regroup:"localCurrency"`
	LocalAmount    float64 `regroup:"localAmount"`
	Instrument     string  `regroup:"instrument"`
	Timestamp      string  `regroup:"timestamp"`
	ChargeCurrency string  `regroup:"chargeCurrency"`
	ChargeAmount   float64 `regroup:"chargeAmount"`
	Purpose        string  `regroup:"purpose"`
}

func (i *instaremMatch) toModel(version int, fingerprint string) (*model.Charge, error) {
	timestamp, err := i.parseTimestamp(version)
	if err != nil {
		return nil, err
	}

	ret := &model.Charge{
		Fingerprint:    fingerprint,
		LocalCurrency:  i.LocalCurrency,
		LocalAmount:    i.LocalAmount,
		Platform:       "instarem",
		Instrument:     i.Instrument,
		Timestamp:      timestamp,
		ChargeCurrency: i.ChargeCurrency,
		ChargeAmount:   i.ChargeAmount,
		Purpose:        i.Purpose,
		IngestedAt:     time.Now(),
	}
	if ret.ChargeCurrency == "" && ret.ChargeAmount == 0 {
		ret.ChargeCurrency = ret.LocalCurrency
		ret.ChargeAmount = ret.LocalAmount
	}
	return ret, nil
}

func (i *instaremMatch) parseTimestamp(version int) (time.Time, error) {
	var ts, layout string
	ts = i.Timestamp

	switch version {
	case 1:
		layout = "2006-01-02 15:04:05 MST"
	case 2:
		layout = "2 Jan, 2006 3:04 PM MST"
		ts = dateOrdinalPrefixPattern.ReplaceAllString(i.Timestamp, "${1}")
	default:
		return time.Time{}, errUnknownParserVersion
	}

	return time.Parse(layout, ts)
}

type Instarem struct{ *mail.Envelope }

func (i *Instarem) Match() bool {
	return i.MailFrom.String() == instaremSender.String()
}

func (i *Instarem) Extract() (charge *model.Charge, err error) {
	charge, err = i.extractV2()
	if err == nil {
		return
	}
	charge, err = i.extractV1()
	return
}

func (i *Instarem) extractV1() (*model.Charge, error) {
	data := Markup(i.Envelope)

	text, err := XPath(data, "//span[2]")
	if err != nil {
		return nil, err
	}

	var match instaremMatch
	if err := instaremV1Pattern.MatchToTarget(text, &match); err != nil {
		return nil, err
	}

	return match.toModel(1, Fingerprint(i.Envelope))
}

func (i *Instarem) extractV2() (*model.Charge, error) {
	table := map[string]string{}
	table["//table[3]//table//table//tr[1]//td[1]"] = "//table[3]//table//table//tr[1]//td[2]"
	table["//table[3]//table//table//tr[2]//td[1]"] = "//table[3]//table//table//tr[2]//td[2]"
	table["//table[3]//table//table//tr[3]//td[1]"] = "//table[3]//table//table//tr[3]//td[2]"
	table["//table[3]//table//table//tr[4]//td[1]"] = "//table[3]//table//table//tr[4]//td[2]"
	table["//table[3]//table//table//tr[5]//td[1]"] = "//table[3]//table//table//tr[5]//td[2]"
	return i.extractV2Table(table)
}

func (i *Instarem) extractV2Table(table map[string]string) (*model.Charge, error) {
	data := Markup(i.Envelope)

	for k, v := range table {
		left := strings.TrimSpace(lib.DropError(XPath(data, k)))
		right := strings.TrimSpace(lib.DropError(XPath(data, v)))

		delete(table, k)
		table[left] = right
	}

	text := strings.Join(
		[]string{
			"TransactionAmount", table["Transaction amount"],
			"AmountPaid", table["Amount paid"],
			"Merchant", table["Merchant"],
			"DateAndTime", table["Date, time"],
			"PaymentSource", table["Payment source"],
		},
		"____",
	)

	var match instaremMatch
	if err := instaremV2Pattern.MatchToTarget(text, &match); err != nil {
		return nil, err
	}

	return match.toModel(2, Fingerprint(i.Envelope))
}
