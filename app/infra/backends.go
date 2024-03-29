package infra

import (
	"context"
	"fmt"

	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/flashmob/go-guerrilla/response"
)

type customCtxKey struct{}

var ctxKey = customCtxKey{}

type CustomContext struct {
	Fingerprint string
}

func fingerprint(envelope *mail.Envelope) string {
	id := envelope.QueuedId
	from := envelope.MailFrom.String()
	subject := envelope.Subject
	return fmt.Sprintf("<Envelope %s from='%s' subj='%s'>", id, from, subject)
}

func deriveContext(envelope *mail.Envelope) (context.Context, error) {
	if envelope.Header == nil {
		if err := envelope.ParseHeaders(); err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey, CustomContext{
		Fingerprint: fingerprint(envelope),
	})
	return ctx, nil
}

//go:generate mockery --name Backend
type Backend interface {
	Name() string
	Initialize(backendConfig backends.BackendConfig) error
	Shutdown() error
	SaveMail(ctx context.Context, e *mail.Envelope) (continueProcessing bool, err error)
}

// FixtureBackend embeds backends.Processor. This is intended for testing purposes only.
//
//go:generate mockery --name FixtureBackend
type FixtureBackend interface{ backends.Processor }

var NoopAndContinue = func() (backends.Result, error, bool) { return nil, nil, true }

func MakeProcessor(b Backend) backends.ProcessorConstructor {
	backends.Svc.AddInitializer(b)
	backends.Svc.AddShutdowner(b)

	return func() backends.Decorator {
		return func(next backends.Processor) backends.Processor {
			return backends.ProcessWith(func(e *mail.Envelope, task backends.SelectTask) (res backends.Result, err error) {
				ctx, err := deriveContext(e)
				if err != nil {
					ctx = context.Background()
				}

				var continueProcessing bool

				switch task {
				// case backends.TaskValidateRcpt:
				// 	continueProcessing, err = b.Validate(e)

				case backends.TaskSaveMail:
					// We only handle save tasks to simplify (just use a single call chain along "save_process")
					continueProcessing, err = b.SaveMail(ctx, e)

				default:
					continueProcessing = false
				}

				if err != nil {
					res = backends.NewResult(fmt.Sprintf("554 Error: %s", err))
					return
				}
				if !continueProcessing {
					res = backends.NewResult(response.Canned.SuccessNoopCmd)
					return
				}

				return next.Process(e, task)
			})
		}
	}
}
