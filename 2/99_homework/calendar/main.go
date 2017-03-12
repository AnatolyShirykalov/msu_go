package main

import (
	// "fmt"
	"time"
)

type calendar struct {
	month int
	year  int
	day   int
}

func NewCalendar(t time.Time) (c calendar) {
	c.year = t.Year()
	c.month = int(t.Month())
	c.day = t.Day()
	return
}

func (c calendar) CurrentQuarter() (quart int) {
	m := c.month
	switch {
	case m < 0:
		panic("not exist month")
	case m < 4:
		quart = 1
	case m < 7:
		quart = 2
	case m < 10:
		quart = 3
	case m < 13:
		quart = 4
	default:
		panic("not exist month")
	}
	return
}
