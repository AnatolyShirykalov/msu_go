package game

import (
	"fmt"
	"sort"
)

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
		if len(keys) == len(p.InRoom.Game.Priory) {
			keys = p.InRoom.Game.Priory
		}
		flag := 0
		for _, link := range keys {
			if flag == 0 {
				msg = fmt.Sprintf("%s можно пройти - ", msg)
			} else {
				msg = fmt.Sprintf("%s, ", msg)
			}
			name := link
			if Aliase, ok := r.Game.Aliases[r.Name]; ok {
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
	return msg
}
