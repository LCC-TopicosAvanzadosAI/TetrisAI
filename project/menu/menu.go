package menu

import (
	"fmt"
	"github.com/faiface/pixel/text"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"

	hp "github.com/TetrisAI/project/helper"
)

const (
	BoardRows = 22
	BoardCols = 10
)

func DisplayMenu(win *pixelgl.Window, windowWidth, windowHeight float64) string {
	face, err := hp.LoadTTF("./../../resources/saarland.ttf", 52) //Loading font and size-font
	if err != nil {
		panic(err)
	}

	Atlas := text.NewAtlas(face, text.ASCII) //Atlas necessary for the font
	//here, I put the coordinates where the
	//texts starts to write

	for !win.Closed() {
		basicTxt := text.New(pixel.V(windowWidth/2, 200), Atlas)

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
		rectLearn := pixel.Rect(basicTxt.Bounds())

		txt = "Cerrar"
		basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2
		basicTxt.Color = colornames.Royalblue
		fmt.Fprintln(basicTxt, txt)
		rectQuit := pixel.Rect(basicTxt.Bounds())

		win.Clear(colornames.Black)
		basicTxt.Draw(win, pixel.IM)
		win.Update()
		if (rectJugar.Contains(win.MousePosition()) && win.JustPressed(pixelgl.MouseButtonLeft)) || win.Pressed(pixelgl.KeyEnter) {
			fmt.Println("returning play")
			return "Play"
		} else if (rectLearn.Contains(win.MousePosition()) && win.JustPressed(pixelgl.MouseButtonLeft)) || win.Pressed(pixelgl.KeyEscape) {
			return "Learn"
		} else if (rectQuit.Contains(win.MousePosition()) && win.JustPressed(pixelgl.MouseButtonLeft)) || win.Pressed(pixelgl.KeyEscape) {
			return "Quit"
		}
	}
	return ""
}
