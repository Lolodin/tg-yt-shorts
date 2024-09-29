package converter

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"regexp"
)

type Bot struct {
	b          *tgbotapi.BotAPI
	convert    Converter
	UpdateTime int
}

func NewBot(convert Converter) *Bot {
	b, err := os.ReadFile("./root/tsconfig.json")
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

	return &Bot{convert: convert, b: bot, UpdateTime: cf.TimeUpdate}
}

type Config struct {
	Token      string `json:"token"`
	TimeUpdate int    `json:"time_update"`
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

func (b *Bot) handleMsg(update tgbotapi.Update) error {

	by, err := b.convert.Covert(update.Message.Text)
	if err != nil {
		url, e := ExtractFirstURL(err.Error())
		if e != nil {
			_, err = b.b.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			return err
		}
		_, err := b.b.Send(tgbotapi.NewMessage(update.Message.Chat.ID, url))

		return err
	}
	vcf := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: by,
		Name: fmt.Sprintf("%p", by)})
	_, err = b.b.Send(vcf)
	fmt.Println(err)
	return err

}

func ExtractFirstURL(text string) (string, error) {
	// Регулярное выражение для поиска URL
	re := regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

	// Находим все совпадения
	matches := re.FindStringSubmatch(text)

	// Если совпадений нет, возвращаем ошибку
	if len(matches) == 0 {
		return "", fmt.Errorf("URL not found")
	}

	// Возвращаем первое совпадение
	return matches[0], nil
}
