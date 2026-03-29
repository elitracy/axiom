package engine

import (
	"sync/atomic"
	"time"
)

const (
	TICKS_PER_SECOND = 1
	TICK_SLEEP       = time.Second / TICKS_PER_SECOND
)

type Tick struct {
	tick atomic.Int64
}

func (t *Tick) Tick() int64 { return t.tick.Load() }

func NewTick() *Tick { return &Tick{tick: atomic.Int64{}} }

type Game interface {
	Init()
	Update(tick *Tick)
}

func RunGame(game Game, startTick *Tick) {
	game.Init()
	for {
		game.Update(startTick)
		startTick.tick.Add(1)
		time.Sleep(TICK_SLEEP)
	}

}
