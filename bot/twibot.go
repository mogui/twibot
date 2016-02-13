package bot

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("twibot")

// Twibot is the main app struct
type Twibot struct {
	conf   *Config
	client *twitter.Client
}

// SetupLog will setup the logger and its verbosity
func (t *Twibot) SetupLog(verbose bool) {

	var format = logging.MustStringFormatter(
		`[%{time:15:04:05.000}] %{level:.5s}: %{message}`,
	)
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	formatt := logging.NewBackendFormatter(backend1, format)
	leveled := logging.AddModuleLevel(formatt)

	logging.SetBackend(leveled)
	if verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
}

// Run will run twibot
func (t *Twibot) Run(path string) {

	// Parse config file
	t.conf = new(Config)
	log.Debugf("Loading config from %s...\n", path)
	err := t.conf.FromJSON(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Lodaded conf, Got %d actions\n", len(t.conf.OnMentions)+len(t.conf.OnDM))

	// Twitter client
	config := oauth1.NewConfig(t.conf.ConsumerKey, t.conf.ConsumerSecret)
	token := oauth1.NewToken(t.conf.Token, t.conf.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	t.client = twitter.NewClient(httpClient)
	stream, err := t.client.Streams.User(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Demux handlers
	demux := twitter.NewSwitchDemux()
	if len(t.conf.OnMentions) > 0 {
		demux.Tweet = t.handleTweet
	}
	if len(t.conf.OnDM) > 0 {
		demux.DM = t.handleDM
	}

	log.Info("start listening...")
	for message := range stream.Messages {
		demux.Handle(message)
	}
}

func (t *Twibot) handleMessage(text string, user *twitter.User, actions []Action) {
	log.Debugf("received from: %s message: %s", user.ScreenName, text)

	// check authorization
	if !t.isAuthorized(user.ScreenName) {
		return
	}

	for _, action := range actions {
		log.Debug("trying action: %s regex: %s case: %t", action.Name, action.Match, action.Case)
		if action.WillMatch(text) {
			log.Debug("OK")

			if action.Bundled {
				// if matched action is BUndled we run the bundled code
				log.Info("Executing Bundled Command: %s", action.Name)
				action.BundledCommand(text, user, t)
			} else {
				// We run the script
				log.Infof("Executing: %s", action.Name)
				out, err := action.Exec()

				var replyMessage string

				if err != nil {
					log.Errorf("executing %s: %s", action.Name, err.Error())
					replyMessage = fmt.Sprintf("[%s] FAILED:\n%s", action.Name, err.Error())
				} else {
					log.Debugf("executed correctly Task %s", action.Name)
					replyMessage = fmt.Sprintf("[%s] OK", action.Name)
					log.Debugf("output:\n%s", out)
					// TODO: decide what to do with output
				}

				if action.Reply {
					log.Debug("sending reply")

					t.client.DirectMessages.New(&twitter.DirectMessageNewParams{ScreenName: user.ScreenName, Text: replyMessage})
				}
			}
			return
		} else {
			log.Debug("KO")
		}
	}
	log.Debugf("not matched: %s ", text)
}

func (t *Twibot) isAuthorized(user string) bool {
	if len(t.conf.AuthorizedUsers) == 0 {
		return true
	}
	for _, a := range t.conf.AuthorizedUsers {
		if a == user {
			return true
		}
	}
	log.Debugf("user: %s is not authorized", user)
	return false
}

func (t *Twibot) handleDM(dm *twitter.DirectMessage) {
	if dm.SenderScreenName == t.conf.BotName {
		// ignore DM sent by ourselves
		return
	}
	t.handleMessage(dm.Text, dm.Sender, t.conf.OnDM)
}

func (t *Twibot) handleTweet(tweet *twitter.Tweet) {
	if tweet.User.ScreenName == t.conf.BotName {
		// ignore Tweet sent by ourselves
		return
	}
	// scan for mentions
	for _, mention := range tweet.Entities.UserMentions {
		if mention.ScreenName == t.conf.BotName {
			t.handleMessage(tweet.Text, tweet.User, t.conf.OnMentions)
			break
		}
	}
}
