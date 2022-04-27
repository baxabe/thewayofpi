package main

import (
	"fmt"
	"time"
	"math"
	"math/rand"
	tb "github.com/nsf/termbox-go"
	"flag"
)

const (
	NPop =		5
	NAxisPop =	9
	ValWidth =	12
	InterLin =	12
	ConWidth =	160
	ConHigh =	58
	MaxDim =	40
	MaxWill =	2 * 1024
	MaxSize =	16 * 1024
	ExpAxisVal =	28
	ExpConfortBorder =	20
	ExpDesire =	24
	InterestBound = 1000
)

type (
	CanAxisStruct struct {
		Polar		float64
		Speed		float64
		Push		float64
		MaxVal	float64
		MinVal	float64
	}
	ValAxisStruct struct {
		Confort	float64
		Radix	float64
		TopConf	float64
		BotConf	float64
		ActVal	float64
		Desire	int64
		Impetus	float64
		weight	float64
		PreWill	float64
	}
	CalAxisStruct struct {
		Delta	float64
		When	int32
		Feje	float64
		Fwill	float64
		Froz	float64
		Rate	float64
		DelRat	float64
		weight	float64
		NewMax	float64
		NewMin	float64
		Desire	int
		Come	bool
	}
	Popul struct {
		Will		float64
		Size		float64
		InAtt		float64
		OutAtt	float64
		VIAAtt	map[int64]float64
		CanAxis	[NAxisPop]CanAxisStruct
		ValAxis	[NAxisPop]ValAxisStruct
		CalAxis	[NAxisPop]CalAxisStruct
	}
	NexusStruct struct {
		InfoDim	float64
		FysDim	float64
		Popul		[NPop]Popul
	}
)

var (
	Nexus		NexusStruct
	Rnd			*rand.Rand
	coord		[NPop][NAxisPop][2]int
	wx			int
	wy			int
	FlagAxis = flag.Bool("e", false, "Calcula solo la fuerza del eje")
	FlagIter = flag.Int("iter", 0, "Número de iteraciones")
	FlagSteep = flag.Int("steep", 100, "Número de pasos para mostrar resultados")
	FlagSleep = flag.Int("sleep", 200, "Milisegundos de espera")
	show		bool
)

func main() {
    flag.Parse()
	if err := tb.Init(); err != nil {
		panic(err)
	}
	defer tb.Close()

	if wx, wy = tb.Size(); wx >= ConWidth && wy >= ConHigh {
		initRnd()
	  initCoord()
 	  initNexus()
 	  
  	event_queue := make(chan tb.Event)
		go func() {
			for {
				event_queue <- tb.PollEvent()
			}
		}()
		iter := 0
		pulse := 0
		show = true
loop:
	  for {
			select {
			case ev := <-event_queue:
				if ev.Type == tb.EventKey && ev.Key == tb.KeyEsc {
					break loop
				} else if ev.Type == tb.EventKey && ev.Key == tb.KeySpace {
					show = !show
				} else if ev.Type == tb.EventKey && ev.Key == tb.KeyArrowUp {
					*FlagSteep = *FlagSteep + 100
					if *FlagSteep > 100000 {
						*FlagSteep = 100000
					}
				} else if ev.Type == tb.EventKey && ev.Key == tb.KeyArrowDown {
					*FlagSteep = *FlagSteep - 100
					if *FlagSteep <= 0 {
						*FlagSteep = 1
					}
				} else if ev.Type == tb.EventKey && ev.Key == tb.KeyArrowRight {
					*FlagSleep = *FlagSleep + 100
					if *FlagSleep < 0 {
						*FlagSleep = 0
					}
				} else if ev.Type == tb.EventKey && ev.Key == tb.KeyArrowLeft {
					*FlagSleep = *FlagSleep - 100
					if *FlagSleep > 100000 {
						*FlagSleep = 100000
					}
				}
			default:
				if iter <= *FlagIter {
					if *FlagAxis {
 			 			calcFEje()
 			 		} else {
 						calcFEndo(pulse)
	  					pulse++
 					}
   					if pulse % *FlagSteep == 0 {
   						printWorld(pulse, *FlagAxis)
  					}
  					if *FlagIter > 0 {
  						iter++
  					}
	  			}
			}
		}
	}
}

func initRnd() {
	t := time.Date(1969, time.January, 1, 0, 0, 0, 0, time.UTC)
	d := time.Since(t)
	s := rand.NewSource(d.Nanoseconds())
	Rnd = rand.New(s)
}

func initCoord() {
	for i := 0; i < NPop; i++ {
		for j := 0; j < NAxisPop; j++ {
			coord[i][j][0] = 8 + j * ValWidth
			coord[i][j][1] = i * InterLin
		}
	}
}

func initNexus() {
	Nexus.InfoDim = float64(Rnd.Intn(MaxDim) + 1)
	Nexus.FysDim = float64(int(Nexus.InfoDim) + Rnd.Intn(8))
	for i := 0; i < NPop; i++ { //Cada Poblacion
		Nexus.Popul[i].Will	= float64(Rnd.Intn(MaxWill) + 1)
		Nexus.Popul[i].Size = float64(Rnd.Intn(MaxSize) + 1)
		//Nexus.Popul[i].InAtt = float64(??)
		//Nexus.Popul[i].OutAtt = float64(??)
		for j := 0; j < NAxisPop; j++ { //Cada Eje
			Nexus.Popul[i].CanAxis[j].Polar = float64(1 - Rnd.Intn(3))
			Nexus.Popul[i].CanAxis[j].Speed = float64(512*16 / (Rnd.Intn(512) + 1))
			Nexus.Popul[i].CanAxis[j].Push = float64((Nexus.InfoDim + Nexus.FysDim) / 2)
			Nexus.Popul[i].CalAxis[j].Feje = math.Trunc(Nexus.Popul[i].CanAxis[j].Polar * Nexus.Popul[i].CanAxis[j].Speed * Nexus.Popul[i].CanAxis[j].Push)
			gap := toint32(1<<(ExpAxisVal-2))
			Nexus.Popul[i].ValAxis[j].PreWill = Nexus.Popul[i].CanAxis[j].Speed * Nexus.Popul[i].Will * Nexus.Popul[i].Size
			Nexus.Popul[i].CanAxis[j].MaxVal = float64((1<<ExpAxisVal) + Rnd.Int31n((1<<ExpAxisVal) - (1<<(ExpAxisVal-1))) + gap)
			Nexus.Popul[i].CanAxis[j].MinVal = float64(0 - ((1<<ExpAxisVal) + Rnd.Int31n((1<<ExpAxisVal) - (1<<(ExpAxisVal-1))) + gap))
			if Rnd.Intn(2) == 0 {
				Nexus.Popul[i].ValAxis[j].Confort = float64(Rnd.Int31n(toint32(Nexus.Popul[i].CanAxis[j].MaxVal - (Nexus.Popul[i].CanAxis[j].MaxVal / (1<<ExpConfortBorder)))))
			} else {
				Nexus.Popul[i].ValAxis[j].Confort = float64(0 - Rnd.Int31n(toint32(math.Abs(float64(Nexus.Popul[i].CanAxis[j].MinVal + (Nexus.Popul[i].CanAxis[j].MinVal / (1<<ExpConfortBorder)))))))
			}
			Nexus.Popul[i].ValAxis[j].Radix = math.Trunc((MaxWill + (MaxWill / 2)) - Nexus.Popul[i].Will)
			Nexus.Popul[i].ValAxis[j].TopConf = Nexus.Popul[i].ValAxis[j].Confort + Nexus.Popul[i].ValAxis[j].Radix
			Nexus.Popul[i].ValAxis[j].BotConf = Nexus.Popul[i].ValAxis[j].Confort - Nexus.Popul[i].ValAxis[j].Radix
			if Nexus.Popul[i].ValAxis[j].Confort < 0 {
				Nexus.Popul[i].ValAxis[j].ActVal = float64(Rnd.Int31n(toint32(1<<(ExpAxisVal-4))))
			} else {
				Nexus.Popul[i].ValAxis[j].ActVal = float64(0 - Rnd.Int31n(toint32(1<<(ExpAxisVal-4))))
			}
			Nexus.Popul[i].ValAxis[j].Desire = int64(Rnd.Int31n(1<<ExpDesire) + 1)
			Nexus.Popul[i].ValAxis[j].Impetus = 0
			Nexus.Popul[i].ValAxis[j].weight = Nexus.Popul[i].Will * float64(Rnd.Intn(InterestBound/5)) / float64(InterestBound)
			Nexus.Popul[i].CalAxis[j].When = -1
			Nexus.Popul[i].CalAxis[j].Froz = math.Max(math.Trunc(Nexus.Popul[i].CanAxis[j].Polar * Nexus.Popul[i].CanAxis[j].Speed * Nexus.FysDim / math.Sqrt(Nexus.Popul[i].Will)),
									 				  Nexus.FysDim)
			Nexus.Popul[i].CalAxis[j].NewMax = Nexus.Popul[i].CanAxis[j].MaxVal - (Nexus.Popul[i].Will * Nexus.Popul[i].Size)
			Nexus.Popul[i].CalAxis[j].NewMin = Nexus.Popul[i].CanAxis[j].MinVal + (Nexus.Popul[i].Will * Nexus.Popul[i].Size)
		}
	}
}

func calcFEje() {
	for i := 0; i < NPop; i++ {
		for j := 0; j < NAxisPop; j++ {
			//
			force := Nexus.Popul[i].CanAxis[j].Polar *
					 Nexus.Popul[i].CanAxis[j].Speed *
					 Nexus.Popul[i].CanAxis[j].Push
			//
			Nexus.Popul[i].ValAxis[j].ActVal = Nexus.Popul[i].ValAxis[j].ActVal + force
			Nexus.Popul[i].CalAxis[j].Feje = force
		}
	}
}

func calcFEndo(pulse int) {
	for i := 0; i < NPop; i++ {
		for j := 0; j < NAxisPop; j++ {
			delta := math.Abs(Nexus.Popul[i].ValAxis[j].ActVal - Nexus.Popul[i].ValAxis[j].Confort)
			if delta <= math.MaxInt32 && delta >= math.MinInt32 {
				if delta > Nexus.Popul[i].ValAxis[j].Radix {

					newImpet := Nexus.Popul[i].ValAxis[j].Impetus
					
					var tideDesire int
					if Nexus.Popul[i].ValAxis[j].Desire == 1 {
						Nexus.Popul[i].ValAxis[j].Desire = int64(Rnd.Int31n(1<<ExpDesire) + 1)
					}
					//var preWill float64
					if Nexus.Popul[i].ValAxis[j].Desire % 2 == 0 {
						tideDesire = 1
						newImpet = math.Max(1, 2*newImpet)
						Nexus.Popul[i].ValAxis[j].Desire = Nexus.Popul[i].ValAxis[j].Desire / 2
						Nexus.Popul[i].CalAxis[j].Come = true
						//preWill = Nexus.Popul[i].ValAxis[j].PreWill * (math.Sqrt(delta) - 1)
					} else {
						tideDesire = -1
						newImpet = 0.25
						Nexus.Popul[i].ValAxis[j].Desire = 3 * Nexus.Popul[i].ValAxis[j].Desire + 1
						Nexus.Popul[i].CalAxis[j].Come = false
						//preWill = Nexus.Popul[i].ValAxis[j].PreWill * (math.Cbrt(delta) - 1)
					}
					
					var populDesire int
					if Nexus.Popul[i].ValAxis[j].ActVal > Nexus.Popul[i].ValAxis[j].TopConf {
						populDesire = -1
					} else if Nexus.Popul[i].ValAxis[j].ActVal < Nexus.Popul[i].ValAxis[j].BotConf {
						populDesire = 1
					} else {
						populDesire = 0
					}
					
					combDesire := tideDesire * populDesire
					
					feje := Nexus.Popul[i].CalAxis[j].Feje

//					log2 := math.Trunc(math.Log2(delta))
//					var rate int32
//					if log2 > 3 {
//						rate = 1<<uint((log2 - 3))
//					} else if log2 > 0 {
//						rate = 1<<uint(log2)
//					} else {
//						rate = 1
//					}
//					
//					delrat := math.Cbrt(delta) / float64(rate)
//					Nexus.Popul[i].CalAxis[j].DelRat = delrat


					weight := Nexus.Popul[i].ValAxis[j].weight
					
					fwill := Nexus.Popul[i].ValAxis[j].PreWill / math.Cbrt(delta)
					
					froz := Nexus.Popul[i].CalAxis[j].Froz
					
					auxVal := weight * float64(combDesire) * newImpet * (feje + fwill + froz)
					actVal := Nexus.Popul[i].ValAxis[j].ActVal

					if math.Abs(auxVal) > delta {
						for (math.Abs(auxVal) <= delta) && (math.Abs(auxVal) > 1) {
							auxVal = math.Max(auxVal * 0.1, 1)
						}
					}

					newVal := actVal + auxVal
					
					if newVal > Nexus.Popul[i].CanAxis[j].MaxVal {
						newVal = Nexus.Popul[i].CalAxis[j].NewMax
					} else if newVal < Nexus.Popul[i].CanAxis[j].MinVal {
						newVal = Nexus.Popul[i].CalAxis[j].NewMin
//					} else if math.Abs(newVal) > delta {
//						newVal = newVal / delta
					}

					Nexus.Popul[i].ValAxis[j].ActVal = newVal
					Nexus.Popul[i].CalAxis[j].weight = weight
					Nexus.Popul[i].ValAxis[j].Impetus = newImpet
					Nexus.Popul[i].CalAxis[j].Delta = delta
					Nexus.Popul[i].CalAxis[j].Fwill = fwill
					Nexus.Popul[i].CalAxis[j].Froz = froz
					Nexus.Popul[i].CalAxis[j].Desire = combDesire
//					Nexus.Popul[i].CalAxis[j].Rate = rate
				} else {
					Nexus.Popul[i].CalAxis[j].Delta = delta
					if Nexus.Popul[i].CalAxis[j].When == -1 {
						Nexus.Popul[i].CalAxis[j].When = int32(pulse)
					}
				}
			//
			} else {
					Nexus.Popul[i].CalAxis[j].Delta = -1
			}
		}
	}
}

func calcFExo() {
	for i := 0; i < NPop; i++ {
		for j := 0; j < NAxisPop; j++ {
		}
	}
}

func printWorld(pulse int, axis bool) {
	if show {
		if err := tb.Clear(tb.ColorDefault, tb.ColorDefault); err != nil {
			fmt.Println("Error")
		}
		for i := 0; i < NPop; i++ {
			print_tb(0, i * InterLin, tb.ColorWhite, tb.ColorBlack, "Del/Conv")
			print_tb(0, i * InterLin +  1, tb.ColorWhite, tb.ColorBlack, "When")
			print_tb(0, i * InterLin +  2, tb.ColorWhite, tb.ColorBlack, "Feje")
			print_tb(0, i * InterLin +  3, tb.ColorWhite, tb.ColorBlack, "Fwill")
			print_tb(0, i * InterLin +  4, tb.ColorWhite, tb.ColorBlack, "Froz")
			print_tb(0, i * InterLin +  5, tb.ColorWhite, tb.ColorBlack, "Impet")
			print_tb(0, i * InterLin +  6, tb.ColorWhite, tb.ColorBlack, "Weight")
			print_tb(0, i * InterLin +  7, tb.ColorWhite, tb.ColorBlack, "Confort")
			print_tb(0, i * InterLin +  8, tb.ColorWhite, tb.ColorBlack, "Radix")
			print_tb(0, i * InterLin +  9, tb.ColorWhite, tb.ColorBlack, "ActVal")
			print_tb(0, i * InterLin + 10, tb.ColorWhite, tb.ColorBlack, "comDes")
		}
		for i := 0; i < NPop; i++ {
			for j := 0; j < NAxisPop; j++ {
				var delta string
				var when string
				if Nexus.Popul[i].CalAxis[j].Delta < 0 {
					delta = "+++++++++"
				} else {
					delta = fmt.Sprintf("%9X", toint32(Nexus.Popul[i].CalAxis[j].Delta))
				}
				var yellow bool
				if Nexus.Popul[i].CalAxis[j].When < 0 {
					when = "*********"
					yellow = false
				} else {
					when = fmt.Sprintf("%9d", Nexus.Popul[i].CalAxis[j].When)
					yellow = true
				}
				feje := fmt.Sprintf("%9X", toint32(Nexus.Popul[i].CalAxis[j].Feje))
				fwill := fmt.Sprintf("%9X", toint32(Nexus.Popul[i].CalAxis[j].Fwill))
				froz := fmt.Sprintf("%9X", toint32(Nexus.Popul[i].CalAxis[j].Froz))
				impet := fmt.Sprintf("%9.6f", Nexus.Popul[i].ValAxis[j].Impetus)
				weight := fmt.Sprintf("%9.6f", Nexus.Popul[i].ValAxis[j].weight)
				conf := fmt.Sprintf("%9X", toint32(Nexus.Popul[i].ValAxis[j].Confort))
				radix := fmt.Sprintf("%9X", toint32(Nexus.Popul[i].ValAxis[j].Radix))
				actval := fmt.Sprintf("%9X", toint32(Nexus.Popul[i].ValAxis[j].ActVal))
				cdesire := fmt.Sprintf("%9d", Nexus.Popul[i].CalAxis[j].Desire)
				
				var fgColor tb.Attribute
				var dtColor tb.Attribute
				if yellow {
					fgColor = tb.ColorYellow				
					dtColor = tb.ColorYellow
				} else {
					fgColor = tb.ColorCyan
					switch {
					case Nexus.Popul[i].ValAxis[j].ActVal >= Nexus.Popul[i].CanAxis[j].MaxVal:
						dtColor = tb.ColorBlue
					case Nexus.Popul[i].ValAxis[j].ActVal <= Nexus.Popul[i].CanAxis[j].MinVal:
						dtColor = tb.ColorMagenta
					case Nexus.Popul[i].CalAxis[j].Come:
						dtColor = tb.ColorGreen
					default:
						dtColor = tb.ColorRed
					}
				}
				print_tb(coord[i][j][0], coord[i][j][1] +  0, dtColor, tb.ColorBlack, delta)
				print_tb(coord[i][j][0], coord[i][j][1] +  1, fgColor, tb.ColorBlack, when)
				print_tb(coord[i][j][0], coord[i][j][1] +  2, fgColor, tb.ColorBlack, feje)
				print_tb(coord[i][j][0], coord[i][j][1] +  3, fgColor, tb.ColorBlack, fwill)
				print_tb(coord[i][j][0], coord[i][j][1] +  4, fgColor, tb.ColorBlack, froz)
				print_tb(coord[i][j][0], coord[i][j][1] +  5, fgColor, tb.ColorBlack, impet)
				print_tb(coord[i][j][0], coord[i][j][1] +  6, fgColor, tb.ColorBlack, weight)
				print_tb(coord[i][j][0], coord[i][j][1] +  7, fgColor, tb.ColorBlack, conf)
				print_tb(coord[i][j][0], coord[i][j][1] +  8, fgColor, tb.ColorBlack, radix)
				print_tb(coord[i][j][0], coord[i][j][1] +  9, fgColor, tb.ColorBlack, actval)
				print_tb(coord[i][j][0], coord[i][j][1] + 10, fgColor, tb.ColorBlack, cdesire)
			}
		}
		print_tb(wx-30, wy-1, tb.ColorWhite, tb.ColorBlack, fmt.Sprintf("Sleep: %20d", *FlagSleep))
		print_tb(wx-30, wy-2, tb.ColorWhite, tb.ColorBlack, fmt.Sprintf("Steep: %20d", *FlagSteep))
		print_tb(wx-30, wy-3, tb.ColorWhite, tb.ColorBlack, fmt.Sprintf("Pulse: %20d", pulse))
		time.Sleep(time.Duration(*FlagSleep) * time.Millisecond)
	} else {
		print_tb(wx-30, wy-1, tb.ColorWhite, tb.ColorBlack, fmt.Sprintf("Sleep: %20d", *FlagSleep))
		print_tb(wx-30, wy-2, tb.ColorWhite, tb.ColorBlack, fmt.Sprintf("Steep: %20d", *FlagSteep))
		print_tb(wx-30, wy-3, tb.ColorWhite, tb.ColorBlack, fmt.Sprintf("Pulse: %20d", pulse))
	}

	if err := tb.Flush(); err != nil {
		fmt.Println("Error")
	}
}

func print_tb(x, y int, fg, bg tb.Attribute, msg string) {
	for _, c := range msg {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
}

func toint32(x float64) int32 {
	if x <= math.MaxInt32 && x >= math.MinInt32 {
		return int32(math.Trunc(x))
	} else if x <= math.MinInt32 {
		print_tb(wx-4, 1, tb.ColorRed, tb.ColorBlack, "Min")
		return math.MinInt32
	}
	print_tb(wx-4, 1, tb.ColorRed, tb.ColorBlack, "Max")
	return math.MaxInt32
}