## mailtl

mailtl = email + transform + load


#### Requirements

- Go 1.20


#### Development

Run this to setup the development env:
```sh
make
```

This will
- install dependencies
- install dev tooling (`mockery`, `gotestsum`)
- setup git hooks
- run unit tests


#### Run

To run server with default config:
```sh
go run main.go
```

To run with config file:
```sh
go run main.go sample.config.json
```

#### Processors

* `filter_by_sender` - Prints some headers and short-circuits if sender is not allowed in config.
* `save_instarem_charge` - Saves email notifications for Instarem charges.


## cmd/sendeml

This command allows you to send emails to the mailtl server.

Usage: `go run sendeml.go --eml ~/Downloads/your-saved-email.eml --server localhost:2525`

See [sendeml README](./cmd/sendeml/README.md) for details.
