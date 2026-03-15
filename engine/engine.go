package engine

import (
	"sync/atomic"
	"time"
)

const (
	TICK_SLEEP = time.Second / 1
)

type Tick struct {
	tick atomic.Int64
}

func (t *Tick) Tick() int64 { return t.tick.Load() }

func NewTick() *Tick { return &Tick{tick: atomic.Int64{}} }

type Game interface {
	Update(tick *Tick)
}

func RunGame(game Game, startTick *Tick) {
	for {
		game.Update(startTick)
		startTick.tick.Add(1)
		time.Sleep(TICK_SLEEP)
	}

}
