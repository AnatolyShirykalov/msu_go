package game

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
	return "вы надели: " + thi
}

func (p *Player) Apply(thi string, subj string) string {
	room := p.InRoom
	//Appliabl
	if p.RefBack.Apply(thi, subj) == "" {
		if room.LinkRoom[room.Subjects[subj].NameRoom] == "lock" {
			room.LinkRoom[room.Subjects[subj].NameRoom] = "exist"
		}
		return subj + " открыта"
	} else {
		return p.RefBack.Apply(thi, subj)
	}
}
