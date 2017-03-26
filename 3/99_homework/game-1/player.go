package main

import (
	"fmt"
	"sort"
	"strings"
	// "time"
)

func (p *Player) HandleInput(command string) {
	// fmt.Println("handl")
	G.wg.Add(2)
	G.msgin <- &Command{
		command: command,
		player:  p,
	}
	// time.Sleep(time.Millisecond)
	G.wg.Wait()
}

func (p *Player) GetOutput() chan string {
	return p.Msg.ChanMsg
}

func (p *Player) don(command string) {
	p.Msg.Lock()
	c := strings.Split(command, " ")
	switch c[0] {
	case "осмотреться":
		p.View()
	case "идти":
		p.MoveTo(GetRoom(c[1]))
	case "надеть":
		if c[1] != "рюкзак" {
			panic("Нельзя надеть " + c[1])
		}
		p.AddBack(c[1])
	case "взять":
		if c[1] == "рюкзак" {
			panic("Нельзя взять " + c[1])
		}
		if p.RefBack == nil {
			p.Msg.ChanMsg <- "некуда класть"
			G.wg.Done()
		} else {
			p.AddThing(c[1])
		}
	case "применить":
		p.Apply(c[1], c[2])
	case "сказать":
		p.Say(c[1:])
	case "сказать_игроку":
		p.Tell(c[1:])
	default:
		p.Msg.ChanMsg <- "неизвестная команда"
		G.wg.Done()
	}
	defer p.Msg.Unlock()
}

// у каждого игрока есть метод, который запускает рутину и возращает канал в которой будет возвращать ответы
func (p *Player) Say(command []string) {
	for name, p2 := range G.Players {
		if p.InRoom.Name == p2.InRoom.Name && name != p.Name {
			p2.Msg.ChanMsg <- p.Name + " говорит: " + strings.Join(command, " ")
			p.Msg.ChanMsg <- p.Name + " говорит: " + strings.Join(command, " ")
			G.wg.Done()
		}
	}
	return
}

func (p *Player) Tell(command []string) {
	if G.Players[command[0]] == nil {
		p.Msg.ChanMsg <- "тут нет такого игрока"
		G.wg.Done()
		return
	}
	if len(command) == 1 {
		if p.InRoom.Name == G.Players[command[0]].InRoom.Name {
			G.Players[command[0]].Msg.ChanMsg <- p.Name + " выразительно молчит, смотря на вас"
			G.wg.Done()
			return
		}
		p.Msg.ChanMsg <- "тут нет такого игрока"
		G.wg.Done()
		return
	}
	G.Players[command[0]].Msg.ChanMsg <- p.Name + " говорит вам: " + strings.Join(command[1:], " ")
	G.wg.Done()
	return
}

func (rfro *Room) Passability(rto *Room) string {
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
	if p.InRoom.Passability(r) != "" {
		p.Msg.ChanMsg <- p.InRoom.Passability(r)
	} else {
		p.InRoom = r
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
		p.Msg.ChanMsg <- msg
	}
	G.wg.Done()
}

func (p *Player) View() {
	msg := ""
	var things []string
	for thing, exist := range p.InRoom.Things {
		if exist && thing != p.RefBack.Type() {
			things = append(things, thing)
		}
	}
	sort.Strings(things)
	for i, thing := range things {
		if i == 0 {
			msg += p.InRoom.Msg["lookaround"]
		}
		if i < len(things)-1 {
			msg += fmt.Sprintf("%s, ", thing)
			continue
		}
		if p.RefBack == nil {
			msg += fmt.Sprintf("%s, ", thing)
		} else {
			if p.InRoom.Act == "" {
				msg += fmt.Sprintf("%s. ", thing)
			} else {
				msg += fmt.Sprintf("%s, надо %s", thing, p.InRoom.Act)
			}
		}
	}
	if p.RefBack == nil {
		msg += p.InRoom.Msg["backact"] + p.RefBack.Type()
		if p.InRoom.Act == "" {
			msg += ". "
		} else {
			msg += fmt.Sprintf(" и %s", p.InRoom.Act)
		}
	}
	if len(things) == 0 {
		msg = "пустая комната. "
	}
	msg += p.InRoom.Msg["end"]

	for name, exist := range G.Players {
		if exist != nil {
			if name != p.Name {
				msg += fmt.Sprintf(". Кроме вас тут ещё %s", name)
				break
			}
		}
	}
	p.Msg.ChanMsg <- msg
	G.wg.Done()
}
