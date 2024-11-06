package model

type AppState struct {
	ReplicationId     string
	ReplicationOffset int
}

var State AppState
