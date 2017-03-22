package game

import (
	"fmt"
	"strings"
)

var G Game

// addPlayer(players["Tristan"])
// g := Game{
// 	Rooms: make(map[string]*Room),
// 	// Links:   make([]Link, 0, 20),
// 	Players: make(map[string]*Player),
// 	Aliases: make(map[string]string),
// }
// addPlayer(players["Izolda"])
//создание игрока
func addPlayer(p *Player) {
	// G.Players[p.Name] = p
	// fmt.Println(p.Name)
	G.Players[p.Name] = p
	p.InRoom = GetRoom("кухня")
	// fmt.Println(len(G.Players))
}
func GetRoom(name string) *Room {
	// fmt.Println(G.Rooms[name].Name)
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

// type Player struct {
// 	InRoom  *Room
// 	RefBack *Back
// 	msg     chan string
// 	Name    string
// }

func NewPlayer(name string) (player *Player) {
	G.Players = map[string]*Player{
		name: {
			Name: name,
			msg:  make(chan string),
		},
	}
	// G.Players[name] = &Player{
	// 	Name: name,
	// 	msg:  make(chan string),
	// }
	// G.Players[name].Name = name
	// fmt.Println(G.Players[na\])
	// G.Players[name]=&Player{Name:}
	// G.Players[name].msg = make(chan string)
	return G.Players[name]
}
func initGame() {
	G = Game{
		Rooms: make(map[string]*Room),
		// Links:   make([]Link, 0, 20),
		Players: make(map[string]*Player),
		Aliases: make(map[string]string),
	}
	G.Rooms = map[string]*Room{
		"кухня": {
			Name: "кухня",
			Act:  "идти в универ. ",
			Msg: map[string]string{
				"notlinked":  "нет пути кухня",
				"enter":      "кухня, ничего интересного.",
				"lookaround": "ты находишься на кухне, на столе ",
				"backact":    "надо собрать ",
				"end":        "можно пройти - коридор",
			}, Things: map[string]bool{
				"чай": true,
			}, Subjects: map[string]ObjSubj{
				"ничего": {Exist: true, Lock: true},
			}, LinkRoom: map[string]string{
				"коридор": "exist"},
		},
		"коридор": {
			Name: "коридор",
			Msg: map[string]string{
				"notlinked": "нет пути коридор",
				"enter":     "ничего интересного.",
				"lookaround": "	Ты нахдишься в коридоре, тут страшно",
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
			Msg: map[string]string{
				"notlinked":  "нет пути в комната",
				"enter":      "ты в своей комнате.",
				"lookaround": "на столе: ",
				"backact":    "на стуле - ",
				"end":        "можно пройти - коридор",
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
			Msg: map[string]string{
				"notlinked":  "нет пути улица",
				"enter":      "на улице весна.",
				"lookaround": "ты находишься на улице, как же прекрасен свежый воздух",
				"locked":     "дверь закрыта",
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
}
func (gamer *Player) HandleInput(command string) {
	// gamer.msg := make(chan string)
	c := strings.Split(command, " ")
	switch c[0] {
	case "осмотреться":
		// fmt.Println("ghbj")
		gamer.msg <- gamer.View()
		return
	case "идти":
		gamer.msg <- gamer.MoveTo(GetRoom(c[1]))
		return
	case "надеть":
		{
			if c[1] != "рюкзак" {
				panic("Нельзя надеть " + c[1])
			}
			gamer.msg <- gamer.AddBack(c[1])
		}
		return
	case "взять":
		{
			if c[1] == "рюкзак" {
				panic("Нельзя взять " + c[1])
			}
			if gamer.RefBack == nil {
				gamer.msg <- "некуда класть"
			} else {
				gamer.msg <- gamer.AddThing(c[1])

			}
			return
		}
	case "применить":
		{
			gamer.msg <- gamer.Apply(c[1], c[2])
			return
		}
	case "сказать":
		gamer.msg <- gamer.Say(c[1:])
		return
	case "сказать_игроку":
		gamer.msg <- gamer.Tell(c[1:])
		return
	default:
		gamer.msg <- "неизвестная команда"
		return
	}
}

// func Run() {
// 	initGame()
// 	p := g.Players[0]
// 	fmt.Println(p.InRoom.LinkedWith(g.GetRoom("коридор")))
// 	fmt.Println(p.MoveTo(g.GetRoom("коридор")))
// 	fmt.Println(p.MoveTo(g.GetRoom("улица")))
// 	fmt.Println(p.MoveTo(g.GetRoom("домой")))
// }
