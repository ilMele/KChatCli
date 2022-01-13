package display

import "github.com/gdamore/tcell"

const (
	msgmark = "> "
	youmark = "you: "
)

type Chat struct {
	You     string
	User    string
	Msg     string
	row     int
	msgrow  int
	left    int
	msgleft int
	top     int

	style tcell.Style
}

func (c *Chat) Defautl(y, u string) {
	c.You = y
	c.User = u
	c.Msg = ""
	c.row = 4
	c.msgrow = 19
	c.left = 3
	c.msgleft = c.left + len(msgmark)
	c.top = 2
}

func (c *Chat) Write(s tcell.Screen, r rune) {
	s.SetContent(c.msgleft, c.msgrow, r, nil, c.style)
	c.msgleft++
	c.Msg += string(r)
}

func (c *Chat) Draw(s tcell.Screen) {
	x := 1
	for _, r := range youmark + c.You {
		s.SetContent(x, c.top, r, nil, c.style)
		x++
	}

	x = len(youmark+c.You) + 2
	for _, r := range "user: " + c.User {
		s.SetContent(x, c.top, r, nil, c.style)
		x++
	}

	s.SetContent(c.left, c.msgrow, rune(msgmark[0]), nil, c.style)
	s.SetContent(c.left+1, c.msgrow, rune(msgmark[1]), nil, c.style)
}

func (c *Chat) Enter() (string, bool) {
	if len(c.Msg) > 0 {
		return c.Msg, true
	}
	return "", false
}

func (c *Chat) Delete() (int, int) {
	dim := len(c.Msg) - 1
	if dim > 0 {
		c.Msg = c.Msg[:dim]
		c.msgleft--
		return c.msgleft, c.msgrow
	}
	if dim == 0 {
		c.Msg = ""
		c.msgleft--
		return c.msgleft, c.msgrow
	}
	return -1, -1
}

func (c *Chat) WriteMessage(s tcell.Screen, m string) {
	x := c.left
	if c.row == c.msgrow-1 {
		c.clearChat(s)
	}
	for _, r := range m {
		s.SetContent(x, c.row, r, nil, c.style)
		x++
	}
	c.row++
}

func (c *Chat) DeleteMsg(s tcell.Screen) {
	min := c.left + len(msgmark)
	for x := c.msgleft; x >= min; x-- {
		s.SetContent(x, c.msgrow, ' ', nil, c.style)
	}
	c.Msg = ""
	c.msgleft = min
}

func (c *Chat) clearChat(s tcell.Screen) {
	w, _ := s.Size()

	for y := c.top + 1; y < c.msgrow; y++ {
		for x := w; x > c.left-1; x-- {
			s.SetContent(x, y, ' ', nil, c.style)
		}
	}
	c.row = 4
}
