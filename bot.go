package backtalk

import (
	"context"
	"fmt"

	"github.com/nlopes/slack"
)

type Bot struct {
	api *slack.Client
	rtm *slack.RTM
}

func New(api *slack.Client) (*Bot, error) {
	b := &Bot{
		api: api,
		rtm: api.NewRTM(),
	}
	return b, nil
}

// RTM provides access to the underlying real-time messaging client
// to Slack.
func (b *Bot) RTM() *slack.RTM {
	return b.rtm
}

// Start starts listening.
func (b *Bot) Start(ctx context.Context, handler Replyer) error {

	go b.rtm.ManageConnection()

Loop:
	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		case msg := <-b.rtm.IncomingEvents:

			// TODO(tylerchr): Make Handler receive ALL events; use a special filter
			//                 to get just the MessageEvents if that's what you want.

			switch evt := msg.Data.(type) {

			case *slack.ConnectedEvent:
				fmt.Printf("Connection counter: %d\n", evt.ConnectionCount)

			case *slack.MessageEvent:

				info := b.rtm.GetInfo()

				// ignore messages from ourself
				if evt.User == info.User.ID {
					break
				}

				if err := handler.Reply(b.rtm, evt); err != nil {
					return err
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", evt.Error())

			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
				break Loop

				// case *slack.UserTypingEvent:
				// 	info := rtm.GetInfo()
				// 	user, channel := info.GetUserByID(evt.User), info.GetChannelByID(evt.Channel)
				// 	fmt.Printf("%s is typing in %#v\n", user.Name, channel)

				// default:
				// 	fmt.Printf("Unknown event: %T\n", evt)

			}

		}
	}

	return nil
}
