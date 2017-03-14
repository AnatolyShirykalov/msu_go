package game

import (
	"fmt"
)

type Object interface {
	Label() string
	Type() string
}

func HaveNotMsg(o Object, key string) string {
	return fmt.Sprintf("Нет сообщения с ключом %s у %s %s", key, o.Type(), o.Label())
}
