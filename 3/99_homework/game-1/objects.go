package main

import "sync"

type Command struct {
	player  *Player
	command string
}

type ObjSubj struct {
	Key      string
	Exist    bool
	Lock     bool
	NameRoom string
}
type Player struct {
	InRoom  *Room
	RefBack *Back
	Name    string
	msg     *Chan
}
type Back struct {
	Things map[string]bool
	Full   bool
}
type Game struct {
	msgin   chan *Command
	wg      sync.WaitGroup
	Priory  []string
	Rooms   map[string]*Room
	Players map[string]*Player
	Aliases map[string]string
}

type Room struct {
	Name       string
	Msg        map[string]string
	Things     map[string]bool
	Subjects   map[string]ObjSubj
	Act        string
	LinkRoom   map[string]string
	Decription func(r *Room, pl *Player) string
}

type Chan struct {
	sync.Mutex
	msg chan string
}

func (r *Room) Label() string {
	return r.Name
}

func (r *Room) Type() string {
	return "комнаты"
}
