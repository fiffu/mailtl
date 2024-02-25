## sendml

**Step 1: Save an email from your mail client.**

On Gmail, you can do this by downloading an email. This will save a `.eml` file to your disk.


**Step 2: Ensure mailtl server is running.**

```sh
go run main.go
```

**Step 3: Send the saved email.**

Send to default local server port:

```sh
go run sendeml.go --eml ~/Downloads/your-saved-email.eml
```

You can use the `--server` option to specify a particular host/port:

```sh
go run sendeml.go --eml ~/Downloads/your-saved-email.eml --server localhost:2525
```
