package currency

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
)

type currencyStorage interface {
	Get(currency string) (float64, error)
}

type currencyNameStorage interface {
	Get(currency string) (name string, symbol string, err error)
}

type Bot struct {
	b                   *tgbotapi.BotAPI
	currencyStorage     currencyStorage
	currencyNameStorage currencyNameStorage
	UpdateTime          int
}

type Config struct {
	Token      string `json:"token"`
	TimeUpdate int    `json:"time_update"`
}

func NewBot(currencyStorage currencyStorage) *Bot {
	b, err := os.ReadFile("./config/tsconfig.json")
	if err != nil {
		return nil
	}
	cf := &Config{}
	err = json.Unmarshal(b, cf)
	if err != nil {
		return nil
	}
	bot, err := tgbotapi.NewBotAPI(cf.Token)
	if err != nil {
		return nil
	}
	bot.Debug = true

	return &Bot{currencyStorage: currencyStorage, b: bot, UpdateTime: cf.TimeUpdate}
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.UpdateTime
	ch := b.b.GetUpdatesChan(u)
	for update := range ch {
		go b.handleMsg(update)
	}

	return nil
}

type Message struct {
	Price    float64
	Symbol   string
	FullName string
}

func (m Message) String() string {
	return "Price: " + strconv.FormatFloat(m.Price, 'f', 2, 64) + "<br>" +
		"Symbol: " + m.Symbol + "<br>" +
		"Full name: " + m.FullName
}

func (b *Bot) handleMsg(update tgbotapi.Update) error {

	switch update.Message.Command() {
	case "p":
		price, err := b.currencyStorage.Get(update.Message.CommandArguments())
		if err != nil {
			_, err := b.b.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			return err
		}
		name, symbol, err := b.currencyNameStorage.Get(update.Message.CommandArguments())
		if err != nil {
			return err
		}
		message := Message{
			Price:    price,
			Symbol:   symbol,
			FullName: name, // Обновите это значение или добавьте логику для получения полного имени
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message.String())
		msg.ParseMode = "HTML"

		_, err = b.b.Send(msg)
		if err != nil {
			return err
		}

	}

	return nil
}
