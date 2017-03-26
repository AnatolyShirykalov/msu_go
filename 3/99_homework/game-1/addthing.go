package main

import "fmt"

func (p *Player) AddThing(thi string) {
	room := p.InRoom
	if !room.Things[thi] {
		p.Msg.ChanMsg <- "нет такого"
	} else {
		room.Things[thi] = false
		p.Msg.ChanMsg <- p.RefBack.Take(thi)
	}
	G.wg.Done()
}

func (p *Player) AddBack(thi string) {
	room := p.InRoom
	if !room.Things[thi] {
		p.Msg.ChanMsg <- "нет такого"
	} else {
		p.RefBack = &Back{
			Things: map[string]bool{
				"ключи":     false,
				"конспекты": false,
			},
		}
		room.Things[thi] = false
		p.Msg.ChanMsg <- fmt.Sprintf("вы надели: %s", thi)
	}
	G.wg.Done()
}

func (p *Player) Apply(thi string, subj string) {
	room := p.InRoom
	//Appliabl
	if p.RefBack.Apply(thi, subj) == "" {
		if room.LinkRoom[room.Subjects[subj].NameRoom] == "lock" {
			room.LinkRoom[room.Subjects[subj].NameRoom] = "exist"
		}
		p.Msg.ChanMsg <- fmt.Sprintf("%s открыта", subj)
	} else {
		p.Msg.ChanMsg <- p.RefBack.Apply(thi, subj)
	}
	G.wg.Done()
}
