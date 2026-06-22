package apperrors

import "errors"

type AppError struct {
	Err     error  // actual error
	Message string // message to display
	Code    int    // internal error code
}

// mandatory method to implement error typing
func (ae *AppError) Error() string {
	if ae.Err != nil {
		return ae.Err.Error()
	}
	return ae.Message
}

func Map(err error) *AppError {

	for target, app := range ErrorMap {
		if errors.Is(err, target) {
			return &AppError{
				Code:    app.Code,
				Message: app.Message,
				Err:     err,
			}
		}
	}
	return &AppError{
		Code:    -1,
		Message: "Unexpected error.",
		Err:     err,
	}
}

// standard errors
var (
	ErrDatabaseDown        = errors.New("security database is temporarily unavailable")
	ErrAccessDenied        = errors.New("insufficient permissions to access this file")
	ErrUnsupportedPlatform = errors.New("Platform unsupported")
	ErrHashing             = errors.New("Encountered error in hashing")
	ErrLocking             = errors.New("error while scanning file")
)

var ErrorMap = map[error]*AppError{
	ErrDatabaseDown: {
		Code:    500,
		Message: "Engine Warning: The signature database is unavailable.",
	},
	ErrUnsupportedPlatform: {
		Code:    600,
		Message: "Unsupported Platform detected!",
	},
	ErrHashing: {
		Code:    100,
		Message: "Hashing failed",
	},
	ErrLocking: {
		Code:    101,
		Message: "File Modified",
	},
}
