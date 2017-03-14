package game

func (r *Room) Label() string {
	return r.Name
}

func (r *Room) Type() string {
	return "комнаты"
}

func (r *Room) Links() []Link {
	ret := make([]Link, 0, 4)
	for _, link := range r.Game.Links {
		if link.Rfrom.Name == r.Name {
			ret = append(ret, link)
		}
	}
	return ret
}

// permitLock == true
func (r *Room) LinkedWith(rto *Room) bool {
	return r.linkTo(rto, true)
}

// permitLock == false
func (r *Room) UnlockedLinkTo(rto *Room) bool {
	return r.linkTo(rto, false)
}

//Проверка того что есть связь между комнатами и нет никакого лока
func (r *Room) linkTo(rto *Room, permitLock bool) bool {
	for _, link := range r.Links() {
		if link.Rfrom.Name == r.Name && link.Rto.Name == rto.Name && (permitLock || !link.Lock) {
			// if link.Lock {
			// 	panic(r.Name + rto.Name)
			// }
			return true
		}
	}
	return false
}
