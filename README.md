## mailtl

mailtl = email + transform + load


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
