package main

import (
	"github.com/fiffu/mailtl/app"
	"github.com/fiffu/mailtl/app/guerrillad"
	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/processors/filter"
	"github.com/fiffu/mailtl/app/processors/save"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Guerrilla daemon and processors
		fx.Provide(guerrillad.NewGuerillaDaemon),
		fx.Provide(filter.NewFilterBySender),
		fx.Provide(save.NewSaveInstaremCharge),

		// Infra glue
		fx.Provide(infra.NewRootConfig),
		fx.Provide(infra.NewRootLogger),
		fx.Provide(app.NewMailTL),

		fx.Invoke(func(app.MailTL) {}),
	).Run()
}
