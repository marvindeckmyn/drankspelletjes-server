// The log package contains a generic logger interface for applications. The
// system will automatically assign codes to logs so that recurring messages
// can easily be filtered irregardles of dynamic data fields.
package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"runtime"
	"strings"
	"time"
)

type Format int64

// Supported log formats
const (
	PlainText Format = iota
	Json
)

type Destination int64

// Supported log destinations
const (
	Stdout Destination = iota
	File
)

type level int64

// Internal log level enumeration
const (
	debug level = iota
	info
	warning
	error
)

// Log struct which contains log settings.
type Log struct {
	format      Format
	destination Destination
	file        string
}

// Package level exported logger
var log = Log{
	format:      PlainText,
	destination: Stdout,
}

// String returns the string representation of a log level.
func (lvl level) string() string {
	switch lvl {
	case debug:
		return "DBUG"

	case info:
		return "INFO"

	case warning:
		return "WARN"

	case error:
		return "ERRO"
	}

	return ""
}

// ColorANSI returns the ANSI color escape string for a given log level.
func (lvl level) colorANSI() string {
	switch lvl {
	case debug:
		return "\033[37m"

	case info:
		return "\033[32m"

	case warning:
		return "\033[33m"

	case error:
		return "\033[31m"
	}

	return "\033[37m"
}

// ShortCode returns the first letter of every log level.
func (lvl level) shortCode() byte {
	switch lvl {
	case debug:
		return 'D'

	case info:
		return 'I'

	case warning:
		return 'W'

	case error:
		return 'E'
	}

	return ' '
}

// Log logs a message to the set up output in the selected format.
func (l *Log) log(lvl level, format string, args ...interface{}) {
	now := time.Now().Format(time.RFC3339)

	crc := crc32.ChecksumIEEE([]byte(format))
	code := fmt.Sprintf("%c%08X", lvl.shortCode(), crc)
	_, file, line, _ := runtime.Caller(2)

	file = file[strings.LastIndex(file, "/"):]

	var buf []byte

	switch l.format {
	case PlainText:
		var b bytes.Buffer
		fmt.Fprintf(&b, "%s[%s][%s][%s][%s:%d] ", lvl.colorANSI(), now, code, lvl.string(), file, line)
		fmt.Fprintf(&b, format, args...)
		fmt.Fprint(&b, "\033[0m")
		buf = b.Bytes()

	case Json:
		data := map[string]string{
			"time": now,
			"lvl":  lvl.string(),
			"msg":  fmt.Sprintf(format, args...),
			"code": code,
			"line": fmt.Sprintf("%d", line),
			"file": file,
		}

		buf, _ = json.Marshal(data)
	}

	switch l.destination {
	case Stdout:
		fmt.Println(string(buf))
	case File:
		//TODO
	}
}

// SetFormat selects the log output format.
func SetFormat(fmt Format) {
	log.format = fmt
}

// SetDestination selects the output destination of the logger.
func SetDestination(dst Destination, file string) {
	log.destination = dst
	log.file = file
}

// Debug logs a debug message.
func Debug(format string, args ...interface{}) {
	log.log(debug, format, args...)
}

// Info logs an informational message.
func Info(format string, args ...interface{}) {
	log.log(info, format, args...)
}

// Warning logs a warning message.
func Warning(format string, args ...interface{}) {
	log.log(warning, format, args...)
}

// Error logs an error message.
func Error(format string, args ...interface{}) {
	log.log(error, format, args...)
}
