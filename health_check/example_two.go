package health_check

type exampleTwoHealthCheck struct {
}

func (dc *exampleTwoHealthCheck) Check() error {
	return nil
}
