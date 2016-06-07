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
			RawID:         "Nstein_GL_AFTM_GL_164835"},
			topic{
				UUID:      "0299feb1-7cb5-3ba2-865d-a2df7d670691",
				PrefLabel: "Africa Inc",
				AlternativeIdentifiers: alternativeIdentifiers{
					TME:   []string{"TnN0ZWluX0dMX0FGVE1fR0xfMTY0ODM1-VG9waWNz"},
					Uuids: []string{"0299feb1-7cb5-3ba2-865d-a2df7d670691"},
				},
				Type: "Topic"}},
	}

	for _, test := range tests {
		expectedTopic := transformTopic(test.term, "Topics")
		assert.Equal(test.topic, expectedTopic, fmt.Sprintf("%s: Expected topic incorrect", test.name))
	}

}
