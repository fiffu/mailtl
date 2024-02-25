package infra

import (
	"fmt"

	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/flashmob/go-guerrilla/response"
)

type Backend interface {
	Name() string
	Initialize(backendConfig backends.BackendConfig) error
	Shutdown() error
	SaveMail(e *mail.Envelope) (continueProcessing bool, err error)
}

var NoopAndContinue = func() (backends.Result, error, bool) { return nil, nil, true }

func MakeProcessor(b Backend) backends.ProcessorConstructor {
	backends.Svc.AddInitializer(b)
	backends.Svc.AddShutdowner(b)

	return func() backends.Decorator {
		return func(next backends.Processor) backends.Processor {
			return backends.ProcessWith(func(e *mail.Envelope, task backends.SelectTask) (res backends.Result, err error) {
				var continueProcessing bool

				switch task {
				// case backends.TaskValidateRcpt:
				// 	continueProcessing, err = b.Validate(e)

				case backends.TaskSaveMail:
					// We only handle save tasks to simplify (just use a single call chain along "save_process")
					continueProcessing, err = b.SaveMail(e)

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
