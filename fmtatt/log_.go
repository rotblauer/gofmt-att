package fmtatt

type LogLevel int
const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
	Ridic
)

type LogConfig struct {
	Level int
	Logger string
}

var DefaultLogConfig = LogConfig{
	Level: 3,
	Logger: "stderr",
}

type Verbosably interface {
	SetLevel(level int)

	F(args []interface{})
	Ff(fmt string, args []interface{})

	E(args []interface{})
	Ef(fmt string, args []interface{})

	W(args []interface{})
	Wf(fmt string, args []interface{})

	I(args []interface{})
	If(fmt string, args []interface{})

	D(args []interface{})
	Df(fmt string, args []interface{})

	R(args []interface{})
	Rf(fmt string, args []interface{})
}
