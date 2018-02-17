package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	_ "image/png"

	//gb "github.com/TetrisAI/project/gameboard"
	hp "github.com/TetrisAI/project/helper"
	mn "github.com/TetrisAI/project/menu"
)

func main() {
	pixelgl.Run(run) //para correr la wea gráfica
}

func run() {
	//Window creation
	windowWidth := 765.0
	windowHeight := 450.0
	cfg := pixelgl.WindowConfig{
		Title:  "Tetris Menu",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	menu := mn.NewMenu()
	//fmt.Println(menu)
	menu.DisplayMenu(win)
	win.Clear(colornames.Black)

	face, err := hp.LoadTTF("./../../resources/saarland.ttf", 52) //Loading font and size-font
	if err != nil {
		panic(err)
	}

	Atlas := text.NewAtlas(face, text.ASCII)                 //Atlas necessary for the font
	basicTxt := text.New(pixel.V(windowWidth/2, 200), Atlas) //here, I put the coordinates where the texts starts to write

	basicTxt.LineHeight = Atlas.LineHeight() * 1.5 // line spacing between strings

	txt := "Jugar"
	basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2 //centralize text
	basicTxt.Color = colornames.Aqua                 //text color
	fmt.Fprintln(basicTxt, txt)                      //put the text in the window
	rectJugar := pixel.Rect(basicTxt.Bounds())       //creation a rectangle around the text

	txt = "Aprender"
	basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2
	basicTxt.Color = colornames.Green
	fmt.Fprintln(basicTxt, txt)
	//rectAprender := pixel.Rect(basicTxt.Bounds())

	txt = "Cerrar"
	basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2
	basicTxt.Color = colornames.Blue
	fmt.Fprintln(basicTxt, txt)
	//rectCerrar := pixel.Rect(basicTxt.Bounds())

	//fmt.Println("cords jugar: ", rectJugar.Min.X, " ", rectJugar.Min.Y, " ", rectJugar.Max.X, " ", rectJugar.Max.Y)
	for !win.Closed() {
		win.Clear(colornames.Black)
		basicTxt.Draw(win, pixel.IM)
		win.Update()

		menu.DisplayMenu(win)
		//fmt.Println(win.MousePosition().X, " ", win.MousePosition().Y)

		if (rectJugar.Contains(win.MousePosition()) && win.JustPressed(pixelgl.MouseButtonLeft)) || win.Pressed(pixelgl.KeyEnter) {
			menu.Jugar(win)
		}
	}
}
