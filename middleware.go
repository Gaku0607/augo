package augo

import (
	"fmt"
	"os"

	"github.com/Gaku0607/augo/logger"
	"github.com/fatih/color"
)

//當所有Handler都執行完成時刪除Request所有檔案
func DeletFiles() HandlerFunc {
	return func(c *Context) {
		c.Next()
		for _, file := range c.Request.Files {
			if err := os.Remove(file); err != nil {
				c.Error(err)
			}
		}
	}
}

func Recovery() HandlerFunc {
	log := logger.NowLogger()
	f := color.New(color.FgHiRed).Fprint
	text := "[Recovery] ID:%d | Method: %s | RecoveryMsg: %s"
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				log.DebugPrint(f, fmt.Sprintf(text, c.Request.Id, c.Request.method, r.(error).Error()))
				c.AbortWithError(r.(error))
				return
			}
		}()
	}
}
