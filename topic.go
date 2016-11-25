package main

type topic struct {
	UUID                   string                 `json:"uuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	PrefLabel              string                 `json:"prefLabel"`
	PrimaryType            string                 `json:"type"`
	TypeHierarchy          []string               `json:"types"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	Uuids []string `json:"uuids,omitempty"`
}

type topicLink struct {
	APIURL string `json:"apiUrl"`
}

var primaryType = "Topic"
var topicTypes = []string{"Thing", "Concept", "Topic"}
