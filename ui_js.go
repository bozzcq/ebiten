// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build js

package ebiten

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/hajimehoshi/ebiten/internal/audio"
	"github.com/hajimehoshi/ebiten/internal/graphics/internal/opengl"
	"github.com/hajimehoshi/ebiten/internal/ui"
	"strconv"
)

var canvas js.Object
var context *opengl.Context

type userInterface struct{}

var currentUI = &userInterface{}

// TODO: This returns true even when the browser is not active.
// The current behavior causes sound noise...
func shown() bool {
	return !js.Global.Get("document").Get("hidden").Bool()
}

func useGLContext(f func(*opengl.Context)) {
	f(context)
}

func vsync() {
	ch := make(chan struct{})
	// TODO: In iOS8, this is called at every 1/30[sec] frame.
	// Can we use DOMHighResTimeStamp?
	js.Global.Get("window").Call("requestAnimationFrame", func() {
		close(ch)
	})
	<-ch
}

func (*userInterface) doEvents() error {
	vsync()
	for !shown() {
		vsync()
	}
	ui.CurrentInput().UpdateGamepads()
	return nil
}

func (*userInterface) terminate() {
	// Do nothing.
}

func (*userInterface) isClosed() bool {
	return false
}

func (*userInterface) swapBuffers() {
	// Do nothing.
}

func init() {
	if js.Global.Get("require") != js.Undefined {
		// Use headless-gl for testing.
		nodeGl := js.Global.Call("require", "gl")
		webglContext := nodeGl.Call("createContext", 16, 16)
		context = opengl.NewContext(&webgl.Context{Object: webglContext})
		return
	}

	doc := js.Global.Get("document")
	window := js.Global.Get("window")
	if doc.Get("body") == nil {
		ch := make(chan struct{})
		window.Call("addEventListener", "load", func() {
			close(ch)
		})
		<-ch
	}

	canvas = doc.Call("createElement", "canvas")
	canvas.Set("width", 16)
	canvas.Set("height", 16)
	doc.Get("body").Call("appendChild", canvas)

	htmlStyle := doc.Get("documentElement").Get("style")
	htmlStyle.Set("height", "100%")
	htmlStyle.Set("margin", "0")
	htmlStyle.Set("padding", "0")

	bodyStyle := doc.Get("body").Get("style")
	bodyStyle.Set("backgroundColor", "#000")
	bodyStyle.Set("position", "relative")
	bodyStyle.Set("height", "100%")
	bodyStyle.Set("margin", "0")
	bodyStyle.Set("padding", "0")
	doc.Get("body").Call("addEventListener", "click", func() {
		canvas.Call("focus")
	})

	canvasStyle := canvas.Get("style")
	canvasStyle.Set("position", "absolute")

	webglContext, err := webgl.NewContext(canvas, &webgl.ContextAttributes{
		Alpha:              true,
		PremultipliedAlpha: true,
	})
	if err != nil {
		panic(err)
	}
	context = opengl.NewContext(webglContext)

	// Make the canvas focusable.
	canvas.Call("setAttribute", "tabindex", 1)
	canvas.Get("style").Set("outline", "none")

	// Keyboard
	canvas.Call("addEventListener", "keydown", func(e js.Object) {
		e.Call("preventDefault")
		code := e.Get("keyCode").Int()
		ui.CurrentInput().KeyDown(code)
	})
	canvas.Call("addEventListener", "keyup", func(e js.Object) {
		e.Call("preventDefault")
		code := e.Get("keyCode").Int()
		ui.CurrentInput().KeyUp(code)
	})

	// Mouse
	canvas.Call("addEventListener", "mousedown", func(e js.Object) {
		e.Call("preventDefault")
		button := e.Get("button").Int()
		ui.CurrentInput().MouseDown(button)
		setMouseCursorFromEvent(e)
	})
	canvas.Call("addEventListener", "mouseup", func(e js.Object) {
		e.Call("preventDefault")
		button := e.Get("button").Int()
		ui.CurrentInput().MouseUp(button)
		setMouseCursorFromEvent(e)
	})
	canvas.Call("addEventListener", "mousemove", func(e js.Object) {
		e.Call("preventDefault")
		setMouseCursorFromEvent(e)
	})
	canvas.Call("addEventListener", "contextmenu", func(e js.Object) {
		e.Call("preventDefault")
	})

	// Touch (emulating mouse events)
	// TODO: Need to create indimendent touch functions?
	canvas.Call("addEventListener", "touchstart", func(e js.Object) {
		e.Call("preventDefault")
		ui.CurrentInput().MouseDown(0)
		touches := e.Get("changedTouches")
		touch := touches.Index(0)
		setMouseCursorFromEvent(touch)
	})
	canvas.Call("addEventListener", "touchend", func(e js.Object) {
		e.Call("preventDefault")
		ui.CurrentInput().MouseUp(0)
		touches := e.Get("changedTouches")
		touch := touches.Index(0)
		setMouseCursorFromEvent(touch)
	})
	canvas.Call("addEventListener", "touchmove", func(e js.Object) {
		e.Call("preventDefault")
		touches := e.Get("changedTouches")
		touch := touches.Index(0)
		setMouseCursorFromEvent(touch)
	})

	// Gamepad
	window.Call("addEventListener", "gamepadconnected", func(e js.Object) {
		// Do nothing.
	})

	audio.Init()
}

func setMouseCursorFromEvent(e js.Object) {
	scale := canvas.Get("dataset").Get("ebitenScale").Int() // TODO: Float?
	rect := canvas.Call("getBoundingClientRect")
	x, y := e.Get("clientX").Int(), e.Get("clientY").Int()
	x -= rect.Get("left").Int()
	y -= rect.Get("top").Int()
	ui.CurrentInput().SetMouseCursor(x/scale, y/scale)
}

func devicePixelRatio() int {
	// TODO: What if ratio is not an integer but a float?
	ratio := js.Global.Get("window").Get("devicePixelRatio").Int()
	if ratio == 0 {
		ratio = 1
	}
	return ratio
}

func (*userInterface) start(width, height, scale int, title string) (actualScale int, err error) {
	doc := js.Global.Get("document")
	doc.Set("title", title)
	actualScale = scale * devicePixelRatio()
	canvas.Set("width", width*actualScale)
	canvas.Set("height", height*actualScale)
	canvas.Get("dataset").Set("ebitenScale", scale)
	canvasStyle := canvas.Get("style")

	cssWidth := width * scale
	cssHeight := height * scale
	canvasStyle.Set("width", strconv.Itoa(cssWidth)+"px")
	canvasStyle.Set("height", strconv.Itoa(cssHeight)+"px")
	// CSS calc requires space chars.
	canvasStyle.Set("left", "calc(50% - "+strconv.Itoa(cssWidth/2)+"px)")
	canvasStyle.Set("top", "calc(50% - "+strconv.Itoa(cssHeight/2)+"px)")

	canvas.Call("focus")

	audio.Start()

	return actualScale, nil
}