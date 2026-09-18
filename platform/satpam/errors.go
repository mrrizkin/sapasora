package satpam

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string { return e.Message }

// Assert checks if the given condition is true, and if not, it panics with a message for ForbiddenError
func Assert(cond bool, errMsg string) {
	if !cond {
		panic(&ForbiddenError{Message: errMsg})
	}
}
