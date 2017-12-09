package bot

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

type (
	Handler interface {
		// Handle responds to an RTM event.
		//
		// Returning an error from a handler will kill the bot.
		Handle(rtm *slack.RTM, evt *slack.MessageEvent) error
	}

	HandlerFunc func(rtm *slack.RTM, evt *slack.MessageEvent) error
)

// Handle implements Handler by invoking the HandlerFunc.
func (hf HandlerFunc) Handle(rtm *slack.RTM, evt *slack.MessageEvent) error {
	return hf(rtm, evt)
}

// DirectFilterHandler is a filtering Handler that excludes all messages that
// aren't explicitly directed at the bot (i.e., @bot).
func DirectFilterHandler(h Handler) Handler {
	return HandlerFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {

		info := rtm.GetInfo()
		prefix := fmt.Sprintf("<@%s> ", info.User.ID)

		if !strings.HasPrefix(evt.Text, prefix) {
			return nil
		}

		return h.Handle(rtm, evt)

	})
}
