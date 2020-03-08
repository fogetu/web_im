package jobs

import (
	"github.com/astaxie/beego/toolbox"
)

func Register() {
	exampleJobOne()
	exampleJobTwo()
	toolbox.StartTask() //真真切切定时执行。
}
