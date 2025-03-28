package handlers

import (
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Start(ctx *th.Context, update telego.Update) error {
	chatID := tu.ID(update.Message.Chat.ID)

	message := tu.Message(chatID, "Hello! I'm a bot that can help you find information in the Internet with AI. Just type what you want to know and I'll find it for you.")
	message.ParseMode = telego.ModeHTML

	log.Printf("start: @%s\n", update.Message.From.Username)

	if _, err := ctx.Bot().SendMessage(ctx, message); err != nil {
		log.Printf("error while sending message: %v", err)
		return err
	}

	return nil
}
