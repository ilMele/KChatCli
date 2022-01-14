package main

import (
	"os"

	"github.com/gdamore/tcell"

	"kiwichat/client/conn"
	"kiwichat/client/display"
)

type Idisplay interface {
	Draw(s tcell.Screen)
	Write(s tcell.Screen, r rune)
	Enter(s tcell.Screen) (string, bool)
	Delete(s tcell.Screen)
}

func main() {

	state := "login"

	msg := make(chan string)

	ser := conn.Init()

	defer ser.Close()

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreenYellow)

	Dlogin := display.Login{}
	Dlogin.Default(defStyle)

	Dchat := display.Chat{}

	var current Idisplay
	current = &Dlogin

	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		panic(err)
	}

	s.SetStyle(defStyle)
	s.Clear()

	current.Draw(s)

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()

		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				quit()
			} else if ev.Key() == tcell.KeyEnter {
				if state == "login" {
					if _, conf := current.Enter(s); conf {
						y, u := Dlogin.You, Dlogin.User
						se := conn.LoginRequest(y, u)
						if se != "" {
							panic(se)
						}
						conn.WaitAddr()
						Dchat.Defautl(y, u)
						current = &Dchat
						state = "chat"
						s.Clear()
						current.Draw(s)
						go conn.ReadResponse(u+": ", msg)
						go readMgsChan(msg, &Dchat, s)

					} else {
						Dlogin.Inputflag(s)
					}
				} else if state == "chat" {
					text, conf := current.Enter(s)
					if !conf {
						continue
					}
					conn.SendMsg(text)
					msg <- "you" + ": " + text
					Dchat.DeleteMsg(s)
				}
			} else if ev.Key() == tcell.KeyDEL || ev.Key() == tcell.KeyBackspace {
				current.Delete(s)
			} else {
				current.Write(s, ev.Rune())
			}
		}
	}
}

func readMgsChan(msg chan string, c *display.Chat, s tcell.Screen) {
	for {
		text := <-msg
		c.WriteMessage(s, text)
		s.Show()
	}
}
