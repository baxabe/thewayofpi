package main

import (
	"fmt"
	"math"
	tb "github.com/nsf/termbox-go"
	"flag"
)

const (
	ConWidth =	100
	ConHigh =	50
	MaxSteep =	1000000
	MaxSleep =	1000000
	MaxRate =	1000000
)

var (
	task = flag.Int("i", 0, "Número de iteraciones")
	steep = flag.Int("st", 100, "Número de pasos para mostrar resultados")
	sleep = flag.Int("sl", 200, "Milisegundos de espera")
	show		bool
)

func main() {
    flag.Parse()
	if err := tb.Init(); err != nil {
		panic(err)
	}
	defer tb.Close()
	var wx, wy int
	if wx, wy = tb.Size(); wx < ConWidth && wy < ConHigh {
		panic(fmt.Sprintf("Poca pantalla. wx = %d, wy = %d", wx, wy))
	}
	initRnd()
	initTreelist()
	scr = screen {
		wx:			wx,
		wy:			wy,
		posJug:		point{0, 0},
		posDims:	point{0, 3},
		posTodo:	point{10, 3},
		posDone:	point{25, 3},
		posRoots:	point{40, 3},
		posCount:	point{0, 1},
		posConf:	point{wx - 50, 0},
	}
  	events := make(chan tb.Event)
	go func() {
		for {
			events <- tb.PollEvent()
		}
	}()
	iter := 0
	pulse := 0
	show = true
	rate = point{1, 1}
	showLabels()
	showConf()
loop:
	for {
		//logn(pulse)
		flush()
		select {
		case ev := <-events:
			if ev.Type == tb.EventKey {
				switch {
				case ev.Key == tb.KeyEsc:
					break loop
				case ev.Key == tb.KeySpace:
					show = !show
				case ev.Key == tb.KeyArrowUp:
					if *steep > 0 {
						n := math.Trunc(math.Log10(float64(*steep)))
						*steep = *steep + int(math.Trunc(math.Max(1, math.Pow10(int(n)))))
						if *steep > MaxSteep {
							*steep = MaxSteep
						}
					} else {
						*steep = 1
					}
				case ev.Key == tb.KeyArrowDown:
					if *steep > 0 {
						n := math.Trunc(math.Log10(float64(*steep)))
						m := math.Trunc(math.Pow10(int(n)))
						if float64(*steep) == m {
							m /= 10
						}
						*steep = *steep - int(math.Trunc(math.Max(1, m)))
						if *steep < 1 {
							*steep = 1
						}
					}
				case ev.Key == tb.KeyArrowRight:
					if *sleep > 0 {
						n := math.Trunc(math.Log10(float64(*sleep)))
						*sleep = *sleep + int(math.Trunc(math.Max(1, math.Pow10(int(n)))))
						if *sleep > MaxSleep {
							*sleep = MaxSleep
						}
					} else {
						*sleep = 1
					}
				case ev.Key == tb.KeyArrowLeft:
					if *sleep > 0 {
						n := math.Trunc(math.Log10(float64(*sleep)))
						m := math.Trunc(math.Pow10(int(n)))
						if float64(*sleep) == m {
							m /= 10
						}
						*sleep = *sleep - int(math.Trunc(math.Max(1, m)))
						if *sleep < 0 {
							*sleep = 0
						}
					}
				case ev.Ch == 'r':
					if rate.x > 0 {
						n := math.Trunc(math.Log10(float64(rate.x)))
						rate.x = rate.x + int(math.Trunc(math.Max(1, math.Pow10(int(n)))))
						if rate.x > MaxRate {
							rate.x = MaxRate
						}
					} else {
						rate.x = 1
					}
				case ev.Ch == 'e':
					if rate.x > 0 {
						n := math.Trunc(math.Log10(float64(rate.x)))
						m := math.Trunc(math.Pow10(int(n)))
						if float64(rate.x) == m {
							m /= 10
						}
						rate.x = rate.x - int(math.Trunc(math.Max(1, m)))
						if rate.x < 0 {
							rate.x = 0
						}
					}
				case ev.Ch == 'f':
					if rate.y > 0 {
						n := math.Trunc(math.Log10(float64(rate.y)))
						rate.y = rate.y + int(math.Trunc(math.Max(1, math.Pow10(int(n)))))
						if rate.y > MaxRate {
							rate.y = MaxRate
						}
					} else {
						rate.y = 1
					}
				case ev.Ch == 'd':
					if rate.y > 0 {
						n := math.Trunc(math.Log10(float64(rate.y)))
						m := math.Trunc(math.Pow10(int(n)))
						if float64(rate.x) == m {
							m /= 10
						}
						rate.y = rate.y - int(math.Trunc(math.Max(1, m)))
						if rate.y < 0 {
							rate.y = 0
						}
					}
				}
				if rate.x > rate.y {
					rate.y = rate.x
				}
				showConf()
			}
		default:
			if iter <= *task {
				doCycle()
  				pulse++
				showPulses(pulse)
			}
			if pulse % *steep == 0 {
				showData()
			}
			if *task > 0 {
				iter++
			}
	 	}
	}
}
