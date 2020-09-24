package health_check

type exampleOneHealthCheck struct {
}

func (dc *exampleOneHealthCheck) Check() error {
	return nil
}
