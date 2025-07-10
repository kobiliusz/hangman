package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"slices"
	"strings"
	"unicode"
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
	message string
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	col := color.White

	text.Draw(screen, g.answer, typeface, 10, 620, col)
	text.Draw(screen, g.message, typeface, 10, 650, col)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return width, height
}

func (g *Game) Init() {
	g.chances = 6
	g.guessed = make([]rune, 0)
	g.word = getWord()
	g.answer = strings.Repeat("_", len(g.word))
	g.message = "Start of game."
}

func (g *Game) Guess(c rune) {
	c = unicode.ToUpper(c)
	if slices.Contains(g.guessed, c) {
		g.message = fmt.Sprintf("Already guessed: '%s'", string(c))
	}
	g.guessed = append(g.guessed, c)
	indexes := runeIndexesOfRune(g.word, c)
	if len(indexes) == 0 {
		g.chances--
		if g.chances == 0 {
			g.answer = g.word
			g.message = "Press space\nfor new game"
			return
		}
		g.message = fmt.Sprintf("Wrong: '%s'", string(c))
	} else {
		for _, index := range indexes {
			g.answer = replaceNthRune(g.answer, index, c)
		}
		g.message = fmt.Sprintf("Correct: '%s'", string(c))
	}
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
		Size:    14,
		DPI:     144,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return face
}

func runeIndexesOfRune(s string, target rune) []int {
	var indexes []int
	pos := 0
	for _, r := range s {
		if r == target {
			indexes = append(indexes, pos)
		}
		pos++
	}
	return indexes
}

func replaceNthRune(s string, n int, newRune rune) string {
	runes := []rune(s)
	if n < 0 || n >= len(runes) {
		return s // poza zakresem
	}
	runes[n] = newRune
	return string(runes)
}
