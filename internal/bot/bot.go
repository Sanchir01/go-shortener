package tgbot

import (
	"context"
	"log/slog"
	"os"

	"github.com/Sanchir01/go-shortener/internal/feature/user"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

type TGBot struct {
	Bot  *bot.Bot
	url  UrlCreator
	user UserService
	l    *slog.Logger
}

type UrlCreator interface {
	CreateUrl(ctx context.Context, userId uuid.UUID, url string) error
}
type UserService interface {
	Register(ctx context.Context, p user.RegisterParams) (*uuid.UUID, error)
}

func New(ctx context.Context, url UrlCreator, user UserService, l *slog.Logger) (*TGBot, error) {
	t := &TGBot{url: url, user: user, l: l}
	opts := []bot.Option{
		bot.WithDefaultHandler(t.UnknownCommand),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		return nil, err
	}

	t.Bot = b

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, t.Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, t.Help)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/ping", bot.MatchTypeExact, t.Ping)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/short", bot.MatchTypePrefix, t.Short)

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Начать"},
			{Command: "help", Description: "Помощь"},
			{Command: "ping", Description: "Проверка связи"},
			{Command: "short", Description: "Сократить ссылку: /short <url>"},
		},
		Scope: &models.BotCommandScopeDefault{},
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TGBot) Start(ctx context.Context, b *bot.Bot, update *models.Update) {

	user, err := t.user.Register(ctx, user.RegisterParams{
		TGID:  &update.Message.From.ID,
		Title: update.Message.From.Username,
	})
	if err != nil {
		t.l.Error("failed register error")
		t.Bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Нажмите комманду /start еще раз произошла ошибка",
		})
	}
	t.l.Info("user register by tg", user)
	t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет! Доступные команды:\n/start — Начать\n/help — Помощь\n/ping — Проверка\n/short <url> — Сократить ссылку",
	})
}

func (t *TGBot) Help(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Помощь:\n/short <url> — отправьте ссылку, чтобы получить короткий вариант.",
	})
}

func (t *TGBot) Ping(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "pong",
	})
}

func (t *TGBot) Short(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	if len(text) <= len("/short ") {
		t.Bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Использование: /short <url>",
		})
		return
	}

	url := text[len("/short "):]
	if t.url != nil {
		_ = t.url.CreateUrl(ctx, uuid.New(), url)
	}

	t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Принял ссылку: " + url,
	})
}

func (t *TGBot) UnknownCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Неизвестная команда. Введите /help",
	})
}
