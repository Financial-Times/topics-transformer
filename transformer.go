package main

import (
	"github.com/pborman/uuid"
)

func transformTopic(t term) topic {
	uuid := uuid.NewMD5(uuid.UUID{}, []byte(t.ID)).String()

	return topic{
		UUID:                   uuid,
		PrefLabel:              t.CanonicalName,
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{t.ID}, Uuids: []string{uuid}},
		Type: "Topic",
	}
}
