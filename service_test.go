package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTopics(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		baseURL  string
		tax      taxonomy
		topics []topicLink
		found    bool
		err      error
	}{
		{"Success", "localhost:8080/transformers/topics/",
			taxonomy{Terms: []term{term{CanonicalName: "Z_Archive", ID: "NWQxOTE4ZGQtOGY1OS00MWY3LTk0ZWEtOWYyOGNmMDg4ZGJk-VG9waWNz", Children: children{[]term{term{CanonicalName: "Africa Inc", ID: "YTcyNWI5YzItOTUwMy00ZWRkLWI0M2YtYzBjZjU5MWNjNTJi-VG9waWNz"}}}}}},
			[]topicLink{topicLink{APIURL: "localhost:8080/transformers/topics/e5360bf3-2068-3660-9a50-9c4d93b4bae1"},
				topicLink{APIURL: "localhost:8080/transformers/topics/c6c9c5f0-b5f6-3392-be0c-f82b6115c40b"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/topics/", taxonomy{}, []topicLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newTopicService(&repo, test.baseURL)
		expectedTopics, found := service.getTopics()
		assert.Equal(test.topics, expectedTopics, fmt.Sprintf("%s: Expected topics link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetTopicByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		tax     taxonomy
		uuid    string
		topic topic
		found   bool
		err     error
	}{
		{"Success", taxonomy{Terms: []term{term{CanonicalName: "Z_Archive", ID: "NWQxOTE4ZGQtOGY1OS00MWY3LTk0ZWEtOWYyOGNmMDg4ZGJk-VG9waWNz", Children: children{[]term{term{CanonicalName: "Africa Inc", ID: "YTcyNWI5YzItOTUwMy00ZWRkLWI0M2YtYzBjZjU5MWNjNTJi-VG9waWNz"}}}}}},
			"e5360bf3-2068-3660-9a50-9c4d93b4bae1", topic{UUID: "e5360bf3-2068-3660-9a50-9c4d93b4bae1", CanonicalName: "Z_Archive", TmeIdentifier: "NWQxOTE4ZGQtOGY1OS00MWY3LTk0ZWEtOWYyOGNmMDg4ZGJk-VG9waWNz", Type: "Topic"}, true, nil},
		{"Not found", taxonomy{Terms: []term{term{CanonicalName: "Z_Archive", ID: "NWQxOTE4ZGQtOGY1OS00MWY3LTk0ZWEtOWYyOGNmMDg4ZGJk-VG9waWNz", Children: children{[]term{term{CanonicalName: "Africa Inc", ID: "YTcyNWI5YzItOTUwMy00ZWRkLWI0M2YtYzBjZjU5MWNjNTJi-VG9waWNz"}}}}}},
			"some uuid", topic{}, false, nil},
		{"Error on init", taxonomy{}, "some uuid", topic{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newTopicService(&repo, "")
		expectedTopic, found := service.getTopicByUUID(test.uuid)
		assert.Equal(test.topic, expectedTopic, fmt.Sprintf("%s: Expected topic incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	tax taxonomy
	err error
}

func (d *dummyRepo) getTopicsTaxonomy() (taxonomy, error) {
	return d.tax, d.err
}
