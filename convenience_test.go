package main

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {

	cases := []struct {
		Input string
		Time  time.Time
	}{
		{
			Input: "1512802153.000011",
			Time:  time.Unix(1512802153, 11000),
		},
		{
			Input: "1512802153.00001100",
			Time:  time.Unix(1512802153, 11000),
		},
	}

	for i, c := range cases {
		if tt := ParseTime(c.Input); !tt.Equal(c.Time) {
			t.Errorf("[case %d] unexpected time for %s: expected %s but got %s\n", i, c.Input, c.Time, tt)
		}
	}

}
