package main

import (
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
)

type CLUI struct {
	*sync.Mutex
	g        *gocui.Gui
	callback func([]byte)
}

func NewCLUI() (*CLUI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	ret := &CLUI{&sync.Mutex{}, g, nil}
	g.Cursor = true
	g.SetManagerFunc(ret.layout)

	if err := ret.setup(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (cl *CLUI) run() error {
	defer cl.g.Close()
	if err := cl.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (cl *CLUI) flush() {
	cl.g.Update(func(g *gocui.Gui) error { return nil })
}

func (cl *CLUI) print(text string, sview string) {
	cl.Lock()
	if view, err := cl.g.View(sview); err == nil {
		view.Write([]byte(text))
	}
	cl.flush()
	cl.Unlock()
}

func (cl *CLUI) printPrompt(text []byte, sview string) {
	cl.Lock()
	if view, err := cl.g.View(sview); err == nil {
		view.Write(append([]byte("> "), text...))
	}
	cl.flush()
	cl.Unlock()
}

func (cl *CLUI) setup() error {
	if err := cl.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := cl.g.SetKeybinding("send", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			text := strings.TrimSpace(v.Buffer())
			content := []byte(text)
			if cl.callback != nil {
				cl.callback(content)
			}
			// cl.printPrompt(content, "input")
			v.Clear()
			v.SetCursor(0, 0)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func (cl *CLUI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("send", 0, maxY-5, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("send"); err != nil {
			return err
		}
		v.Title = "Send"
		v.Editable = true
		v.Wrap = true
	}

	if v, err := g.SetView("input", 0, 0, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Overwrite = true
		v.Wrap = true
		v.Autoscroll = true
	}

	return nil
}
