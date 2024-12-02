package validator

type ValidationErrors struct {
	Errors map[string]string
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make(map[string]string),
	}
}

func (v *ValidationErrors) Add(field, message string) {
	v.Errors[field] = message
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *ValidationErrors) Error() string {
	return "Validation failed."
}
