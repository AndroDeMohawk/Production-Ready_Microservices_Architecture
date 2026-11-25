package model

type compatibilityChecker struct{}

func NewCompatibilityChecker() *compatibilityChecker {
	return &compatibilityChecker{}
}

func (c *compatibilityChecker) Check(parts []Part) error {
	panic("не реализовано")
}
