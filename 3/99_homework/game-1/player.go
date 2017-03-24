package game

import (
	"fmt"
	"sort"
	"strings"
)

// players := map[string]*Player{
// 	"Tristan": NewPlayer("Tristan"),
// 	"Izolda":  NewPlayer("Izolda"),
// }
// go func() {
// 	output := players["Tristan"].GetOutput()
// 	for lastOutput["Tristan"] = range output {
// 	}
// }()

func (p *Player) GetOutput() chan string {
	return p.msg
}

func (p *Player) don(command string) {
	c := strings.Split(command, " ")
	switch c[0] {
	case "осмотреться":
		fmt.Println("ghbj")
		p.msg <- p.View()
		return
	case "идти":
		fmt.Println(p.MoveTo(GetRoom(c[1])))
		p.msg <- p.MoveTo(GetRoom(c[1]))
		return
	case "надеть":
		if c[1] != "рюкзак" {
			panic("Нельзя надеть " + c[1])
		}
		p.msg <- p.AddBack(c[1])
		return
	case "взять":
		if c[1] == "рюкзак" {
			panic("Нельзя взять " + c[1])
		}
		if p.RefBack == nil {
			p.msg <- "некуда класть"
		} else {
			p.msg <- p.AddThing(c[1])

		}
		return
	case "применить":
		p.msg <- p.Apply(c[1], c[2])
		return
	case "сказать":
		p.msg <- p.Say(c[1:])
		return
	case "сказать_игроку":
		p.msg <- p.Tell(c[1:])
		return
	default:
		p.msg <- "неизвестная команда"
		return
	}
}

// у каждого игрока есть метод, который запускает рутину и возращает канал в которой будет возвращать ответы
func (p *Player) Say(comma []string) string {

	return ""
}

func (p *Player) Tell(command []string) string {

	if G.Players[command[0]] == nil {
		return "тут нет такого игрока"
	}

	if len(command) == 1 {
		return p.Name + " выразительно молчит, смотря на вас"
	}

	return p.Name + "говорит вам: " + strings.Join(command, " ")
}

func (rfro *Room) Passability(rto *Room) string {
	Flag, ok := rfro.LinkRoom[rto.Name]
	if !ok {
		msg, ok := rto.Msg["notlinked"]
		if !ok {
			panic(HaveNotMsg(rfro, "notlinked"))
		} else {
			return msg
		}
	}
	if Flag == "lock" {
		if msg, ok := rto.Msg["locked"]; !ok {
			panic(HaveNotMsg(rto, "locked"))
		} else {
			return msg
		}
	}
	if _, ok := rto.Msg["enter"]; !ok {
		panic(HaveNotMsg(rto, "enter"))
	}
	return ""
}

func (p *Player) MoveTo(r *Room) string {
	passability := p.InRoom.Passability(r)
	if passability == "" {
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
		if p.InRoom.Act == "" {

			msg += p.InRoom.Msg["backact"] + p.RefBack.Type() + ". "
		} else {
			msg += p.InRoom.Msg["backact"] + p.RefBack.Type() + " и " + p.InRoom.Act
		}
	}
	if flag == 0 {
		msg = "пустая комната. "
	}
	msg += p.InRoom.Msg["end"]
	// if len(G.Players) > 1 {
	// fmt.Println("dd")

	for name, exist := range G.Players {
		if exist != nil {
			if name != p.Name {
				msg += ". Кроме вас тут ещё " + name
				break
			}
		}

	}
	// fmt.Println(msg)
	// }
	return msg
}
