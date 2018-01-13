package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"

	"github.com/tylerchr/backtalk"
)

var SlackToken = os.Getenv("SLACK_TOKEN")

type (
	Accomplishment struct {
		User      string
		Timestamp time.Time
		Text      string
	}

	Storage interface {
		Save(a Accomplishment) error
		Range(user string, from, to time.Time, cb func(i int, a Accomplishment) bool) error
	}

	memoryStorage []Accomplishment
)

func (ms *memoryStorage) Save(a Accomplishment) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}
	*ms = append(*ms, a)
	return nil
}

func (ms memoryStorage) Range(u string, from, to time.Time, cb func(i int, a Accomplishment) bool) error {
	var i int
RangeLoop:
	for _, a := range ms {

		switch {
		case a.User != u:
			// not the right user
			continue
		case !from.IsZero() && a.Timestamp.Before(from):
			// happened too early
			continue
		case !to.IsZero() && a.Timestamp.After(to):
			// happened too late
			continue
		default:
			if stop := cb(i, a); stop {
				break RangeLoop
			}
			i++
		}
	}

	return nil
}

func main() {

	var s Storage = &memoryStorage{}

	api := slack.New(SlackToken)

	// start the bot
	ijust, _ := backtalk.New(api)
	ijust.Start(context.TODO(), backtalk.ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {

		// look up the user
		info := rtm.GetInfo()
		user := info.GetUserByID(evt.User)
		// fmt.Printf("From: %#v\n", user)

		prefix := fmt.Sprintf("<@%s> ", info.User.ID)
		text := strings.TrimSpace(strings.TrimPrefix(evt.Text, prefix))

		switch {

		// Special command: list accomplishments
		case strings.HasPrefix(strings.ToLower(text), "what"):
			fmt.Printf("Looking up accomplishments of %s\n", user.Name)

			var all []Accomplishment

			s.Range(user.Name, time.Now().Add(-24*time.Hour), time.Now(), func(i int, a Accomplishment) bool {
				fmt.Printf("%d\t%s\t%s\t%s\n", i, a.User, a.Timestamp, a.Text)
				all = append(all, a)
				return false
			})

			if len(all) == 0 {
				rtm.SendMessage(rtm.NewOutgoingMessage(
					fmt.Sprintf("<@%s> Looks like you've gotten nothing done!", evt.User),
					evt.Channel,
				))
			} else {
				var list string
				for _, c := range all {
					list = fmt.Sprintf("%s\n- %s (%s ago)", list, c.Text, time.Now().Sub(c.Timestamp).Round(1*time.Second))
				}

				rtm.SendMessage(rtm.NewOutgoingMessage(
					fmt.Sprintf("<@%s> Here's the rundown:\n%s", evt.User, strings.TrimSpace(list)),
					evt.Channel,
				))
			}

		// New accomplishment
		default:

			a := Accomplishment{
				User:      user.Name,
				Timestamp: backtalk.ParseTime(evt.Timestamp),
				Text:      text,
			}

			fmt.Printf("New accomplishment from %s: %s\n", user.Name, a)

			s.Save(a)

			// add a reaction
			if err := api.AddReaction("thumbsup", slack.NewRefToMessage(evt.Channel, evt.Timestamp)); err != nil {
				fmt.Printf("error adding reaction: %s\n", err)
			}

		}

		return nil
	}))

}
