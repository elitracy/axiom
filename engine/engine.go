package engine

import (
	"time"
)

const (
	TICK_SLEEP = time.Second / 1
)

type Tick struct {
	tick int64
}

func (t *Tick) Tick() int64 { return t.tick }

func NewTick(startTick int64) *Tick { return &Tick{tick: startTick} }

type Game interface {
	Update(tick Tick)
}

func RunGame(game Game, startTick *Tick) {
	tick := startTick
	for {
		game.Update(*tick)
		tick.tick++
		time.Sleep(TICK_SLEEP)
	}

}
