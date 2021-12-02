package augo

import "errors"

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

				var (
					key LogKey
					err error
				)

				switch r.(type) {
				case string:
					err = errors.New(r.(string))
					key = LogKey{PANIC_ERROR: err.Error()}

				case error:
					err = r.(error)
					key = LogKey{PANIC_ERROR: err.Error()}

				default:
					err = errors.New("UNKNOW_ERROR")
					key = LogKey{UNKNOW_ERROR: r}

				}

				deletFiles(c.Request.Files)

				log.Log(CreateLogParms(c.Request.Id, PANIC, c.Request.FilesName(), c.Request.method, key))
				c.AbortWithError(err)
			}
		}()
		c.Next()
	}
}
