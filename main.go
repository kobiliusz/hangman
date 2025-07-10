package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"strings"
)

//go:embed assets/*
var assets embed.FS

const width = 350
const height = 700

var typeface = loadFont()

type Game struct {
	chances int
	guessed []rune
	answer  string
	word    string
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	col := color.White

	text.Draw(screen, g.answer, typeface, 10, 620, col)
}

func (g *Game) Layout(w, h int) (int, int) {
	return width, height
}

func (g *Game) Init() {
	g.chances = 6
	g.guessed = make([]rune, 0)
	g.word = getWord()
	g.answer = strings.Repeat("_", len(g.word))
}

func getWord() string {
	//TODO
	return "TUSIDELEC"
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Hangman")
	game := &Game{}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func loadFont() font.Face {
	bytes, err := assets.ReadFile("assets/font.ttf")
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(bytes)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     144,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return face
}
