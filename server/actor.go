package server

import xactor "github.com/75912001/xlib/actor"

type Actor struct {
	*xactor.Actor[uint64] // actor, 该 id 为 server id
}

// NewActor creates a new Actor with the given id and behavior.
//
//	id: server id.
func NewActor(id uint64, behavior xactor.Behavior) *Actor {
	return &Actor{
		Actor: xactor.NewActor(id, nil, behavior),
	}
}
