package main

import (
	"github.com/codegangsta/cli"
	"github.com/mogui/twibot/bot"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "twibot"
	app.Usage = "A twitter bot to execute commands by DM or @mentions"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "conf.json",
			Usage: "config file for the bot",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Spit verbose logging",
		},
	}
	app.Action = func(c *cli.Context) {

		verbose := c.Bool("verbose")
		path := c.String("config")

		// Run twibot
		twibot := new(bot.Twibot)
		twibot.SetupLog(verbose)
		twibot.Run(path)
	}
	app.Run(os.Args)
}
