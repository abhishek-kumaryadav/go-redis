package model

import "github.com/google/uuid"

type AppState struct {
	ReplicationId     uuid.UUID
	ReplicationOffset int
}

var State AppState
