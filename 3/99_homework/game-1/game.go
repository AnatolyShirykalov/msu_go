package main

import "fmt"

var G Game

func addPlayer(p *Player) {
	G.Players[p.Name] = p
	p.InRoom = GetRoom("кухня")
}
func GetRoom(name string) *Room {
	r, ok := G.Rooms[name]
	if ok {
		return r
	} else {
		alias, ok1 := G.Aliases[name]
		if ok1 {
			return GetRoom(alias)
		} else {
			panic(fmt.Sprintf("Не могу найти комнату по ключу %s", name))
		}
	}
}

func NewPlayer(name string) (player *Player) {
	G.Players = map[string]*Player{
		name: {
			Name: name,
			Msg: &Chan{
				ChanMsg: make(chan string),
			},
		},
	}
	return G.Players[name]
}
func initGame() {
	G = Game{
		Rooms:   make(map[string]*Room),
		Players: make(map[string]*Player),
		Aliases: make(map[string]string),
		msgin:   make(chan *Command)}

	G.Rooms = InitRoom()
	for name, room := range G.Rooms {
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
	go run()
}

func run() {
	for cmd := range G.msgin {
		cmd.player.don(cmd.command)
		G.wg.Done()
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
				"lookaround": "	Ты нахдишься в коридоре, тут страшно",
			}, Things: map[string]bool{
				"ничего": true,
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

