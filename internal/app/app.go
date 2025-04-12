package app

import (
	"runtime"
	"shadxel/internal/luaengine"
	"shadxel/internal/render"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	window    *sdl.Window
	renderer  *render.Renderer
	rotation  float32
	mouseHeld bool
	lastX     int32
	engine    *luaengine.LuaEngine
}

// Call this first — SDL requires main thread
func init() {
	runtime.LockOSThread()
}

func NewApp() (*App, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	window, err := sdl.CreateWindow("Voxel Viewer", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1600, 1200, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return nil, err
	}

	_, err = window.GLCreateContext()
	if err != nil {
		return nil, err
	}

	if err := gl.Init(); err != nil {
		return nil, err
	}

	lua, err := luaengine.NewLuaEngine("script.lua")
	if err != nil {
		return nil, err
	}

	renderer, err := render.NewRenderer()
	if err != nil {
		return nil, err
	}

	return &App{
		window:   window,
		renderer: renderer,
		engine:   lua,
	}, nil
}

func (a *App) Run() {
	defer sdl.Quit()
	defer a.window.Destroy()

	start := time.Now()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_LEFT {
					a.mouseHeld = e.State == sdl.PRESSED
					a.lastX = e.X
				}
			case *sdl.MouseMotionEvent:
				if a.mouseHeld {
					dx := e.X - a.lastX
					a.rotation += float32(dx) * 0.01
					a.lastX = e.X
				}
			}
		}

		elapsed := int(time.Since(start).Seconds())
		grid := a.engine.GenerateGrid(50, 50, elapsed)
		a.renderer.DrawGrid(grid, a.rotation)

		a.window.GLSwap()
		time.Sleep(16 * time.Millisecond) // ~60 FPS
	}
}
