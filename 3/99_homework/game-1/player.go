package main

import (
	"fmt"
	"sort"
	"strings"
)

func (p *Player) HandleInput(command string) {
	// fmt.Println("handl")
	G.wg.Add(2)
	G.msgin <- &Command{
		command: command,
		player:  p,
	}
	G.wg.Wait()
}

func (p *Player) GetOutput() chan string {
	return p.msg.msg
}

func (p *Player) don(command string) {
	p.msg.Lock()
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
			p.msg.msg <- "некуда класть"
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
		p.msg.msg <- "неизвестная команда"
		G.wg.Done()
	}
	defer p.msg.Unlock()
}

// у каждого игрока есть метод, который запускает рутину и возращает канал в которой будет возвращать ответы
func (p *Player) Say(command []string) {
	for name, p2 := range G.Players {
		if p.InRoom.Name == p2.InRoom.Name && name != p.Name {
			p2.msg.msg <- p.Name + " говорит: " + strings.Join(command, " ")
			p.msg.msg <- p.Name + " говорит: " + strings.Join(command, " ")
			G.wg.Done()
		}
	}
	return
}

func (p *Player) Tell(command []string) {
	if G.Players[command[0]] == nil {
		p.msg.msg <- "тут нет такого игрока"
		G.wg.Done()
		return
	}
	if len(command) == 1 {
		if p.InRoom.Name == G.Players[command[0]].InRoom.Name {
			G.Players[command[0]].msg.msg <- p.Name + " выразительно молчит, смотря на вас"
			G.wg.Done()
			return
		}
		p.msg.msg <- "тут нет такого игрока"
		G.wg.Done()
		return
	}
	G.Players[command[0]].msg.msg <- p.Name + " говорит вам: " + strings.Join(command[1:], " ")
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
		p.msg.msg <- p.InRoom.Passability(r)
	} else {
		p.InRoom = r
		msg := r.Msg["enter"]
		var keys []string
		for k, value := range r.LinkRoom {
			if value != "" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		if len(keys) == len(G.Priory) {
			keys = G.Priory
		}
		flag := 0
		for _, link := range keys {
			if flag == 0 {
				msg = fmt.Sprintf("%s можно пройти - ", msg)
			} else {
				msg = fmt.Sprintf("%s, ", msg)
			}
			name := link
			if Aliase, ok := G.Aliases[r.Name]; ok {
				name = Aliase
			}
			msg = fmt.Sprintf("%s%s", msg, name)
			flag += 1
		}
		p.msg.msg <- msg
	}
	G.wg.Done()
}

func (p *Player) View() {
	msg := ""
	var keys []string
	for k, value := range p.InRoom.Things {
		if value && k != p.RefBack.Type() {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	flag := 0
	for _, k := range keys {
		value := p.InRoom.Things[k]
		if flag == 0 {
			msg += p.InRoom.Msg["lookaround"]
			flag = 1
		}
		if value {
			if flag < len(keys) {
				msg += k + ", "
				flag += 1
			} else {
				if p.RefBack == nil {
					msg += k + ", "
					flag += 1
				} else {
					if p.InRoom.Act == "" {
						msg += k + ". "
					} else {
						msg += k + ", надо " + p.InRoom.Act

					}
				}
			}
		}
	}
	if p.RefBack == nil {
		msg += p.InRoom.Msg["backact"] + p.RefBack.Type()
		if p.InRoom.Act == "" {
			msg += ". "
		} else {
			msg += " и " + p.InRoom.Act
		}
	}
	if flag == 0 {
		msg = "пустая комната. "
	}
	msg += p.InRoom.Msg["end"]

	for name, exist := range G.Players {
		if exist != nil {
			if name != p.Name {
				msg += ". Кроме вас тут ещё " + name
				break
			}
		}

	}
	p.msg.msg <- msg
	G.wg.Done()
}
