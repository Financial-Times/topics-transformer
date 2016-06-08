package main

type topic struct {
	UUID                   string                 `json:"uuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	PrefLabel              string                 `json:"prefLabel"`
	Type                   string                 `json:"type"`
}

type alternativeIdentifiers struct {
	TME               []string `json:"TME,omitempty"`
	FactsetIdentifier string   `json:"factsetIdentifier,omitempty"`
	LeiCode           string   `json:"leiCode,omitempty"`
	Uuids             []string `json:"uuids,omitempty"`
}

type topicLink struct {
	APIURL string `json:"apiUrl"`
}
