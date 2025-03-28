package app

import (
	"context"
	"fmt"
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/nktauserum/aisearch-telegram/internal/handlers"
)

type Application struct {
	token string
}

func NewApplication(token string) *Application {
	return &Application{token: token}
}

func NewBot(token string) (*telego.Bot, error) {
	return telego.NewBot(token, telego.WithDefaultLogger(false, true))
}

func (a *Application) Run() error {
	ctx := context.Background()
	log.Println("Starting the bot...")

	bot, err := NewBot(a.token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)
	bh, _ := th.NewBotHandler(bot, updates)
	bh.Handle(handlers.Start, th.CommandEqual("start"))
	bh.HandleMessage(handlers.Search)

	defer func() { _ = bh.Stop() }()
	return bh.Start()
}
