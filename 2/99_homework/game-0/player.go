package game

import (
	"fmt"
	"sort"
)

type Player struct {
	InRoom  *Room
	RefBack *Back
}

func (p *Player) Apply(thi string, subj string) string {
	room := p.InRoom
	if room.Things[thi] {
		return "не к чему применить"
	}
	if p.RefBack.Apply(thi, subj) == "" {
		flag := -1
		for i, link := range p.InRoom.Links() {
			if link.Rfrom.Name == room.Name && link.Lock {
				flag = i
			}
		}
		if flag > -1 {
			room.Game.Links[flag+1].Lock = false
		}
		return subj + " открыта"
	} else {
		return p.RefBack.Apply(thi, subj)
	}
}

func (p *Player) AddThing(thi string) string {
	room := p.InRoom
	if !room.Things[thi] {
		return "нет такого"
	}
	room.Things[thi] = false
	return p.RefBack.Take(thi)
}

func (p *Player) AddBack(thi string) string {
	room := p.InRoom
	if !room.Things[thi] {
		return "нет такого"
	}
	p.RefBack = &Back{
		Things: map[string]bool{
			"ключи":     false,
			"конспекты": false,
		},
	}
	room.Things[thi] = false
	return "вы одели: " + thi
}

func (p *Player) MoveTo(r *Room) string {
	if !p.InRoom.LinkedWith(r) {
		msg, ok := r.Msg["notlinked"]
		if !ok {
			panic(HaveNotMsg(p.InRoom, "notlinked"))
		} else {
			return msg
		}
	}
	if !p.InRoom.UnlockedLinkTo(r) {
		if msg, ok := r.Msg["locked"]; !ok {
			panic(HaveNotMsg(r, "locked"))
		} else {
			return msg
		}
	}
	if msg, ok := r.Msg["enter"]; !ok {
		panic(HaveNotMsg(r, "enter"))
	} else {
		// p.InRoom.Watch = false
		p.InRoom = r
		for i, link := range r.Links() {
			if i == 0 {
				msg = fmt.Sprintf("%s можно пройти - ", msg)
			} else {
				msg = fmt.Sprintf("%s, ", msg)
			}
			name := link.Name
			if len(name) == 0 {
				name = link.Rto.Name
			}
			msg = fmt.Sprintf("%s%s", msg, name)
		}
		return msg
	}
}

func (p *Player) View() string {
	msg := ""
	var keys []string
	flag := 0
	for k, value := range p.InRoom.Things {
		if value && k != p.RefBack.Type() {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
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
