package app

import (
	"log"
	"runtime"
	"shadxel/internal/camera"
	"shadxel/internal/config"
	"shadxel/internal/luaengine"
	"shadxel/internal/render"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth = 1600.
	WindowHeigh = 1200.
	Aspect      = WindowWidth / WindowHeigh
)

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

func NewApp(c config.Config) (*App, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	window, err := sdl.CreateWindow("Shadxel", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WindowWidth, WindowHeigh, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
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

	lua, err := luaengine.NewLuaEngine(c.Script)
	if err != nil {
		return nil, err
	}

	renderer, err := render.NewRenderer(Aspect)
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

	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	running := true
	var frame int

	grid := a.engine.GenerateGrid(50, frame)
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
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_r {
					err := a.engine.Load()
					if err != nil {
						log.Println("Failed to reload script:", err)
					} else {
						log.Println("Script reloaded!")
					}
				}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					width := e.Data1
					height := e.Data2
					a.renderer.Resize(width, height)
					log.Printf("Window resized: %dx%d\n", width, height)

					// Optionally update your projection matrix here if it depends on aspect ratio
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
			grid = a.engine.GenerateGrid(50, frame)
			frame++
		default:
			// Avoid maxing CPU
			sdl.Delay(1)
		}
		view := a.camera.ViewMatrix()
		a.renderer.Draw(grid, view)
		a.window.GLSwap()
	}
}
