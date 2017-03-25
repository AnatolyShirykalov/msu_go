package game

import (
	"fmt"
)

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
			Name:   name,
			msg:    make(chan string),
			msgout: make(chan string, 1),
		},
	}
	return G.Players[name]
}
func initGame() {
	G = Game{
		Rooms: make(map[string]*Room),
		// Links:   make([]Link, 0, 20),
		Players: make(map[string]*Player),
		Aliases: make(map[string]string),
		msgin:   make(chan *Command),
	}

	G.Rooms = map[string]*Room{
		"кухня": {
			Name: "кухня",
			Act:  "идти в универ. ",
			Msg: map[string]*subjLock{
				"notlinked":  {lable: "нет пути кухня"},
				"enter":      {lable: "кухня, ничего интересного."},
				"lookaround": {lable: "ты находишься на кухне, на столе "},
				"backact":    {lable: "надо собрать "},
				"end":        {lable: "можно пройти - коридор"},
			}, Things: map[string]bool{
				"чай": true,
			}, Subjects: map[string]ObjSubj{
				"ничего": {Exist: true, Lock: true},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
		"коридор": {
			Name: "коридор",
			Msg: map[string]*subjLock{
				"notlinked": {lable: "нет пути коридор"},
				"enter":     {lable: "ничего интересного."},
				"lookaround": {lable: "	Ты нахдишься в коридоре, тут страшно"},
			}, Things: map[string]bool{
				"ничего": true,
			}, Subjects: map[string]ObjSubj{
				"шкаф":  {Exist: true, Lock: true, Key: "кот"},
				"дверь": {Exist: true, Lock: true, Key: "ключи", NameRoom: "улица"},
			}, LinkRoom: map[string]string{
				"кухня":   "exist",
				"комната": "exist",
				"улица":   "lock"},
		},
		"комната": {
			Name: "комната",
			Msg: map[string]*subjLock{
				"notlinked":  {lable: "нет пути в комната"},
				"enter":      {lable: "ты в своей комнате."},
				"lookaround": {lable: "на столе: "},
				"backact":    {lable: "на стуле - "},
				"end":        {lable: "можно пройти - коридор"},
			}, Things: map[string]bool{
				"ключи":     true,
				"конспекты": true,
				"рюкзак":    true,
			}, Subjects: map[string]ObjSubj{
				"тумбочка": {Exist: true, Lock: true, Key: "хлопок в ладоши"},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
		"улица": {
			Name: "улица",
			Msg: map[string]*subjLock{
				"notlinked":  {lable: "нет пути улица"},
				"enter":      {lable: "на улице весна."},
				"lookaround": {lable: "ты находишься на улице, как же прекрасен свежый воздух"},
				"locked":     {lable: "дверь закрыта"},
			}, Subjects: map[string]ObjSubj{
				"шкаф": {Exist: true, Lock: true, Key: "кот"},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
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
