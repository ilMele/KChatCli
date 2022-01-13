package display

import "github.com/gdamore/tcell"

const(
    inputmark = "-> "
    usermark = "user: "
)

type Login struct{
    You_set bool
    You string
    User string
    left int
    X int
    Y int
    Style tcell.Style
}

func (l *Login) Default(s tcell.Style){
    l.You_set = false
    l.You = ""
    l.User = ""
    l.left = len(youmark)
    l.X = 4
    l.Y = 2
    l.Style = s 
}

func (l *Login) Write(s tcell.Screen, r rune){
    if !l.You_set{
        l.You += string(r)
        s.SetContent(l.X + l.left, l.Y, r, nil, l.Style)
        l.left++
        return
    }

    l.User += string(r)
    s.SetContent(l.X + l.left, l.Y+1, r, nil, l.Style)
    l.left++
}

func (l *Login) Draw(s tcell.Screen){
    x := l.X
    y := l.Y
    
    for _, r := range youmark + l.You{
        s.SetContent(x, y, r, nil, l.Style)
        x++
    }
    x = l.X
    y++
    for _, r := range usermark + l.User{
        s.SetContent(x, y, r, nil, l.Style)
        x++
    }
    l.Inputflag(s)
}

func (l *Login) Enter() (string, bool){
    if !l.You_set{
        l.You_set = true
        l.left = len(usermark) 
        return "", false
    }
    return "", true
}

func (l *Login) Delete() (int, int){
    if l.You_set{
        dim := len(l.User) -1
        if dim > 0{
            l.User = l.User[:dim]
            l.left--
            return l.X + l.left, l.Y + 1
        }
        if dim == 0{
            l.User = ""
            l.left--
            return l.X + l.left, l.Y + 1
        }
    }

    dim := len(l.You) - 1
    if dim > 0{
        l.You = l.You[:dim]
        l.left--
        return l.X + l.left, l.Y
    }
    if dim == 0{
        l.You = ""
        l.left--
        return l.X + l.left, l.Y
    }
    return -1, -1
}

func (l *Login) Inputflag(s tcell.Screen){
    x := l.X - len(inputmark)
    if l.You_set{
        for _, r := range inputmark{
            s.SetContent(x, l.Y + 1, r, nil, l.Style)
            x++
        }
        x--
        for i := len(inputmark); i > 0; i--{
            s.SetContent(x, l.Y, ' ', nil, l.Style)
            x--
        }
    return
    }
    for _, r := range inputmark{
        s.SetContent(x, l.Y, r, nil, l.Style)
        x++
    }
    
    
}
