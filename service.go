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
	getTopicCount() int
	getTopicIds() []string
	reload() error
}

type topicServiceImpl struct {
	repository    tmereader.Repository
	baseURL       string
	topicsMap     map[string]topic
	topicLinks    []topicLink
	taxonomyName  string
	maxTmeRecords int
}

func newTopicService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (topicService, error) {
	s := &topicServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.reload()
	if err != nil {
		return &topicServiceImpl{}, err
	}
	return s, nil
}

func (s *topicServiceImpl) reload() error {
	s.topicsMap = make(map[string]topic)
	var links []topicLink
	s.topicLinks = links
	responseCount := 0
	log.Println("Fetching topics from TME")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Println("Finished fetching topics from TME")
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
	// TODO: Can we use a count from our responses while actually in use to trigger a healthcheck?
	//	_, err := s.repository.GetTmeTermsFromIndex(1)
	//	if err != nil {
	//		return err
	//	}
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

func (s *topicServiceImpl) getTopicCount() int {
	return len(s.topicLinks)
}

func (s *topicServiceImpl) getTopicIds() []string {
	i := 0
	keys := make([]string, len(s.topicsMap))

	for k := range s.topicsMap {
		keys[i] = k
		i++
	}
	return keys
}
