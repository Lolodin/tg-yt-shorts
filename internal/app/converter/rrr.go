package converter

import "errors"

type RRR interface {
	Register() error
	Run() error
	Resolve() error
}

type Telegram interface {
	Run() error
}
type App struct {
	telegram Telegram
}

func (a *App) Register() error {
	var c Converter
	c = NewVideo()

	a.telegram = NewBot(c)
	if a.telegram == nil {
		return errors.New("cant create telegram")
	}

	return nil
}

func (a *App) Run() error {
	return a.telegram.Run()
}

func (a *App) Resolve() error {
	panic("RESOLVER RUN")
}
