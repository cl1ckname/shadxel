package app

import (
	"runtime"
	"shadxel/internal/camera"
	"shadxel/internal/luaengine"
	"shadxel/internal/render"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

var projection = mgl32.Perspective(mgl32.DegToRad(45.0), 1, 0.1, 100.0)

type App struct {
	window    *sdl.Window
	renderer  *render.Renderer
	camera    *camera.OrbitCamera
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
	gl.Enable(gl.DEPTH_TEST)

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
		camera:   camera.NewOrbitCamera(),
	}, nil
}

func (a *App) Run() {
	defer sdl.Quit()
	defer a.window.Destroy()

	ticker := time.NewTicker(time.Second / 2) // Redraw once per second
	defer ticker.Stop()
	running := true
	var frame int

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
					dx := e.XRel
					dy := e.YRel
					a.camera.Rotate(float32(dx), float32(dy))
				}
			}
		}
		select {
		case <-ticker.C:
			gl.Viewport(0, 0, 1600, 1200) // or update this dynamically on resize
			gl.ClearColor(0.9, 0.9, 0.9, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			grid := a.engine.GenerateGrid(50, frame)

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			view := a.camera.ViewMatrix()
			a.renderer.Draw(grid, view, projection)
			a.window.GLSwap()
			frame++
		default:
			// Avoid maxing CPU
			sdl.Delay(1)
		}
	}
}
