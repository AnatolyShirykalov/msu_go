package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

func (p *Player) HandleInput(command string) {
	G.Msgin <- &Command{
		Command: command,
		Player:  p,
	}
	time.Sleep(time.Millisecond / 100)
}

func (p *Player) GetOutput() chan string {
	return p.Msg
}

func (p *Player) don(command string) {
	c := strings.Split(command, " ")
	switch c[0] {
	case "осмотреться":
		p.View()
	case "идти":
		p.MoveTo(GetRoom(c[1]))
	case "надеть":
		if c[1] != "" {
			if c[1] != "рюкзак" {
				p.Msg <- "Нельзя надеть " + c[1]
			} else {
				p.AddBack(c[1])
			}
		} else {
			p.Msg <- "не чего надеть"
		}
	case "взять":
		if c[1] != "" {
			if c[1] == "рюкзак" {
				p.Msg <- "Нельзя взять " + c[1]
			}
			if p.RefBack == nil {
				p.Msg <- "некуда класть"
			} else {
				p.AddThing(c[1])
			}
		} else {
			p.Msg <- "не чего взять"
		}
	case "применить":
		if c[2] != "" {
			p.Apply(c[1], c[2])
		} else {
			p.Msg <- "не к чему применить"
		}
	case "сказать":
		p.Say(c[1:])
	case "сказать_игроку":
		p.Tell(c[1:])
	default:
		p.Msg <- "неизвестная команда"
	}
}

// у каждого игрока есть метод, который запускает рутину и возращает канал в которой будет возвращать ответы
func (p *Player) Say(command []string) {
	for _, p2 := range GetRoom(p.InRoom).Players {
		// fmt.Println("...")
		// if name != p.Name {
		p2.Msg <- p.Name + " говорит: " + strings.Join(command, " ")
		// p.Msg <- p.Name + " говорит: " + strings.Join(command, " ")
		// }
		log.Println("User registered " + p2.Name + p.Name)
	}
	return
}
func (p *Player) Tell(command []string) {
	RoomPlayer := GetRoom(p.InRoom)
	if RoomPlayer.Players[command[0]] == nil {
		p.Msg <- "тут нет такого игрока"
		return
	}
	if len(command) == 1 {
		RoomPlayer.Players[command[0]].Msg <- p.Name + " выразительно молчит, смотря на вас"
		return
	}
	RoomPlayer.Players[command[0]].Msg <- p.Name + " говорит вам: " + strings.Join(command[1:], " ")
	return
}

func (rfro *Room) notPassability(rto *Room) string {
	Flag, ok := rfro.LinkRoom[rto.Name]
	if !ok {
		notlinked := "нет пути "
		if rto.Name == "комната" {
			notlinked += "в "
		}
		return notlinked + rto.Name
	}
	if Flag == "lock" {
		return rto.Msg["locked"]
	}
	return ""
}

func (p *Player) MoveTo(r *Room) {
	if GetRoom(p.InRoom).notPassability(r) != "" {
		p.Msg <- GetRoom(p.InRoom).notPassability(r)
		return
	}
	G.Rooms[p.InRoom].Players[p.Name] = nil
	G.Rooms[r.Name].Players[p.Name] = p
	p.InRoom = r.Name
	msg := r.Msg["enter"]
	var rooms []string
	for roomin, roomout := range r.LinkRoom {
		if roomout != "" {
			rooms = append(rooms, roomin)
		}
	}
	sort.Strings(rooms)
	if len(rooms) == len(G.Priory) {
		rooms = G.Priory
	}
	for i, link := range rooms {
		if i == 0 {
			msg = fmt.Sprintf("%s можно пройти - ", msg)
		} else {
			msg = fmt.Sprintf("%s, ", msg)
		}
		name := link
		if Aliase, ok := G.Aliases[r.Name]; ok {
			name = Aliase
		}
		msg = fmt.Sprintf("%s%s", msg, name)
	}
	p.Msg <- msg
}

func (p *Player) View() {
	msg := ""
	var things []string
	RoomPlayer := GetRoom(p.InRoom)
	for thing, exist := range RoomPlayer.Things {
		if exist && thing != p.RefBack.Type() {
			things = append(things, thing)
		}
	}
	sort.Strings(things)
	msg += RoomPlayer.Msg["lookaround"]
	for i, thing := range things {
		if i < len(things)-1 {
			msg += fmt.Sprintf("%s, ", thing)
			continue
		}
		if p.RefBack == nil {
			msg += fmt.Sprintf("%s, ", thing)
		} else {
			if RoomPlayer.Act == "" {
				msg += fmt.Sprintf("%s. ", thing)
			} else {
				msg += fmt.Sprintf("%s, надо %s", thing, RoomPlayer.Act)
			}
		}
	}
	if len(RoomPlayer.Things) != 0 {
		if p.RefBack == nil {
			msg += RoomPlayer.Msg["backact"] + p.RefBack.Type()
			if GetRoom(p.InRoom).Act == "" {
				msg += ". "
			} else {
				msg += fmt.Sprintf(" и %s", RoomPlayer.Act)
			}
		}
		if len(things) == 0 {
			msg = "пустая комната. "
		}
	}
	msg += RoomPlayer.Msg["end"]

	flag := true
	for name, exist := range RoomPlayer.Players {
		if exist != nil {
			if name != p.Name {
				if flag {
					msg += fmt.Sprintf(". кроме вас тут ещё %s", name)
					flag = false
				} else {
					msg += fmt.Sprintf(", %s", name)
				}
			}
		}
	}
	p.Msg <- msg
}
