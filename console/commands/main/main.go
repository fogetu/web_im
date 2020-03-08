package main

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
)

func main() {
	var endRunning = make(chan bool, 1)
	tk1 := toolbox.NewTask("tk1", "0/10 * * * * * ", func() error {
		fmt.Println("tk1111")
		return nil
	})
	toolbox.AddTask("mytask", tk1)
	toolbox.StartTask() //真真切切定时执行。
	<-endRunning
}
