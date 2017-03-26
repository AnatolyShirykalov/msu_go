package main

import "fmt"

func (b *Back) Type() string {
	return "рюкзак"
}

func (b *Back) Apply(thi string, subj string) string {
	if b == nil {
		return fmt.Sprintf("нет предмета в инвентаре - %s", thi)

	}
	if _, ok := b.Things[thi]; !ok {
		return fmt.Sprintf("нет предмета в инвентаре - %s", thi)
	}
	if subj != "дверь" {
		return "не к чему применить"
	}
	return ""
}

func (b *Back) Take(thi string) string {
	if b.Things[thi] {
		return fmt.Sprintf("%s уже есть", thi)
	}
	b.Things[thi] = true
	CountThi := 0
	for _, value := range b.Things {
		if value {
			CountThi += 1
		}
	}
	if CountThi == 2 {
		b.Full = true
	}
	return fmt.Sprintf("предмет добавлен в инвентарь: %s", thi)
}
