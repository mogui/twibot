package bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Config structure representing the json man config file
type Config struct {
	ConsumerKey     string   `json:"consumer_key"`
	ConsumerSecret  string   `json:"consumer_secret"`
	Token           string   `json:"token"`
	TokenSecret     string   `json:"token_secret"`
	BotName         string   `json:"screen_name"`
	AuthorizedUsers []string `json:"authorized_account"`
	OnMentions      []Action `json:"on_mentions"`
	OnDM            []Action `json:"on_dm"`
}

// FromJSON will parse the given json file and fill the Config structure
func (c *Config) FromJSON(path string) error {
	// parse hook file for hooks
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("ERROR loading file: ", err.Error())
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return err
	}

	err = c.validate()
	if err != nil {
		return err
	}

	// Load Bundled Commands
	c.OnDM = append(BundledCommands, c.OnDM...)
	return nil
}

// IsValid will check the integrity of the config struct
func (c *Config) validate() error {
	if !(c.ConsumerKey != "" && c.ConsumerSecret != "" && c.Token != "" && c.TokenSecret != "") {
		return errors.New("You haven't sepcified all the necessary Twitter API keys")
	}
	actions := append(c.OnDM, c.OnMentions...)
	for _, x := range actions {
		if err := x.Validate(); err != nil {
			return err
		}
	}

	return nil
}
