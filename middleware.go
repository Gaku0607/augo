package augo

import (
	"os"

	"github.com/Gaku0607/augo/logger"
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

func Recovery(log *logger.Logger) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Log(logger.CreateLogParms(c.Request.Id, logger.Recovery, c.Request.FilesName(), c.Request.method, logger.LogKey{"RecoveryMsg:": r.(error).Error()}))
				c.AbortWithError(r.(error))
				return
			}
		}()
		c.Next()
	}
}
