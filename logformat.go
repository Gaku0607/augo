package augo

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
)

//打印LOG產生的標頭
var LogTitle = "AUGO"

//設置標頭
func SetLogTitle(title string) {
	LogTitle = title
	setDeBugTitle(LogTitle)
}

//獲取標頭
func GetLogTitle() string {
	return LogTitle
}

//關閉Log的顏色輸出
func NoColor() {
	color.NoColor = true
}

//默認LogTimeFormat
var defaultTimeFormat = "2006/01/02-15:04:05"

//獲取當前的TimeFormat
func SetTimeFormat(format string) {
	defaultTimeFormat = format
}

//默認輸出Log方式
var defaultWriter io.Writer = color.Output

//默認輸出Log方式
var defaultErrWriter io.Writer = color.Error

//設置Log的輸出方式
func SetWriter(w io.Writer) {
	defaultWriter = w
}

//設置ErrorLog的輸出方式
func SetErrWriter(w io.Writer) {
	defaultErrWriter = w
}

//log輸出時的格式 可以參考
type LoggerFormatter func(*LoggerParms) string

var (
	defaultFormatter = func(f *LoggerParms) string {
		str := fmt.Sprintf(
			"[%s] %s |%s| ID:%d | %s | Files:%s ",
			LogTitle,
			f.Time,
			f.Type,
			f.RequsetId,
			f.Method,
			f.Files,
		)
		for key, val := range f.Key {
			str = strings.Join([]string{str, fmt.Sprintf("| %s: %v ", key, val)}, "")
		}
		return str
	}

	NilFormatter = func(f *LoggerParms) string {
		return ""
	}

	JSONFormatter = func(f *LoggerParms) string {
		s, _ := json.Marshal(f)
		return string(s)
	}
)

//自定義Log輸出參數
type LogKey map[string]interface{}

const (
	ERROR   = " ERROR "
	REQUEST = " REQUEST "
	RESULT  = " RESULT "
	PANIC   = " PANIC "
)

//Logger.Log()輸出時所帶有的參數 key參數能夠自定義
type LoggerParms struct {
	Type      string `json:"type"`
	Time      string `json:"time"`
	RequsetId int64  `json:"request_Id"`
	Files     string `json:"files"`
	Method    string `json:"method"`
	Key       LogKey `josn:"key"`
}

func (l *LoggerParms) setTypeColor() func(w io.Writer, a ...interface{}) (int, error) {
	switch l.Type {
	case ERROR:
		return color.New(color.FgRed).Fprint
	case REQUEST:
		return color.New(color.FgHiBlue).Fprint
	case RESULT:
		return color.New(color.FgHiGreen).Fprint
	case PANIC:
		return color.New(color.FgHiRed).Fprint
	default:
		return color.New(color.FgHiWhite).Fprint
	}
}

func CreateLogParms(reqId int64, Type, filename, method string, key map[string]interface{}) *LoggerParms {
	l := &LoggerParms{}
	l.RequsetId = reqId
	l.Time = time.Now().Format(defaultTimeFormat)
	l.Type = Type
	l.Files = filename
	l.Method = method
	l.Key = key
	return l
}
