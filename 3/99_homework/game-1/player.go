package game

import (
	"fmt"
	"sort"
	"strings"
)

func (p *Player) HandleInput(command string) {
	G.wg.Add(2)
	G.msgin <- &Command{
		command: command,
		player:  p,
	}
	G.wg.Wait()
}

func (p *Player) GetOutput() chan string {
	func() {
		for {
			select {
			case msg := <-p.msg:
				p.msgout <- msg
			}
		}
	}()
	return p.msgout
}

func (p *Player) don(command string) {
	c := strings.Split(command, " ")
	switch c[0] {
	case "осмотреться":
		p.msg <- p.View()
		G.wg.Done()
	case "идти":
		p.msg <- p.MoveTo(GetRoom(c[1]))
		G.wg.Done()
	case "надеть":
		if c[1] != "рюкзак" {
			panic("Нельзя надеть " + c[1])
		}
		p.msg <- p.AddBack(c[1])
		G.wg.Done()
	case "взять":
		if c[1] == "рюкзак" {
			panic("Нельзя взять " + c[1])
		}
		if p.RefBack == nil {
			p.msg <- "некуда класть"
		} else {
			p.msg <- p.AddThing(c[1])

		}
		G.wg.Done()

	case "применить":
		p.msg <- p.Apply(c[1], c[2])
		G.wg.Done()
	case "сказать":
		p.msg <- p.Say(c[1:])
		G.wg.Done()
	case "сказать_игроку":
		p.msg <- p.Tell(c[1:])
		G.wg.Done()
	default:
		p.msg <- "неизвестная команда"
		G.wg.Done()
	}
}

// у каждого игрока есть метод, который запускает рутину и возращает канал в которой будет возвращать ответы
func (p *Player) Say(command []string) string {
	msg := p.Name + " говорит: " + strings.Join(command, " ") + " "
	for name, p2 := range G.Players {
		if p.InRoom.Name == p2.InRoom.Name && name != p.Name {
			msg += name + " говорит: " + strings.Join(command, " ")
		}
	}
	return msg
}

func (p *Player) Tell(command []string) string {
	if G.Players[command[0]] == nil {
		return "тут нет такого игрока"
	}

	if len(command) == 1 {
		fmt.Println(command[0])
		if p.InRoom.Name == G.Players[command[0]].InRoom.Name {
			return p.Name + " выразительно молчит, смотря на вас"
		}
		return "тут нет такого игрока"
	}

	return p.Name + "говорит вам: " + strings.Join(command, " ")
}

func (rfro *Room) Passability(rto *Room) string {
	Flag, ok := rfro.LinkRoom[rto.Name]
	if !ok {
		return rto.Msg["notlinked"].lable
	}
	if Flag == "lock" {
		return rto.Msg["locked"].lable
	}
	return ""
}

func (p *Player) MoveTo(r *Room) string {
	passability := p.InRoom.Passability(r)
	if passability == "" {
		p.InRoom = r
		msg := r.Msg["enter"].lable
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
		return msg
	} else {
		return passability
	}
}

func (p *Player) View() string {
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
			msg += p.InRoom.Msg["lookaround"].lable
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
		msg += p.InRoom.Msg["backact"].lable + p.RefBack.Type()
		if p.InRoom.Act == "" {
			msg += ". "
		} else {
			msg += " и " + p.InRoom.Act
		}
	}
	if flag == 0 {
		msg = "пустая комната. "
	}
	msg += p.InRoom.Msg["end"].lable

	for name, exist := range G.Players {
		if exist != nil {
			if name != p.Name {
				msg += ". Кроме вас тут ещё " + name
				break
			}
		}

	}
	return msg
}
