package jobs

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
)

func exampleJobTwo() {
	tk1 := toolbox.NewTask("tk2", "0/20 * * * * * ", func() error {
		fmt.Println("tk2222")
		return nil
	})
	toolbox.AddTask("mytask2", tk1)
}
