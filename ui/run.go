package ui

import (
	"log"
	"runtime"

	"../nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gordonklaus/portaudio"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "NES"
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}

func Run(path string) {
	// initialize audio
	portaudio.Initialize()
	defer portaudio.Terminate()

	audio := NewAudio()
	if err := audio.Start(); err != nil {
		log.Fatalln(err)
	}
	defer audio.Stop()

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width*scale, height*scale, title, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	// run
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}

	view := NewGameView(window, audio, console, path)
	view.Enter()
	timestamp := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		now := glfw.GetTime()
		dt := now - timestamp
		timestamp = now
		if view != nil {
			view.Update(now, dt)
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
	view.Exit()
}
