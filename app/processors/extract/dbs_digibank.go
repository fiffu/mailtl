package extract

import (
	"fmt"
	"time"

	"github.com/fiffu/mailtl/app/model"
	"github.com/flashmob/go-guerrilla/mail"
	"golang.org/x/net/html"
)

var (
	dbsDigibankSender = mail.Address{
		User:        "ibanking.alert",
		Host:        "dbs.com",
		DisplayName: "",
	}
)

type dbsDigibankMatch struct {
	DateAndTime string
	LocalAmount string
	Instrument  string
	Purpose     string
}

func (d *dbsDigibankMatch) toModel(fingerprint string) (*model.Charge, error) {
	currency, amount, err := ParseMoney(d.LocalAmount)
	if err != nil {
		return nil, err
	}

	timestamp, err := d.parseTime()
	if err != nil {
		return nil, err
	}

	ret := &model.Charge{
		Fingerprint:    fingerprint,
		LocalCurrency:  currency,
		LocalAmount:    amount,
		Platform:       "dbs_digibank",
		Instrument:     d.Instrument,
		Timestamp:      timestamp,
		ChargeCurrency: currency,
		ChargeAmount:   amount,
		Purpose:        d.Purpose,
		IngestedAt:     time.Now(),
	}
	return ret, nil
}

func (d *dbsDigibankMatch) parseTime() (time.Time, error) {
	// 18 May 19:05 SGT 2025
	timestamp := fmt.Sprintf("%s %d", d.DateAndTime, time.Now().Year())
	layout := "02 Jan 15:04 MST 2006"
	return time.Parse(layout, timestamp)
}

type DBSDigibank struct{ *mail.Envelope }

func (i *DBSDigibank) Match() bool {
	return i.MailFrom.String() == dbsDigibankSender.String()
}

func (i *DBSDigibank) Extract() (charge *model.Charge, err error) {
	data := Markup(i.Envelope)

	texts, err := XPathEach(data, "//strong", func(n *html.Node) string {
		return n.NextSibling.Data
	})
	if err != nil {
		return nil, err
	}

	match := &dbsDigibankMatch{
		DateAndTime: texts[0],
		LocalAmount: texts[1],
		Instrument:  texts[2],
		Purpose:     texts[3],
	}
	return match.toModel(Fingerprint(i.Envelope))
}
