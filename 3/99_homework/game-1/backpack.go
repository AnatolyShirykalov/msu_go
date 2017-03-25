package main

func (b *Back) Type() string {
	return "рюкзак"
}

func (b *Back) Apply(thi string, subj string) string {
	if b == nil {
		return "нет предмета в инвентаре - " + thi

	}
	if _, ok := b.Things[thi]; !ok {
		return "нет предмета в инвентаре - " + thi
	}
	if subj != "дверь" {
		return "не к чему применить"
	}
	return ""
}

func (b *Back) Take(thi string) string {
	if b.Things[thi] {
		return thi + " уже есть"
	}
	b.Things[thi] = true
	flag := 0
	for _, value := range b.Things {
		if value {
			flag += 1
		}
	}
	if flag == 2 {
		b.Full = true
	}
	return "предмет добавлен в инвентарь: " + thi
}
