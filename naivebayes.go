package backtalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

var (
	// ErrUnknownIntent indicates that an intent was claimed in training
	// that wasn't provided to the classifier at construction time.
	ErrUnknownIntent = errors.New("unknown intent")

	// ErrIntentUnclear indicates that the model is unable to provide a
	// classification with high enough confidence to clear the minimum
	// threshold.
	ErrIntentUnclear = errors.New("unable to classify intent")
)

type NaiveBayesClassifier struct {
	model     *text.NaiveBayes
	threshold float64

	labels      []string
	labelLookup map[string]uint8

	trainStream chan base.TextDatapoint
	trainErrors chan error
}

func NewNaiveBayesClassifier(intents []string, threshold float64) *NaiveBayesClassifier {

	stream := make(chan base.TextDatapoint)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, uint8(len(intents)), base.OnlyWordsAndNumbers)
	model.Output = ioutil.Discard
	go model.OnlineLearn(errors)

	// create a reverse index for looking up intent numbers by name
	labelLookup := make(map[string]uint8)
	for i, intent := range intents {
		labelLookup[intent] = uint8(i)
	}

	return &NaiveBayesClassifier{
		model:       model,
		threshold:   threshold,
		labels:      intents,
		labelLookup: labelLookup,
		trainStream: stream,
		trainErrors: errors,
	}

}

// NewNaiveBayesClassifierFromModel restores a serialized NaiveBayes model.
func NewNaiveBayesClassifierFromModel(model []byte) (*NaiveBayesClassifier, error) {

	var deserialized struct {
		Type      string
		Labels    []string
		Threshold float64
		Model     json.RawMessage
	}

	// read the model file
	if err := json.Unmarshal(model, &deserialized); err != nil {
		return nil, err
	}

	// sanity-check the model type
	if deserialized.Type != "NaiveBayes" {
		return nil, errors.New("unknown model type")
	}

	// create a reverse index for looking up intent numbers by name
	labelLookup := make(map[string]uint8)
	for i, intent := range deserialized.Labels {
		labelLookup[intent] = uint8(i)
	}

	// restore the model itself
	naiveBayes := new(text.NaiveBayes)
	if err := naiveBayes.Restore([]byte(deserialized.Model)); err != nil {
		return nil, err
	}

	return &NaiveBayesClassifier{
		model:       naiveBayes,
		threshold:   deserialized.Threshold,
		labels:      deserialized.Labels,
		labelLookup: labelLookup,
		trainStream: nil,
		trainErrors: nil,
	}, nil

}

func (nbc *NaiveBayesClassifier) Train(intent, input string) error {

	intentID, ok := nbc.labelLookup[intent]

	// verify that the intent to train for is known to the classifier
	if !ok {
		return ErrUnknownIntent
	}

	nbc.trainStream <- base.TextDatapoint{
		X: input,
		Y: intentID,
	}
	return nil

}

func (nbc *NaiveBayesClassifier) DoneTraining() error {
	close(nbc.trainStream)

	var errors []error
	for err := range nbc.trainErrors {
		log.Printf("training error: %#v\n", err)
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("experienced %d training errors", len(errors))
	}

	return nil
}

func (nbc *NaiveBayesClassifier) Classify(input string) (string, error) {

	intentID, probability := nbc.model.Probability(input)

	log.Printf("Classified to %s [%.3f]\n", nbc.labels[intentID], probability)

	if probability < nbc.threshold {
		return "", ErrIntentUnclear
	}

	return nbc.labels[intentID], nil
}

func (nbc *NaiveBayesClassifier) WriteTo(w io.Writer) (n int64, err error) {

	data, err := json.Marshal(struct {
		Type      string
		Labels    []string
		Threshold float64
		Model     *text.NaiveBayes
	}{
		Type:      "NaiveBayes",
		Labels:    nbc.labels,
		Threshold: nbc.threshold,
		Model:     nbc.model,
	})

	if err != nil {
		return 0, err
	}

	nn, err := w.Write(data)
	return int64(nn), err
}
