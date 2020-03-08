package jobs

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
)

func exampleJobOne() {
	tk1 := toolbox.NewTask("tk1", "0/10 * * * * * ", func() error {
		fmt.Println("tk1111")
		return nil
	})
	toolbox.AddTask("mytask1", tk1)
}
