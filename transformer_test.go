package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name  string
		term  term
		topic topic
	}{
		{"Transform term to topic", term{
			CanonicalName: "Africa Inc",
			ID:            "YTcyNWI5YzItOTUwMy00ZWRkLWI0M2YtYzBjZjU5MWNjNTJi-VG9waWNz"},
			topic{
				UUID:      "c6c9c5f0-b5f6-3392-be0c-f82b6115c40b",
				PrefLabel: "Africa Inc",
				AlternativeIdentifiers: alternativeIdentifiers{
					TME:   []string{"YTcyNWI5YzItOTUwMy00ZWRkLWI0M2YtYzBjZjU5MWNjNTJi-VG9waWNz"},
					Uuids: []string{"c6c9c5f0-b5f6-3392-be0c-f82b6115c40b"},
				},
				Type: "Topic"}},
	}

	for _, test := range tests {
		expectedTopic := transformTopic(test.term)
		assert.Equal(test.topic, expectedTopic, fmt.Sprintf("%s: Expected topic incorrect", test.name))
	}

}
