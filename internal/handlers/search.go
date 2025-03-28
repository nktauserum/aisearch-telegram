package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nktauserum/aisearch-telegram/pkg/stream"
)

type Response struct {
	Content string `json:"content"`
}

func Search(ctx *th.Context, message telego.Message) error {
	status := stream.Ping()
	if !status {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(message.Chat.ID), "We're sorry, but search service is inaccessible. Please contact the developer: @Auserum"))
		return fmt.Errorf("search service inaccessible")
	}

	result := stream.Request(message.Text)
	log.Printf("search: @%s\n", message.From.Username)

	msg_params := tu.Message(tu.ID(message.Chat.ID), "Searching...")
	msg_params.ParseMode = telego.ModeHTML

	msg, err := ctx.Bot().SendMessage(ctx, msg_params)
	if err != nil {
		log.Printf("error while sending a message: %v", err)
		return err
	}

	currentText := ""
	for res := range result {
		if res.Event == stream.Message {
			response := Response{}
			err := json.Unmarshal([]byte(res.Content), &response)
			if err != nil {
				_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(tu.ID(message.Chat.ID), msg.MessageID, fmt.Sprintf("We're sorry, but json response from search service didn't want unmarshalling. Please contact the developer: @Auserum\n\n``` json\n%s\n```", res.Content)))
				log.Printf("error while unmarshalling a response json: %v", err)
				return err
			}

			if len(response.Content) == 0 && len(currentText) == 0 {
				_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(tu.ID(message.Chat.ID), msg.MessageID, fmt.Sprintf("We're sorry, but the model did not give us a correct answer. Please contact the developer: @Auserum\n\nNow answer is: %v", response)))
				log.Printf("error getting a model's answer: %v", err)
				return err
			}

			currentText = currentText + response.Content
			ctx.Bot().EditMessageText(ctx, tu.EditMessageText(tu.ID(msg.Chat.ID), msg.MessageID, currentText))
		} else if res.Event == stream.Source {
			currentText = fmt.Sprintf("%s\n%s", currentText, res.Content)
			ctx.Bot().EditMessageText(ctx, tu.EditMessageText(tu.ID(msg.Chat.ID), msg.MessageID, currentText))
		}
	}

	return nil
}
