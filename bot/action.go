package bot

import (
	"fmt"
	"github.com/dougnukem/go-twitter/twitter"
	"os/exec"
	"regexp"
	"strings"
)

// Action is a struct that contains details on the comman to launch
type Action struct {
	Name           string `json:"name"`
	Match          string `json:"match"`
	Script         string `json:"script"`
	Reply          bool   `json:"reply,omitempty"`
	Case           bool   `json:"case,omitempty"`
	Bundled        bool
	BundledCommand func(text string, user *twitter.User, app *Twibot)
}

// WillMatch will detect if a given tweet match with its action and if it has
// groups it will substitute them in the script command
func (a *Action) WillMatch(tweet string) bool {
	var regex *regexp.Regexp

	if a.Case {
		regex = regexp.MustCompile(a.Match)
	} else {
		regex = regexp.MustCompile("(?i)" + a.Match)
	}

	match := regex.FindStringSubmatch(tweet)

	if len(match) > 0 {
		if len(match) > 1 {
			// Got groups gonna substitute the placeholder in the command
			for i := 1; i < len(match); i++ {
				a.Script = strings.Replace(a.Script, fmt.Sprintf("{%d}", i), match[i], -1)
			}
		}
		return true
	}
	return false
}

// Validate will check if the config is the minimum acceptable conf
func (a *Action) Validate() error {
	if !(a.Name != "" && a.Match != "" && (a.Script != "" || a.Bundled)) {
		return fmt.Errorf("command: %s failed validating -> name, match, script are all mandatory fields", a.Name)
	}
	return nil
}

// Exec will execute his inner command returning the output or the occurred error
func (a *Action) Exec() (string, error) {
	var (
		cmdOut []byte
		err    error
	)
	parts := strings.Fields(a.Script)
	head := parts[0]
	parts = parts[1:len(parts)]
	if cmdOut, err = exec.Command(head, parts...).Output(); err != nil {
		return "", err
	}
	return string(cmdOut), nil
}
