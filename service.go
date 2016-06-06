package main

import (
	"net/http"
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type topicService interface {
	getTopics() ([]topicLink, bool)
	getTopicByUUID(uuid string) (topic, bool)
	checkConnectivity() error
}

type topicServiceImpl struct {
	repository tmereader.Repository
	baseURL    string
	topicsMap  map[string]topic
	topicLinks []topicLink
	taxonomyName  string
	maxTmeRecords int
}

func newTopicService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (topicService, error) {
	s := &topicServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.init()
	if err != nil {
		return &topicServiceImpl{}, err
	}
	return s, nil
}

func (s *topicServiceImpl) init() error {
	s.topicsMap = make(map[string]topic)
	responseCount := 0
	log.Printf("Fetching locations from TME\n")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Printf("Finished fetching locations from TME\n")
			break
		}
		s.initTopicsMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d topic links\n", len(s.topicLinks))

	return nil
}

func (s *topicServiceImpl) getTopics() ([]topicLink, bool) {
	if len(s.topicLinks) > 0 {
		return s.topicLinks, true
	}
	return s.topicLinks, false
}

func (s *topicServiceImpl) getTopicByUUID(uuid string) (topic, bool) {
	topic, found := s.topicsMap[uuid]
	return topic, found
}

func (s *topicServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back? Maybe a healthcheck or gtg endpoint?
	//_, err := s.repository.getTopicsTaxonomy()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (s *topicServiceImpl) initTopicsMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		top := transformTopic(t, s.taxonomyName)
		s.topicsMap[top.UUID] = top
		s.topicLinks = append(s.topicLinks, topicLink{APIURL: s.baseURL + top.UUID})
	}
}
