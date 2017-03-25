package main

func (p *Player) AddThing(thi string) string {
	room := p.InRoom
	if !room.Things[thi] {
		return "нет такого"
	}
	room.Things[thi] = false
	return p.RefBack.Take(thi)
}

func (p *Player) AddBack(thi string) {
	room := p.InRoom
	if !room.Things[thi] {
		p.msg <- "нет такого"
	} else {
		p.RefBack = &Back{
			Things: map[string]bool{
				"ключи":     false,
				"конспекты": false,
			},
		}
		room.Things[thi] = false
		p.msg <- "вы надели: " + thi
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
		p.msg <- subj + " открыта"
	} else {
		p.msg <- p.RefBack.Apply(thi, subj)
	}
	G.wg.Done()
}
