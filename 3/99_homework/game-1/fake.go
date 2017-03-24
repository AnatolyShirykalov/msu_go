package main

import (
	"fmt"
)

var G Game

type Player struct {
	msg chan string
}

type Game struct {
	msgin chan string
}

func initGame() {
	G = Game{
		msgin: make(chan string),
	}
}
func NewPlayer(name string) *Player {
	G.Players = map[string]*Player{}
}
func addPlayer(name string) {

}

func (p *Player) HandleInput(msg string) {
	G.msgin <- msg
}

func (p *Player) GetOutput() chan string {
	return p.msg
}
