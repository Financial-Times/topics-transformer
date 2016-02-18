package main

type topic struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier"`
	Type          string `json:"type"`
}

type topicLink struct {
	APIURL string `json:"apiUrl"`
}
