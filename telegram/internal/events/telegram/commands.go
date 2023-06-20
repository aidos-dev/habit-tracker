package telegram

import (
	"errors"
	"log"
	"strings"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage"
	"github.com/gin-gonic/gin"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	Habit    = "/habit"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from [%s]", text, username)

	switch text {
	case StartCmd:
		return p.sendHello(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case Habit:
		return p.sendRandom(chatID, username)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = errs.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = errs.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int, username string) error {
	if !p.userExists(username) {
		p.signUp(chatID, username)
	}

	log.Printf("user [%v] started bot\n", username)

	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) userExists(username string) bool {
	// p.adapter.FindTgUser(username)
	// p.ginEng.GET("tg_user_name", p.adapter.FindTgUser)

	log.Print("method userExists called")

	var userExists bool

	log.Printf("userExists value: %v", userExists)

	urlPath := "/tgUser"

	p.adapter.Router.GET(urlPath, func(c *gin.Context) {
		p.adapter.FindTgUser(c, username, &userExists)
		name := c.Query("tg_user_name")

		c.Set(models.CtxUsername, name)
	})

	return userExists
}

func (p *Processor) signUp(chatID int, username string) {
	// defer func() { err = errs.WrapIfErr("can't do command: can't sign up", err) }()

	log.Print("method signUp called")

	urlPath := "/auth/sign-up"

	p.adapter.Router.POST(urlPath, func(c *gin.Context) {
		p.adapter.SignUp(c, username)
	})
}
