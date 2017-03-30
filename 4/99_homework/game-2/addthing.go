package main

import "fmt"

func (p *Player) AddThing(thi string) {
	room := GetRoom(p.InRoom)
	if !room.Things[thi] {
		p.Msg <- "нет такого"
	} else {
		room.Things[thi] = false
		p.Msg <- p.RefBack.Take(thi)
	}
}

func (p *Player) AddBack(thi string) {
	room := GetRoom(p.InRoom)
	if !room.Things[thi] {
		p.Msg <- "нет такого"
	} else {
		p.RefBack = &Back{
			Things: map[string]bool{
				"ключи":     false,
				"конспекты": false,
			},
		}
		room.Things[thi] = false
		p.Msg <- fmt.Sprintf("вы надели: %s", thi)
	}
}

func (p *Player) Apply(thi string, subj string) {
	room := GetRoom(p.InRoom)
	//Appliabl
	if p.RefBack.Apply(thi, subj) == "" {
		if room.LinkRoom[room.Subjects[subj].NameRoom] == "lock" {
			room.LinkRoom[room.Subjects[subj].NameRoom] = "exist"
		}
		p.Msg <- fmt.Sprintf("%s открыта", subj)
	} else {
		p.Msg <- p.RefBack.Apply(thi, subj)
	}
}
