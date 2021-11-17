package augo

//當所有Handler都執行完成時刪除Request所有檔案
func DeletFiles() HandlerFunc {
	return func(c *Context) {
		c.Next()
		if err := deletFiles(c.Request.Files); err != nil {
			c.AbortWithError(err)
		}
	}
}

func Recovery(log *Logger) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Log(CreateLogParms(c.Request.Id, PANIC, c.Request.FilesName(), c.Request.method, LogKey{"RecoveryMsg:": r.(error).Error()}))
				c.AbortWithError(r.(error))
				return
			}
		}()
		c.Next()
	}
}
