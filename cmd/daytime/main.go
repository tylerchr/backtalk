package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/nlopes/slack"

	"github.com/tylerchr/backtalk"
)

const (
	IntentMorning   = "Morning"
	IntentAfternoon = "Afternoon"
	IntentEvening   = "Evening"
	IntentNight     = "Night"
)

var SlackToken = os.Getenv("SLACK_TOKEN")

func main() {

	api := slack.New(SlackToken)

	im, _ := backtalk.NewIntentMux(LoadClassifier("model.json"))

	im.Intent(IntentMorning, backtalk.ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {
		log.Printf("Looks like morning (%s)", evt.Text)
		return nil
	}))

	im.Intent(IntentAfternoon, backtalk.ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {
		log.Printf("Looks like afternoon (%s)", evt.Text)
		return nil
	}))

	im.Intent(IntentEvening, backtalk.ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {
		log.Printf("Looks like evening (%s)", evt.Text)
		return nil
	}))

	im.Intent(IntentNight, backtalk.ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {
		log.Printf("Looks like nighttime (%s)", evt.Text)
		return nil
	}))

	im.Intent(backtalk.UnableToClassifyIntent, backtalk.ReplyFunc(func(rtm *slack.RTM, evt *slack.MessageEvent) error {
		log.Printf("Don't know what to make of it (%s)", evt.Text)
		return nil
	}))

	// start the bot
	ijust, _ := backtalk.New(api)
	ijust.Start(context.TODO(), backtalk.DirectFilterReplyer(im))

}

func LoadClassifier(path string) backtalk.Classifier {
	model, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	c, err := backtalk.NewNaiveBayesClassifierFromModel(model)
	if err != nil {
		panic(err)
	}
	return c
}
