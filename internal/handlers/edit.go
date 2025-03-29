package handlers

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type editor struct {
	ctx     *th.Context
	params  *telego.LinkPreviewOptions
	chat_id int64
}

func NewEditor(ctx *th.Context, chat_id int64) *editor {
	if ctx == nil {
		return nil
	}

	return &editor{ctx: ctx, chat_id: chat_id, params: &telego.LinkPreviewOptions{IsDisabled: true}}
}

func (e *editor) Edit(message_id int64, text string) error {
	_, err := e.ctx.Bot().EditMessageText(e.ctx, &telego.EditMessageTextParams{ChatID: tu.ID(e.chat_id), MessageID: int(message_id), LinkPreviewOptions: e.params, Text: text})
	return err
}

func (e *editor) EditHTML(message_id int64, text string) error {
	_, err := e.ctx.Bot().EditMessageText(e.ctx, &telego.EditMessageTextParams{ChatID: tu.ID(e.chat_id), MessageID: int(message_id), LinkPreviewOptions: e.params, Text: text, ParseMode: telego.ModeHTML})
	return err
}

func (e *editor) EditMD(message_id int64, text string) error {
	_, err := e.ctx.Bot().EditMessageText(e.ctx, &telego.EditMessageTextParams{ChatID: tu.ID(e.chat_id), MessageID: int(message_id), LinkPreviewOptions: e.params, Text: text, ParseMode: telego.ModeMarkdownV2})
	return err
}
