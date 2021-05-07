package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	cirleColor   = color.RGBA{3, 194, 252, 255}
	outerColor   = color.RGBA{252, 98, 3, 255}
	currentTicks = 0
	printer      = message.NewPrinter(language.English)
)

const (
	screenWidth       = 640
	screenHeight      = 640
	noisePerTick      = 8000
	pointsToCalculate = 2000000
	ticksBeforeStart  = 120
)

// Game is our main game object.
type Game struct {
	totalPoints    int32
	pointsInCircle int32
	estimatedPi    float64
	finished       bool
	img            *image.RGBA
}

// Update updates game state. Called every tick.
func (g *Game) Update() error {
	// We delay so the user just sees circle for a bit before we fill it in.
	if currentTicks < ticksBeforeStart {
		currentTicks++
		return nil
	}
	if g.finished {
		return nil
	}
	// Generate new noise
	g.GenerateNoise()
	return nil
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.img.Pix)
	delta := math.Abs(g.estimatedPi - math.Pi)
	ebitenutil.DebugPrint(screen, printer.Sprintf("TPS: %0.2f. \nPoints: %d. \nEstimated Pi: %f. \nActual Pi: %f. \nDelta: %f", ebiten.CurrentTPS(), g.totalPoints, g.estimatedPi, math.Pi, delta))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// GenerateNoise generates new noise in
func (g *Game) GenerateNoise() {
	i := 0
	for i < noisePerTick {
		// We're done.
		if g.totalPoints == pointsToCalculate {
			g.finished = true
			return
		}
		// Generate random x/y coordinates
		x := rand.Intn(screenWidth)
		y := rand.Intn(screenHeight)
		g.totalPoints++

		// In my case, x-center, y-center and radius are the same.
		if withinCircle(x, y, screenWidth/2, screenHeight/2, screenWidth/2) {
			g.img.Set(x, y, cirleColor)
			g.pointsInCircle++
		} else {
			// We are outside the circle
			g.img.Set(x, y, outerColor)
		}
		g.estimatedPi = 4 * float64(g.pointsInCircle) / float64(g.totalPoints)
		i++
	}
}

// withinCircle is a very simple Pythagoras implementation.
// See also: https://stackoverflow.com/questions/481144/equation-for-testing-if-a-point-is-inside-a-circle
func withinCircle(x, y, centerX, centerY, radius int) bool {
	return (x-centerX)*(x-centerX)+(y-centerY)*(y-centerY) < (radius * radius)
}

// DrawCircle draws a circle for visualization.
func (g *Game) DrawCircle() {
	offsetX := screenWidth / 2
	offsetY := screenHeight / 2
	radius := screenWidth / 2
	x, y, dx, dy := radius-1, 0, 1, 1
	err := dx - (radius * 2)
	// See https://stackoverflow.com/questions/51626905/drawing-circles-with-two-radius-in-golang
	for x > y {
		g.img.Set(offsetX+x, offsetY+y, cirleColor)
		g.img.Set(offsetX+y, offsetY+x, cirleColor)
		g.img.Set(offsetX-y, offsetY+x, cirleColor)
		g.img.Set(offsetX-x, offsetY+y, cirleColor)
		g.img.Set(offsetX-x, offsetY-y, cirleColor)
		g.img.Set(offsetX-y, offsetY-x, cirleColor)
		g.img.Set(offsetX+y, offsetY-x, cirleColor)
		g.img.Set(offsetX+x, offsetY-y, cirleColor)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (radius * 2)
		}
	}
}

func main() {
	// Initialize our game object.
	game := &Game{
		img:            image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		totalPoints:    0,
		pointsInCircle: 0,
		finished:       false,
	}
	// Set window size and tile
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Monte carlo simulation to estimate Pi.")
	// Draw initial circle
	game.DrawCircle()
	// Start game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
