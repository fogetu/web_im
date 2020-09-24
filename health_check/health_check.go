package health_check

import (
	"github.com/astaxie/beego/toolbox"
)

func Register() {
	toolbox.AddHealthCheck("exampleOneHealthCheck", &exampleOneHealthCheck{})
	toolbox.AddHealthCheck("exampleTwoHealthCheck", &exampleTwoHealthCheck{})
}
