package ansi

import (
	"fmt"
	"os"
)

var (
	start     = "\033["
	reset     = "\033[0m"
	bold      = "1;"
	blink     = "5;"
	underline = "4;"
	inverse   = "7;"
)

type Style byte

const (
	Normal    Style = 0x00
	Bold      Style = 0x01
	Blink     Style = 0x02
	Underline Style = 0x04
	Invert    Style = 0x08
	Intense   Style = 0x10
)

type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Colorless
)

const (
	normalFg  = 30
	intenseFg = 90
	normalBg  = 40
	intenseBg = 100
)

type Colorer struct {
	enabled bool
}

func NewFor(f *os.File) *Colorer {
	return &Colorer{isTTY(f.Fd())}
}

func (c *Colorer) Fg(s string, color Color, style Style) string {
	return c.FgBg(s, color, style, Colorless, Normal)
}

func (c *Colorer) FgBg(s string, fgColor Color, fgStyle Style, bgColor Color, bgStyle Style) string {
	if !c.enabled {
		return s
	}

	buf := make([]byte, 0, 24)
	buf = append(buf, start...)

	if fgStyle&Bold != 0 {
		buf = append(buf, bold...)
	}

	if fgStyle&Blink != 0 {
		buf = append(buf, blink...)
	}

	if fgStyle&Underline != 0 {
		buf = append(buf, underline...)
	}

	if fgStyle&Invert != 0 {
		buf = append(buf, inverse...)
	}

	var fgBase int
	if fgStyle&Intense == 0 {
		fgBase = normalFg
	} else {
		fgBase = intenseFg
	}
	buf = append(buf, fmt.Sprintf("%d;", fgBase+int(fgColor))...)

	if bgColor != Colorless {
		var bgBase int
		if bgStyle&Intense == 0 {
			bgBase = normalBg
		} else {
			bgBase = intenseBg
		}
		buf = append(buf, fmt.Sprintf("%d;", bgBase+int(bgColor))...)
	}

	buf = append(buf[:len(buf)-1], "m"...)
	return string(buf) + s + reset
}
