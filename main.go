package main

import (
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"go-platform/drawer"
	"go-platform/generator"
	"go-platform/sprite"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	gopherSpritePath  = "./assets/sprites/gopher.png"
	spriteCoordinates = "./assets/sprites/coordinates/gopher.csv"
)

type GoalDrawer interface {
	Draw(imd *imdraw.IMDraw)
	Update(dt float64, randomColor pixel.RGBA)
}

func run() {
	rand.Seed(time.Now().UnixNano())

	sheet, anims, err := sprite.Load(gopherSpritePath, spriteCoordinates, 12)
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "SuperGopher",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	physics := drawer.NewSpritePhysicsUpdater(-300, 64, 192, pixel.R(-6, -7, 6, 7), pixel.ZV, false)
	color := generator.NewRandomColor()
	gopherDrawer := drawer.NewSpriteDrawer(sheet, anims, 1.0/10, 0, 0, +1, anims["Front"][0])
	goal := drawer.NewGoalDrawer(pixel.V(70, 40), 10, 1.0/7, 0, [10]pixel.RGBA{})

	// hardcoded level for now
	platforms := []drawer.Platform{
		{Rect: pixel.R(-1000, -11, 100, -10), Color: color.Generate()},
		{Rect: pixel.R(150, -11, 300, -10), Color: color.Generate()},
		{Rect: pixel.R(270, -30, 400, -29), Color: color.Generate()},
		{Rect: pixel.R(-10, 1.5, 20, 2), Color: color.Generate()},
		{Rect: pixel.R(40, 10, 60, 12), Color: color.Generate()},
		{Rect: pixel.R(40, 10, 60, 12), Color: color.Generate()},
	}

	canvas := pixelgl.NewCanvas(pixel.R(-160/2, -120/2, 160/2, 120/2))
	imd := imdraw.New(sheet)
	imd.Precision = 32

	camPos := pixel.ZV

	last := time.Now()

	// Game loop
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// lerp the camera position towards the gopher
		camPos = pixel.Lerp(camPos, physics.Rect.Center(), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos.Scaled(-1))
		canvas.SetMatrix(cam)

		// slow motion with tab
		if win.Pressed(pixelgl.KeyTab) {
			dt /= 8
		}

		// restart the level on pressing enter
		if win.JustPressed(pixelgl.KeyEnter) {
			physics.Rect = physics.Rect.Moved(physics.Rect.Center().Scaled(-1))
			physics.Vel = pixel.ZV
		}

		// control the gopher with keys
		ctrl := pixel.ZV
		if win.Pressed(pixelgl.KeyLeft) {
			ctrl.X--
		}
		if win.Pressed(pixelgl.KeyRight) {
			ctrl.X++
		}
		if win.JustPressed(pixelgl.KeySpace) {
			ctrl.Y = 1
		}

		// update the physics and animation
		physics.Update(dt, ctrl, platforms)
		goal.Update(dt, color.Generate())
		gopherDrawer.Update(dt, physics)

		// draw the scene to the canvas using IMDraw
		canvas.Clear(colornames.Black)
		imd.Clear()

		for _, p := range platforms {
			p.Draw(imd)
		}

		goal.Draw(imd)
		gopherDrawer.Draw(imd, physics)
		imd.Draw(canvas)

		// stretch the canvas to the window
		win.Clear(colornames.White)
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
