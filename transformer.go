package main

import (
	"github.com/pborman/uuid"
)

func transformTopic(t term) topic {
	return topic{
		UUID:          uuid.NewMD5(uuid.UUID{}, []byte(t.ID)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.ID,
		Type:          "Topic",
	}
}
