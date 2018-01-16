package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/tylerchr/backtalk"
)

func main() {

	var trainingPath string
	var modelPath string
	var threshold float64

	flag.StringVar(&trainingPath, "samples", "", "path to training data file")
	flag.StringVar(&modelPath, "o", "model.json", "path to write final model to")
	flag.Float64Var(&threshold, "threshold", 0.7, "acceptance threshold")
	flag.Parse()

	if trainingPath == "" {
		panic("no training data provided")
	}

	samples, err := ReadSamples(trainingPath)
	if err != nil {
		panic(err)
	}

	// extract intent names
	var intents []string
	for k := range samples {
		intents = append(intents, k)
	}

	t0 := time.Now()

	// create new model
	classifier := backtalk.NewNaiveBayesClassifier(intents, threshold)

	// train the model
	for intent, strings := range samples {
		for _, data := range strings {
			if err := classifier.Train(intent, data); err != nil {
				panic(fmt.Errorf("unexpected training error: %s", err))
			}
		}
	}

	if err := classifier.DoneTraining(); err != nil {
		panic(fmt.Errorf("unexpected training error: %s", err))
	}

	f, err := os.Create(modelPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := classifier.WriteTo(f); err != nil {
		panic(err)
	}

	fmt.Printf("%s (trained on %d intents in %s)\n", modelPath, len(intents), time.Since(t0))

}

func ReadSamples(file string) (map[string][]string, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var model struct {
		Intents map[string][]string
	}

	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return model.Intents, nil

}
