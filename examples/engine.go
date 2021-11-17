package main

import (
	"fmt"
	"time"

	"github.com/Gaku0607/augo"
)

func main() {
	//設置系統
	augo.SetSystemVersion(augo.MacOS)

	augo.SetLogTitle("GAKU")

	c := augo.NewCollector()
	c.Use(augo.Recovery(c.Logger), print1())
	{
		g1 := c.Group("/Users/YourPath")
		g1.HandlerWithVisit("/Method", print2())

	}

	{
		c.Handler("/Users/YourPath", false, print3())
	}

	engine := augo.NewEngine(
		augo.MaxThread(5),                      //線程數
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
		fmt.Println("test:", "3")
		// d.Error(errors.New("err3"))
	}
}
