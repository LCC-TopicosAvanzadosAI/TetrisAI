package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	_ "image/png"
	"io/ioutil"
	"os"

	//gb "github.com/TetrisAI/project/gameboard"
	//hp "github.com/TetrisAI/project/helper"
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
	fmt.Println(menu)

	menu.DisplayMenu(win)
	win.Clear(colornames.Black)

	face, err := loadTTF("saarland.ttf", 52)
	if err != nil {
		panic(err)
	}

	Atlas := text.NewAtlas(face, text.ASCII)
	basicTxt := text.New(pixel.V(windowWidth/2, 200), Atlas)

	basicTxt.LineHeight = Atlas.LineHeight() * 1.5

	txt := "Jugar"
	basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2
	basicTxt.Color = colornames.Aqua
	fmt.Fprintln(basicTxt, txt)
	rectJugar := pixel.Rect(basicTxt.Bounds())
	txt = "Aprender"
	basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2
	basicTxt.Color = colornames.Green
	fmt.Fprintln(basicTxt, txt)
	//rectAprender := basicTxt.R()
	txt = "Cerrar"
	basicTxt.Dot.X -= basicTxt.BoundsOf(txt).W() / 2
	basicTxt.Color = colornames.Blue
	fmt.Fprintln(basicTxt, txt)
	//rectCerrar := basicTxt.R()

	//fmt.Println("cords jugar: ", rectJugar.Min.X, " ", rectJugar.Min.Y, " ", rectJugar.Max.X, " ", rectJugar.Max.Y)
	for !win.Closed() {
		win.Clear(colornames.Black)
		basicTxt.Draw(win, pixel.IM)
		//basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
		win.Update()

		menu.DisplayMenu(win)
		//fmt.Println(win.MousePosition().X, " ", win.MousePosition().Y)

		if rectJugar.Contains(win.MousePosition()) && win.JustPressed(pixelgl.MouseButtonLeft) {
			menu.Jugar(win)
			//draw.Draw(rectJugar, pixel.IM)
			//draw.Draw(rectJugar)
		}
	}
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
