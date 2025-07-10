package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"strings"
	"unicode"
)

//go:embed assets/*
var assets embed.FS

//go:embed assets/dict.txt
var wordsFile string

const width = 350
const height = 750

var typeface = loadFont()

type Game struct {
	chances int
	guessed []rune
	answer  string
	word    string
	message string
	win     bool
}

func (g *Game) Update() error {
	if g.chances > 0 && !g.win {
		chars := ebiten.AppendInputChars(nil)
		if len(chars) > 0 {
			g.Guess(chars[0])
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Init()
		}
	}
	return nil
}

//goland:noinspection GoDeprecation
func (g *Game) Draw(screen *ebiten.Image) {
	col := color.White

	text.Draw(screen, g.answer, typeface, 10, 650, col)
	text.Draw(screen, g.message, typeface, 10, 680, col)

	ebitenutil.DrawRect(screen, 50, 600, 250, 10, col)
	ebitenutil.DrawRect(screen, 100, 200, 20, 400, col)
	ebitenutil.DrawRect(screen, 100, 200, 150, 20, col)
	ebitenutil.DrawRect(screen, 240, 220, 10, 40, col)

	x := 245.0
	y := 260.0

	if g.chances <= 5 {
		// Głowa (koło)
		ebitenutil.DrawCircle(screen, x, y+20, 20, col)
		ebitenutil.DrawCircle(screen, x, y+20, 18, color.Black)
	}
	if g.chances <= 4 {
		// Tułów (linia w dół)
		ebitenutil.DrawLine(screen, x, y+40, x, y+100, col)
	}
	if g.chances <= 3 {
		// Lewa ręka
		ebitenutil.DrawLine(screen, x, y+60, x-30, y+80, col)
	}
	if g.chances <= 2 {
		// Prawa ręka
		ebitenutil.DrawLine(screen, x, y+60, x+30, y+80, col)
	}
	if g.chances <= 1 {
		// Lewa noga
		ebitenutil.DrawLine(screen, x, y+100, x-25, y+140, col)
	}
	if g.chances == 0 {
		// Prawa noga
		ebitenutil.DrawLine(screen, x, y+100, x+25, y+140, col)
	}
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
	g.win = false
}

func (g *Game) Guess(c rune) {
	c = unicode.ToUpper(c)
	if slices.Contains(g.guessed, c) {
		g.message = fmt.Sprintf("Already guessed:'%s'", string(c))
		return
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
		if len(runeIndexesOfRune(g.answer, '_')) > 0 {
			g.message = fmt.Sprintf("Correct: '%s'", string(c))
		} else {
			g.win = true
			g.message = "Win! Press space\nfor new game"
		}
	}
}

func getWord() string {
	words := getWordList()
	index := rand.Intn(len(words))
	return strings.ToUpper(words[index])
}

func getWordList() []string {
	words := make([]string, 0)
	for _, line := range strings.Split(wordsFile, "\n") {
		words = append(words, strings.TrimSpace(line))
	}
	return words
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
