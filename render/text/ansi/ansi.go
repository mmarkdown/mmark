// Copied from https://github.com/go-playground/ansi
// MIT licensed.

package ansi

// ANSI escape sequences
// NOTE: in a standard xterm terminal the light colors will appear BOLD instead of the light variant
const (
	Reset             = "\x1b[0m"
	Italics           = "\x1b[3m"
	Underline         = "\x1b[4m"
	Blink             = "\x1b[5m"
	Inverse           = "\x1b[7m"
	ItalicsOff        = "\x1b[23m"
	UnderlineOff      = "\x1b[24m"
	BlinkOff          = "\x1b[25m"
	InverseOff        = "\x1b[27m"
	Black             = "\x1b[30m"
	DarkGray          = "\x1b[30;1m"
	Red               = "\x1b[31m"
	LightRed          = "\x1b[31;1m"
	Green             = "\x1b[32m"
	LightGreen        = "\x1b[32;1m"
	Yellow            = "\x1b[33m"
	LightYellow       = "\x1b[33;1m"
	Blue              = "\x1b[34m"
	LightBlue         = "\x1b[34;1m"
	Magenta           = "\x1b[35m"
	LightMagenta      = "\x1b[35;1m"
	Cyan              = "\x1b[36m"
	LightCyan         = "\x1b[36;1m"
	Gray              = "\x1b[37m"
	White             = "\x1b[37;1m"
	ResetForeground   = "\x1b[39m"
	BlackBackground   = "\x1b[40m"
	RedBackground     = "\x1b[41m"
	GreenBackground   = "\x1b[42m"
	YellowBackground  = "\x1b[43m"
	BlueBackground    = "\x1b[44m"
	MagentaBackground = "\x1b[45m"
	CyanBackground    = "\x1b[46m"
	GrayBackground    = "\x1b[47m"
	ResetBackground   = "\x1b[49m"
	Bold              = "\x1b[1m"
	BoldOff           = "\x1b[22m"
	Strikethrough     = "\x1b[9m"
	StrikethroughOff  = "\x1b[29m"
)
