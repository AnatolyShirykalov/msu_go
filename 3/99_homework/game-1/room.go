package main

func DefaultFunc(r *Room, pl *Player) string {
	return ""
}
func InitRoom(nameRoom string) *Room {
	return &Room{
		Name: nameRoom,
		Act:  "идти в универ. ",
		Msg: map[string]string{
			"enter":      nameRoom + ", ничего интересного.",
			"lookaround": "ты находишься на кухне, на столе ",
			"backact":    "надо собрать ",
			"end":        "можно пройти - коридор",
		}, Things: map[string]bool{
			"чай": true,
		}, Subjects: map[string]ObjSubj{
			"ничего": {Exist: true, Lock: true},
		}, LinkRoom: map[string]string{
			"коридор": "exist"},
	}
}
