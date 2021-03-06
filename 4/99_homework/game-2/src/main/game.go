package main

import "fmt"

var G Game

func AddPlayer(p *Player) {
	G.Rooms["кухня"].Players[p.Name] = p
	p.InRoom = "кухня"
}

func GetRoom(name string) *Room {
	r, ok := G.Rooms[name]
	if ok {
		return r
	}
	alias, ok1 := G.Aliases[name]
	if ok1 {
		return GetRoom(alias)
	} else {
		panic(fmt.Sprintf("Не могу найти комнату по ключу %s", name))
	}
}

func NewPlayer(name string) (player *Player) {
	p := &Player{
		Name: name,
		Msg:  make(chan string),
	}
	return p
}
func InitGame() {
	G = Game{
		Rooms:   make(map[string]*Room),
		Aliases: make(map[string]string),
		Msgin:   make(chan *Command),
	}

	G.Rooms = InitRoom()
	for name, room := range G.Rooms {
		room.Players = make(map[string]*Player)
		room.Name = name
	}

	G.Priory = []string{
		"кухня",
		"комната",
		"улица",
	}

	G.Aliases = map[string]string{
		"улица": "домой",
	}

	go Run()
}

func Run() {
	for cmd := range G.Msgin {
		cmd.Player.don(cmd.Command)
	}
}

func InitRoom() map[string]*Room {
	return map[string]*Room{
		"кухня": {
			Act: "идти в универ. ",
			Msg: map[string]string{
				"enter":      "кухня, ничего интересного.",
				"lookaround": "ты находишься на кухне, на столе ",
				"backact":    "надо собрать ",
				"end":        "можно пройти - коридор",
			}, Things: map[string]bool{
				"чай": true,
			}, Subjects: map[string]ObjSubj{
				"ничего": {
					Exist: true, Lock: true},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
		"коридор": {
			Msg: map[string]string{
				"enter": "ничего интересного.",
				"lookaround": "	ты находишься в коридоре, тут страшно ",
				"end": "можно пройти - кухня, комната, улица",
			}, Subjects: map[string]ObjSubj{
				"шкаф": {
					Exist: true, Lock: true, Key: "кот"},
				"дверь": {
					Exist: true, Lock: true, Key: "ключи", NameRoom: "улица"},
			}, LinkRoom: map[string]string{
				"кухня":   "exist",
				"комната": "exist",
				"улица":   "lock"},
		},
		"комната": {
			Msg: map[string]string{
				"enter":      "ты в своей комнате.",
				"lookaround": "на столе: ",
				"backact":    "на стуле - ",
				"end":        "можно пройти - коридор",
			}, Things: map[string]bool{
				"ключи":     true,
				"конспекты": true,
				"рюкзак":    true,
			}, Subjects: map[string]ObjSubj{
				"тумбочка": {
					Exist: true, Lock: true, Key: "хлопок в ладоши"},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
		"улица": {
			Msg: map[string]string{
				"enter":      "на улице весна.",
				"lookaround": "ты находишься на улице, как же хорошо",
				"locked":     "дверь закрыта",
			}, Subjects: map[string]ObjSubj{
				"шкаф": {
					Exist: true, Lock: true, Key: "кот"},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
	}
}
