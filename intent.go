package backtalk

import (
	"log"

	"github.com/nlopes/slack"
)

// UnableToClassifyIntent is a meta-intent for which a Replyer can be registered
// to react to unclassifiable inputs.
const UnableToClassifyIntent = "Backtalk_MetaIntent_UnableToClassifyIntent"

type (
	IntentMux struct {
		classifier Classifier
		intents    map[string]Replyer
	}

	Classifier interface {
		Train(intent, input string) error
		Classify(input string) (string, error)
	}
)

func NewIntentMux(c Classifier) (*IntentMux, error) {
	return &IntentMux{
		classifier: c,
		intents:    make(map[string]Replyer),
	}, nil
}

// Intent registers a new intent.
func (im *IntentMux) Intent(intent string, reply Replyer) {

	if im.intents == nil {
		im.intents = make(map[string]Replyer)
	}

	im.intents[intent] = reply

}

func (im *IntentMux) Reply(rtm *slack.RTM, evt *slack.MessageEvent) error {

	// run the classifier on the text to decide on the intent

	intent, err := im.classifier.Classify(evt.Text)
	if err != nil {
		log.Printf("unable to classify text: '%s'", evt.Text)
		intent = UnableToClassifyIntent
	}

	// delegate to registered replyer
	if replyer, ok := im.intents[intent]; ok && replyer != nil {
		return replyer.Reply(rtm, evt)
	}

	return nil

}
