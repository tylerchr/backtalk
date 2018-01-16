# `btmodel` Model Generator

The `btmodel` model generator creates models suitable for use in the NaiveBayesClassifier. This allows for ahead-of-time model generation, which can be desirable for time and determinism reasons.

## Usage

Use of `btmodel` requires training samples that are pre-classified by intent. These will be used to calculate the Bayesian probabilities that define the model.

```
// training.json

{
  "Intents": {
    "Morning": [
      "good morning",
      "the sun is coming up",
      "what is for breakfast"
    ],
    "Afternoon": [
      "good afternoon",
      "good day",
      "it's hot outside",
      "time for lunch"
    ],
    "Evening": [
      "what is for dinner",
      "the sun is going down"
    ],
    "Night": [
      "it's cold outside",
      "the moon is out",
      "the look at the stars"
    ]
  }
}
```

```bash
$ go install github.com/tylerchr/backtalk/cmd/btmodel
$ btmodel -samples training.json
model.json (trained on 4 intents in 774.872µs)
```

At this point, the resulting `model.json` file can be loaded and its contents given to `backtalk.NewNaiveBayesClassifierFromModel` instead of training the model at runtime directly.

## Tests

Since machine learning based models can be rather unpredictable, it's recommended that you write a test to ensure that your critical inputs get classified property by the model. This strategy may help prevent broken or bad models from being deployed with your application.

```go
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
```