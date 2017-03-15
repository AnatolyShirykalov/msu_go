package game

import (
	"fmt"
	"strings"
)

var g Game

func (g *Game) GetLinks() (Ret []*Link) {
	Ret = make([]*Link, 0, len(g.Links))
	for i := range g.Links {
		Ret = append(Ret, &g.Links[i])
	}
	return
}

func (g *Game) GetRoom(name string) *Room {
	r, ok := g.Rooms[name]
	if ok {
		return r
	} else {
		alias, ok1 := g.Aliases[name]
		if ok1 {
			return g.GetRoom(alias)
		} else {
			panic(fmt.Sprintf("Не могу найти комнату по ключу %s", name))
		}
	}
}

func handleCommand(command string) string {
	gamer := &g.Players[0]
	c := strings.Split(command, " ")
	switch c[0] {
	case "осмотреться":
		return gamer.View()
	case "идти":
		return gamer.MoveTo(g.GetRoom(c[1]))
	case "одеть":
		{
			if c[1] != "рюкзак" {
				panic("Нельзя надеть " + c[1])
			}
			return gamer.AddBack(c[1])
		}
	case "взять":
		{
			if c[1] == "рюкзак" {
				panic("Нельзя взять " + c[1])
			}
			if gamer.RefBack == nil {
				return "некуда класть"
			}
			return gamer.AddThing(c[1])
		}
	case "применить":
		{
			return gamer.Apply(c[1], c[2])
		}
	default:
		return "неизвестная команда"
	}
}

func initGame() {
	g = Game{
		Rooms:   make(map[string]*Room),
		Links:   make([]Link, 0, 20),
		Players: make([]Player, 0, 1),
		Aliases: make(map[string]string),
	}
	g.Rooms = map[string]*Room{
		"кухня": {
			Game: &g,
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
				"ничего": {Exist: true, Lock: true}},
		},
		"коридор": {
			Game: &g,
			Name: "коридор",
			Msg: map[string]string{
				"notlinked": "нет пути коридор",
				"enter":     "ничего интересного.",
				"lookaround": "	Ты нахдишься в коридоре, тут страшно",
			}, Things: map[string]bool{
				"ничего": true,
			}, Subjects: map[string]ObjSubj{
				"шкаф": {Exist: true, Lock: true, Key: "кот"}},
		},
		"комната": {
			Game: &g,
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
				"тумбочка": {Exist: true, Lock: true, Key: "хлопок в ладоши"}},
		},
		"улица": {
			Game: &g,
			Name: "улица",
			Msg: map[string]string{
				"notlinked":  "нет пути улица",
				"enter":      "на улице весна.",
				"lookaround": "ты находишься на улице, как же прекрасен свежый воздух",
				"locked":     "дверь закрыта",
			}, Subjects: map[string]ObjSubj{
				"дверь": {Exist: true, Lock: true, Key: "ключи"},
				"шкаф":  {Exist: true, Lock: true, Key: "кот"}},
		},
	}

	g.Aliases = map[string]string{
		"домой": "коридор",
	}
	g.Players = []Player{
		{InRoom: g.GetRoom("кухня")},
	}

	g.Links = []Link{
		{Rfrom: g.GetRoom("кухня"), Rto: g.GetRoom("коридор")},
		{Rfrom: g.GetRoom("коридор"), Rto: g.GetRoom("кухня")},
		{Rfrom: g.GetRoom("коридор"), Rto: g.GetRoom("комната")},
		{Rfrom: g.GetRoom("коридор"), Rto: g.GetRoom("улица"), Lock: true},
		{Rfrom: g.GetRoom("комната"), Rto: g.GetRoom("коридор")},
		{Rfrom: g.GetRoom("улица"), Rto: g.GetRoom("коридор"), Name: "домой"},
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
