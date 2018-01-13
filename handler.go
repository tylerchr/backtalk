package backtalk

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

type (
	Replyer interface {
		// Handle responds to an RTM event.
		//
		// Returning an error from a handler will kill the bot.
		Reply(rtm *slack.RTM, evt *slack.MessageEvent) error
	}

	ReplyFunc func(rtm *slack.RTM, evt *slack.MessageEvent) error
)

// Handle implements Replyer by invoking the ReplyFunc.
func (hf ReplyFunc) Reply(rtm *slack.RTM, evt *slack.MessageEvent) error {
	return hf(rtm, evt)
}

// DirectFilterReplyer is a filtering Replyer that excludes all messages that
// aren't explicitly directed at the bot (i.e., @bot).
func DirectFilterReplyer(r Replyer) Replyer {
	return ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {

		// look up the sending user
		info := rtm.GetInfo()

		var isDirectMessage bool
		for _, im := range info.IMs {
			if evt.Channel == im.ID {
				isDirectMessage = true
				break
			}
		}

		// ignore messages not directed at us
		prefix := fmt.Sprintf("<@%s> ", info.User.ID)
		if !isDirectMessage && !strings.HasPrefix(evt.Text, prefix) {
			return nil
		}

		return r.Reply(rtm, evt)

	})
}
