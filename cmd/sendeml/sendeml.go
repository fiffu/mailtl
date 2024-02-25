package main

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/urfave/cli/v2"
)

const (
	flagServer  = "server"
	flagEMLPath = "eml"
)

type Args struct {
	EML struct {
		Path string
		File parsemail.Email
	}
	SMTPAddr string // mail.example.com:port
}

func parseArgs(ctx *cli.Context) (args Args, err error) {
	addr, err := parseServerAddr(ctx.String(flagServer))
	if err != nil {
		return
	}

	path := ctx.String(flagEMLPath)
	eml, err := parseEML(path)
	if err != nil {
		return
	}

	args.SMTPAddr = addr
	args.EML.Path = path
	args.EML.File = eml
	return
}

func parseServerAddr(s string) (string, error) {
	substrs := strings.Split(s, ":")
	if len(substrs) < 2 {
		return "", fmt.Errorf("server address should be in format hostname:port")
	}

	host := substrs[0]
	port := substrs[1]
	if host == "" {
		return "", fmt.Errorf("missing host component in '%s'", s)
	}
	if port == "" {
		port = "25"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

func parseEML(emlPath string) (_ parsemail.Email, err error) {
	emlFile, err := os.Open(emlPath)
	if err != nil {
		return
	}
	return parsemail.Parse(emlFile)
}

func sendEML(ctx *cli.Context) error {
	args, err := parseArgs(ctx)
	if err != nil {
		return err
	}

	if len(args.EML.File.From) == 0 {
		return fmt.Errorf("%s needs to contain at least one address in the 'From' header", args.EML.Path)
	}

	return smtp.SendMail(
		args.SMTPAddr,
		nil,
		stringifyAddrs(args.EML.File.From)[0],
		stringifyAddrs(args.EML.File.To),
		[]byte(args.EML.File.HTMLBody),
	)
}

func stringifyAddrs(addrs []*mail.Address) []string {
	strs := make([]string, len(addrs))
	for i, addr := range addrs {
		strs[i] = addr.Address
	}
	return strs
}

func main() {
	app := &cli.App{
		Name:        "sendeml",
		Description: "sends an .eml file to the given destination",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flagServer,
				Value: "localhost:2525",
			},
			&cli.StringFlag{
				Name:     flagEMLPath,
				Usage:    "path to .eml file to be sent",
				Required: true,
			},
		},

		Action: sendEML,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
