package main

import tb "github.com/nsf/termbox-go"
import "fmt"

func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()

	x, y := tb.Size()

	err = tb.Clear(tb.ColorDefault, tb.ColorDefault)
	if err != nil {
		fmt.Println("Error")
	}
	err = tb.Flush()
	if err != nil {
		fmt.Println("Error")
	}
	tb.SetCursor(0, 0)
	fmt.Printf("**x: %d\n**y: %d\n", x, y)

	fmt.Println("Done")
}