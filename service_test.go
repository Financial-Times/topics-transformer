package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTopics(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		baseURL string
		terms   []term
		topics  []topicLink
		found   bool
		err     error
	}{
		{"Success", "localhost:8080/transformers/topics/",
			[]term{term{CanonicalName: "Z_Archive", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Africa Inc", RawID: "Nstein_GL_AFTM_GL_164835"}},
			[]topicLink{topicLink{APIURL: "localhost:8080/transformers/topics/81918290-8f91-3722-9ef3-aa9f31cf9e43"},
				topicLink{APIURL: "localhost:8080/transformers/topics/0299feb1-7cb5-3ba2-865d-a2df7d670691"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/topics/", []term{}, []topicLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newTopicService(&repo, test.baseURL, "Topics", 10000)
		expectedTopics, found := service.getTopics()
		assert.Equal(test.topics, expectedTopics, fmt.Sprintf("%s: Expected topics link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetTopicByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name  string
		terms []term
		uuid  string
		topic topic
		found bool
		err   error
	}{
		{"Success", []term{term{CanonicalName: "Z_Archive", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Africa Inc", RawID: "Nstein_GL_AFTM_GL_164835"}},
			"81918290-8f91-3722-9ef3-aa9f31cf9e43", getDummyTopic("81918290-8f91-3722-9ef3-aa9f31cf9e43", "Z_Archive", "ODQ1ZGM3ZDctYWU4OS00ZmVkLWE4MTktOWVkY2JiM2ZlNTA3-VG9waWNz"), true, nil},
		{"Not found", []term{term{CanonicalName: "Z_Archive", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Africa Inc", RawID: "Nstein_GL_AFTM_GL_164835"}},
			"some uuid", topic{}, false, nil},
		{"Error on init", []term{}, "some uuid", topic{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newTopicService(&repo, "", "Topics", 10000)
		expectedTopic, found := service.getTopicByUUID(test.uuid)
		assert.Equal(test.topic, expectedTopic, fmt.Sprintf("%s: Expected topic incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	terms []term
	err   error
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	if startRecord > 0 {
		return nil, d.err
	}
	var interfaces = make([]interface{}, len(d.terms))
	for i, data := range d.terms {
		interfaces[i] = data
	}
	return interfaces, d.err
}
func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return d.terms[0], d.err
}

func getDummyTopic(uuid string, prefLabel string, tmeId string) topic {
	return topic{
		UUID:                   uuid,
		PrefLabel:              prefLabel,
		PrimaryType:            primaryType,
		TypeHierarchy:          topicTypes,
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeId}, Uuids: []string{uuid}}}
}
