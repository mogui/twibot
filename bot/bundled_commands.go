package bot

import (
	"bytes"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
)

// BundledCommands is an array of bundled actions of the bot
var BundledCommands = []Action{

	//  will just reply PONG
	Action{
		Name:    "PING",
		Match:   "ping",
		Bundled: true,
		BundledCommand: func(text string, user *twitter.User, app *Twibot) {
			// reply with pong
			app.client.DirectMessages.New(&twitter.DirectMessageNewParams{
				ScreenName: user.ScreenName,
				Text:       "PONG",
			})
		},
	},

	// will reply with a DM conatining all the available commands
	Action{
		Name:    "HELP",
		Match:   "help|\\?",
		Bundled: true,
		BundledCommand: func(text string, user *twitter.User, app *Twibot) {
			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprintf("[info]\nDirect Messages (%d):\n", len(app.conf.OnDM)))
			for i, action := range app.conf.OnDM {
				buffer.WriteString(fmt.Sprintf("%02d. [%s]: %s\n", i+1, action.Name, action.Match))
			}
			buffer.WriteString(fmt.Sprintf("\n Mentions (%d):\n", len(app.conf.OnMentions)))
			for i, action := range app.conf.OnMentions {
				buffer.WriteString(fmt.Sprintf("%02d. [%s]: %s\n", i+1, action.Name, action.Match))
			}

			app.client.DirectMessages.New(&twitter.DirectMessageNewParams{
				ScreenName: user.ScreenName,
				Text:       buffer.String(),
			})
		},
	},
}
