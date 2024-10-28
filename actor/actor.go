package actor

// Actor 角色抽象
type Actor interface {
	Start()
	Stop()
	GetRoleID() uint64
	Logout() bool
	GetMsgChan() chan Msg
}

type Msg interface {
	Proc()
}

type Actors struct {
	Actors map[uint64]Actor
}

func (actors *Actors) AddActor(id uint64, actor Actor) {
	if actors.Actors == nil {
		actors.Actors = make(map[uint64]Actor)
	}
	actors.Actors[id] = actor
}

func (actors *Actors) RemoveActor(id uint64) {
	delete(actors.Actors, id)
}

func (actors *Actors) ForeachActor(f func(actor Actor)) {
	for _, actor := range actors.Actors {
		f(actor)
	}
}

func CreateRoleActors() *Actors {
	return &Actors{Actors: make(map[uint64]Actor)}
}
