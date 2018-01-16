package main

import (
	"io/ioutil"
	"testing"

	"github.com/tylerchr/backtalk"
)

func TestModel(t *testing.T) {

	cases := []struct {
		Intent string
		Data   string
	}{
		{Intent: "Morning", Data: "good morning"},
		{Intent: "Morning", Data: "the sun is coming up"},
		{Intent: "Morning", Data: "what is for breakfast"},
		{Intent: "Afternoon", Data: "good afternoon"},
		{Intent: "Afternoon", Data: "good day"},
		{Intent: "Afternoon", Data: "it's hot outside"},
		{Intent: "Afternoon", Data: "time for lunch"},
		{Intent: "Evening", Data: "what is for dinner"},
		{Intent: "Evening", Data: "the sun is going down"},
		{Intent: "Night", Data: "it's cold outside"},
		{Intent: "Night", Data: "the moon is out"},
		{Intent: "Night", Data: "the look at the stars"},
	}

	model, err := ioutil.ReadFile("model.json")
	if err != nil {
		t.Fatalf("unable to load model file: %s", err)
	}
	classifier, err := backtalk.NewNaiveBayesClassifierFromModel(model)
	if err != nil {
		t.Fatalf("unable to restore classifier: %s", err)
	}

	// train the model
	for i, c := range cases {
		if intent, err := classifier.Classify(c.Data); err != nil {
			t.Errorf("[case %d] unexpected training error: %s", i, err)
		} else if intent != c.Intent {
			t.Errorf("[case %d] unexpected classification: expected %s but got %s [\"%s\"", i, c.Intent, intent, c.Data)
		}
	}
}
