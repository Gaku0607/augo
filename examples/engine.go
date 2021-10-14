package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gaku0607/augo"
	"github.com/Gaku0607/augo/logger"
)

func main() {
	//設置系統
	augo.SetSystemVersion(augo.MacOS)

	c := augo.NewCollector()
	c.Logger.Config.Format = logger.JSONFormatter
	c.Use(augo.Recovery(), print1())
	{
		g1 := c.Group("/Users/YourPath1")
		g1.Handler("/Test1", print2())

	}

	{
		g2 := c.Group("/Users/gaku/IRIS系統測試檔案/出倉單")
		g2.Handler("/個例", print3())
	}
	//設置定時
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	engine := augo.NewEngine(
		augo.MaxThread(5), //線程數
		augo.SetContext(ctx),
		augo.ScanIntval(time.Millisecond*1000), //提交
		augo.SetCollector(c),
	)

	engine.Run()

	engine.Wait()

	if err := engine.Close(); err != nil {
		panic(err.Error())
	}

}

func print1() augo.HandlerFunc {
	return func(d *augo.Context) {
		fmt.Println("test:", "1")
	}
}

func print2() augo.HandlerFunc {
	return func(d *augo.Context) {
		fmt.Println("test:", "2")
	}
}

func print3() augo.HandlerFunc {
	return func(d *augo.Context) {
		// fmt.Println("test:", "3")
		d.Error(errors.New("err3"))
	}
}
