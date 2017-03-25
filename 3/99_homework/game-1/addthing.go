package main

func (p *Player) AddThing(thi string) {
	room := p.InRoom
	if !room.Things[thi] {
		p.msg.msg <- "нет такого"
	} else {
		room.Things[thi] = false
		p.msg.msg <- p.RefBack.Take(thi)
	}
	G.wg.Done()
}

func (p *Player) AddBack(thi string) {
	room := p.InRoom
	if !room.Things[thi] {
		p.msg.msg <- "нет такого"
	} else {
		p.RefBack = &Back{
			Things: map[string]bool{
				"ключи":     false,
				"конспекты": false,
			},
		}
		room.Things[thi] = false
		p.msg.msg <- "вы надели: " + thi
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
		p.msg.msg <- subj + " открыта"
	} else {
		p.msg.msg <- p.RefBack.Apply(thi, subj)
	}
	G.wg.Done()
}
