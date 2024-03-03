package testfixtures

import (
	"bytes"

	"github.com/flashmob/go-guerrilla/mail"
)

func NewEnvelope() *mail.Envelope {
	return &mail.Envelope{
		MailFrom: NewAddress(),
		Data:     NewData(),
	}
}

func NewAddress() mail.Address {
	return mail.Address{
		User: "user",
		Host: "example.com",
	}
}

var headerData = "Subject: testfixtures.go"
var headerEnd = "\n\n"

func NewData() bytes.Buffer {
	return *bytes.NewBuffer([]byte(headerData + headerEnd))
}
