package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nktauserum/aisearch-telegram/pkg/parse"
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

	msg_params := tu.Message(tu.ID(message.Chat.ID), "🔎 <i>Searching...</i>")
	msg_params.ParseMode = telego.ModeHTML
	msg_params.LinkPreviewOptions = &telego.LinkPreviewOptions{IsDisabled: true}

	editor := NewEditor(ctx, message.Chat.ID)

	log.Printf("search request from @%s request: %s\n", message.From.Username, message.Text)

	msg, err := ctx.Bot().SendMessage(ctx, msg_params)
	if err != nil {
		log.Printf("error while sending a message: %v", err)
		return err
	}

	result := stream.Request(message.Text)

	var currentText string
	var buffer string
	for res := range result {
		if res.Event == stream.Message {
			response := Response{}
			err := json.Unmarshal([]byte(res.Content), &response)
			if err != nil {
				_ = editor.Edit(int64(msg.MessageID), fmt.Sprintf("We're sorry, but json response from search service didn't want unmarshalling. Please contact the developer: @Auserum\n\n``` json\n%s\n```", res.Content))
				log.Printf("error while unmarshalling a response json: %v", err)
				return err
			}

			// Ответ не содержит полезной информации
			if len(response.Content) == 0 && len(currentText) == 0 {
				_ = editor.Edit(int64(msg.MessageID), fmt.Sprintf("We're sorry, but the model did not give us a correct answer. Please contact the developer: @Auserum\n\nNow answer is: %v", response))
				log.Printf("error getting a model's answer: %v", err)
				return err
			}

			buffer += response.Content
			// Ограничение символов требуется для API Telegram
			if len(buffer) > 90 {
				currentText += buffer
				err = editor.Edit(int64(msg.MessageID), currentText)
				if err != nil {
					log.Printf("error editing message: %v", err)
					return err
				}
				buffer = ""
			}
		} else if res.Event == stream.Source {
			currentText += fmt.Sprintf("\n%s", res.Content)
			_ = editor.Edit(int64(msg.MessageID), currentText)
		}
	}

	// Если после обработки полученных сообщений у нас остался буфер, добавляем его к сообщению
	if len(buffer) > 0 {
		currentText += buffer
		err := editor.Edit(int64(msg.MessageID), currentText)
		if err != nil {
			log.Printf("error editing message: %v", err)
			return err
		}
	}

	parsed_content, err := parse.ConvertToHTML(currentText)
	if err != nil {
		log.Printf("error parsing final message: %v", err)
		return err
	}
	parsed_content = strings.TrimPrefix(parsed_content, "<html><head></head><body>")
	parsed_content = strings.TrimSuffix(parsed_content, "</body></html>")

	err = editor.EditHTML(int64(msg.MessageID), parsed_content)
	if err != nil {
		log.Printf("error editing final message: %v", err)
		return err
	}

	return nil
}
