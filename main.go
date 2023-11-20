package main


import (
	"os"
	"log"
	"github.com/omnipunk/tg"
)

type Context = tg.Context
type RenderFunc = tg.RenderFunc
type Func = tg.Func

var beh = tg.NewBehaviour().WithRootNode(tg.NewRootNode(
	Func(func(c *Context){
		c.Sendf("The bot started! Now send me photo of QR code to recognize it.")
		for _ = range c.Input() {
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
