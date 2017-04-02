package main

//import "encoding/json"

type Command struct {
	Player  *Player `json:"player"`
	Command string  `json:"command"`
}

type ObjSubj struct {
	Key      string `json:"key"`
	Exist    bool   `json:"exist"`
	Lock     bool   `json:"lock"`
	NameRoom string `json:"nameroom"`
}
type Player struct {
	InRoom  string      `json:"inroom"`
	RefBack *Back       `json:"refback"`
	Name    string      `json:"name"`
	Msg     chan string `json:"msg"`
}
type Back struct {
	Things map[string]bool `json:"things"`
	Full   bool            `json:"full"`
}
type Game struct {
	Msgin   chan *Command     `json:"msgin"`
	Priory  []string          `json:"priory"`
	Rooms   map[string]*Room  `json:"rooms"`
	Aliases map[string]string `json:"aliases"`
}

type Room struct {
	Players  map[string]*Player `json:"players"`
	Name     string             `json:"name"`
	Msg      map[string]string  `json:"msg"`
	Things   map[string]bool    `json:"things"`
	Subjects map[string]ObjSubj `json:"subject"`
	Act      string             `json:"act"`
	LinkRoom map[string]string  `json:"linkroom"`
	// Decription func(r *Room, pl *Player) string `json:"dec"`
}

func (r *Room) Label() string {
	return r.Name
}

func (r *Room) Type() string {
	return "комнаты"
}
