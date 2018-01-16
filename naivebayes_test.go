package backtalk

import (
	"testing"
)

func TestClassifier(t *testing.T) {

	// get intent labels
	var intents []string
	for k := range TrainingMaterial {
		intents = append(intents, k)
	}

	classifier := NewNaiveBayesClassifier(intents, 0.5)

	// train the model
	for intent, strings := range TrainingMaterial {
		for _, data := range strings {
			if err := classifier.Train(intent, data); err != nil {
				t.Fatalf("unexpected training error: %s", err)
			}
		}
	}

	if err := classifier.DoneTraining(); err != nil {
		t.Fatalf("unexpected training error: %s", err)
	}

	// re-classify the training data (yeah, I know... but it's a naive test)
	for actualIntent, strings := range TrainingMaterial {
		for _, data := range strings {
			if classifiedIntent, err := classifier.Classify(data); err != nil {
				t.Errorf("unexpected classification error: %s", err)
			} else if classifiedIntent != actualIntent {
				t.Errorf("unexpected classification: expected %s but got %s for '%s'", actualIntent, classifiedIntent, data)
			}
		}
	}

}

var TrainingMaterial = map[string][]string{
	"WhatServiceInDC": []string{
		"What {Service} is out?",
		"What {Service} is in {Datacenter}?",
		"What {Service} is out in {Datacenter}?",
	},
	"WhenServiceReleased": []string{
		"When was {Service} [last] [released|deployed] [to {Datacenter}?",
		"When was {Service} released?",
		"When was {Service} deployed?",
		"When was {Service} last released?",
		"When was {Service} last deployed?",
		"When was {Service} released to {Datacenter}?",
		"When was {Service} deployed to {Datacenter}?",
		"When was {Service} last released to {Datacenter}?",
		"When was {Service} last deployed to {Datacenter}?",
	},
	"WhatReleaseHistory": []string{
		"What is the release history for {Service} in {Datacenter}?",
		"What is the release history of {Service} to {Datacenter}?",
		"What is the release history of {Service}?",
	},
	"WhatsWeird": []string{
		"What is weird with {Service}?",
		"What is weird in {Datacenter}?",
		"What is weird with {Service} in {Datacenter}?",
		"What's wrong with {Service}?",
		"What's suspicious with {Service}?",
		"What's unusual with {Service}?",
	},
}
