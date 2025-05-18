package main

import (
	"github.com/fiffu/mailtl/app"
	"github.com/fiffu/mailtl/app/guerrillad"
	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/processors"
	"github.com/fiffu/mailtl/app/storage"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Guerrilla daemon and processors
		fx.Provide(guerrillad.NewGuerillaDaemon),
		fx.Provide(processors.NewFilterBySender),
		fx.Provide(processors.NewSaveInstaremCharge),
		fx.Provide(processors.NewSaveDBSDigibankCharge),

		// Infra glue
		fx.Provide(infra.NewRootConfig),
		fx.Provide(infra.NewRootLogger),
		fx.Provide(storage.NewStorage),

		fx.Provide(app.NewMailTL),
		fx.Invoke(func(app.MailTL) {}),
	).Run()
}
