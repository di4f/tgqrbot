package main


import (
	"os"
	"log"
	"github.com/omnipunk/tg"
	qr "github.com/omnipunk/qr"
	"net/http"
	//"io"
	"fmt"
	"image/jpeg"
)

type Context = tg.Context
type RenderFunc = tg.RenderFunc
type Func = tg.Func

var beh = tg.NewBehaviour().WithRootNode(tg.NewRootNode(
	Func(func(c *Context){
		c.Sendf("The bot started! Now send me photo of QR code to recognize it.")
		for u := range c.Input() {
			if u.Message != nil && u.Message.Photo != nil {
				photo := u.Message.Photo[len(u.Message.Photo) - 1]
				file, err := c.Bot.Api.GetFile(tg.FileConfig{FileID: photo.FileID})
				if err != nil {
					c.Sendf("err: %q", err.Error())
					continue
				}

				r, err := http.Get(fmt.Sprintf(
					"https://api.telegram.org/file/bot%s/%s",
					c.Bot.Api.Token,
					file.FilePath,
				))
				if err != nil || r.StatusCode != 200 {
					if err != nil {
						c.Sendf("err: %q", err.Error())
					}
					continue
				}
				defer r.Body.Close()

				img, err := jpeg.Decode(r.Body)
				if err != nil {
					c.Sendf("err: %q", err.Error())
					continue
				}
				codes, err := qr.Recognize(img)
				if err != nil {
					c.Sendf("err: %q", err.Error())
					continue
				}
				for _, code := range codes {
					c.Sendf2("`%s`", tg.Escape2(string(code.Payload)))
				}
			}
		}
	}),
)).WithRoot(
	tg.NewCommandCompo().WithUsage(tg.Func(func(c *Context) {
		c.Sendf("There is no such command %q", c.Message.Command())
	})).WithPreStart(tg.Func(func(c *Context) {
		c.Sendf("Please, use /start ")
	})).WithCommands(
		tg.NewCommand(
			"start",
			"start the bot",
		).Go("/"),
	),
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	bot, err := tg.NewBot(token)
	if err != nil {
		log.Panic(err)
	}
	bot = bot.
		WithBehaviour(beh).
		Debug(true)

	log.Printf("Authorized on account %s", bot.Api.Self.UserName)
	err = bot.Run()
	if err != nil {
		panic(err)
	}
}
