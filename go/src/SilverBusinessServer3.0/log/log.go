package log

import (
	"SilverBusinessServer/env"
	"os"
	"path"

	"code.google.com/p/log4go"
)

var (
	clevel     = log4go.Level(env.Get("logger.level").(int64))
	cpath      = CreateLogDir(env.Get("logger.file.path").(string))
	cformatc   = env.Get("logger.console.format").(string)
	cformatf   = env.Get("logger.file.format").(string)
	csize      = int(env.Get("logger.file.size").(int64))
	clines     = int(env.Get("logger.file.lines").(int64))
	cmaxbackup = int(env.Get("logger.file.maxbackup").(int64))
)

var (
	wfroot      = newFileLogWriter("root")
	wfhttp      = newFileLogWriter("http")
	wfwebsocket = newFileLogWriter("websocket")
	wform       = newFileLogWriter("orm")
	wfdao       = newFileLogWriter("dao")
	wferror     = newFileLogWriter("error")
	wffs        = newFileLogWriter("fs")
	wc          = newConsoleLogWriter()
)

var (
	// Root Logger
	Root logger = log4go.Logger{
		//"stdout": &log4go.Filter{clevel, wc},
		"froot":  &log4go.Filter{clevel, wfroot},
		"ferror": &log4go.Filter{log4go.ERROR, wferror},
	}
	// ORM Logger
	ORM logger = log4go.Logger{
		//"stdout": &log4go.Filter{clevel, wc},
		"froot":  &log4go.Filter{clevel, wfroot},
		"ferror": &log4go.Filter{log4go.ERROR, wferror},
		"form":   &log4go.Filter{clevel, wform},
	}
	// Dao Logger
	Dao logger = log4go.Logger{
		//"stdout": &log4go.Filter{clevel, wc},
		"froot":  &log4go.Filter{clevel, wfroot},
		"ferror": &log4go.Filter{log4go.ERROR, wferror},
		"fdao":   &log4go.Filter{clevel, wfdao},
	}
	// HTTP Logger
	HTTP logger = log4go.Logger{
		//"stdout": &log4go.Filter{clevel, wc},
		"froot":  &log4go.Filter{clevel, wfroot},
		"ferror": &log4go.Filter{log4go.ERROR, wferror},
		"fhttp":  &log4go.Filter{clevel, wfhttp},
	}
	// WebSocket Logger
	WebSocket logger = log4go.Logger{
		//"stdout":     &log4go.Filter{clevel, wc},
		"froot":      &log4go.Filter{clevel, wfroot},
		"ferror":     &log4go.Filter{log4go.ERROR, wferror},
		"fwebsocket": &log4go.Filter{clevel, wfwebsocket},
	}
	// FS Logger
	FS logger = log4go.Logger{
		//"stdout": &log4go.Filter{clevel, wc},
		"froot":  &log4go.Filter{clevel, wfroot},
		"ferror": &log4go.Filter{log4go.ERROR, wferror},
		"ffs":    &log4go.Filter{clevel, wffs},
	}
)

// Quick use for Root.(method)
var (
	Finest   = Root.Finest
	Fine     = Root.Fine
	Debug    = Root.Debug
	Trace    = Root.Trace
	Info     = Root.Info
	Warn     = Root.Warn
	Error    = Root.Error
	Critical = Root.Critical
)

type logger interface {
	Finest(arg0 interface{}, args ...interface{})
	Fine(arg0 interface{}, args ...interface{})
	Debug(arg0 interface{}, args ...interface{})
	Trace(arg0 interface{}, args ...interface{})
	Info(arg0 interface{}, args ...interface{})
	Warn(arg0 interface{}, args ...interface{}) error
	Error(arg0 interface{}, args ...interface{}) error
	Critical(arg0 interface{}, args ...interface{}) error
}

func newFileLogWriter(name string) *log4go.FileLogWriter {
	fname := path.Join(cpath, name+".log")
	flw := log4go.NewFileLogWriter(fname, false)
	flw.SetFormat(cformatf)
	flw.SetRotate(false)
	flw.SetRotateSize(csize)
	flw.SetRotateLines(clines)
	flw.SetRotateDaily(true)
	flw.SetRotateMaxBackup(cmaxbackup)
	return flw
}

func newConsoleLogWriter() *log4go.ConsoleLogWriter {
	clw := log4go.NewConsoleLogWriter()
	clw.SetFormat(cformatc)
	return clw
}

func CreateLogDir(dir string) string {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			if err = os.Mkdir(dir, os.ModePerm); err != nil {
				return "."
			} else {
				return dir
			}
		} else {
			// other error
			return "."
		}
	} else {
		//exist
		return dir
	}
}

func init() {

	debugMode := env.Get("debug").(bool)
	if debugMode {
		// Root Logger
		Root = log4go.Logger{
			"stdout": &log4go.Filter{clevel, wc},
			"froot":  &log4go.Filter{clevel, wfroot},
			"ferror": &log4go.Filter{log4go.ERROR, wferror},
		}
		// ORM Logger
		ORM = log4go.Logger{
			"stdout": &log4go.Filter{clevel, wc},
			"froot":  &log4go.Filter{clevel, wfroot},
			"ferror": &log4go.Filter{log4go.ERROR, wferror},
			"form":   &log4go.Filter{clevel, wform},
		}
		// Dao Logger
		Dao = log4go.Logger{
			"stdout": &log4go.Filter{clevel, wc},
			"froot":  &log4go.Filter{clevel, wfroot},
			"ferror": &log4go.Filter{log4go.ERROR, wferror},
			"fdao":   &log4go.Filter{clevel, wfdao},
		}
		// HTTP Logger
		HTTP = log4go.Logger{
			"stdout": &log4go.Filter{clevel, wc},
			"froot":  &log4go.Filter{clevel, wfroot},
			"ferror": &log4go.Filter{log4go.ERROR, wferror},
			"fhttp":  &log4go.Filter{clevel, wfhttp},
		}
		// WebSocket Logger
		WebSocket = log4go.Logger{
			"stdout":     &log4go.Filter{clevel, wc},
			"froot":      &log4go.Filter{clevel, wfroot},
			"ferror":     &log4go.Filter{log4go.ERROR, wferror},
			"fwebsocket": &log4go.Filter{clevel, wfwebsocket},
		}
		// FS Logger
		FS = log4go.Logger{
			"stdout": &log4go.Filter{clevel, wc},
			"froot":  &log4go.Filter{clevel, wfroot},
			"ferror": &log4go.Filter{log4go.ERROR, wferror},
			"ffs":    &log4go.Filter{clevel, wffs},
		}
	}
}
