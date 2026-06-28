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
	ErrFileOpen            = errors.New("file failed to open")
	ErrFileWatch           = errors.New("file watch failed")
	ErrDatabaseDown        = errors.New("security database is temporarily unavailable")
	ErrAccessDenied        = errors.New("insufficient permissions to access this file")
	ErrUnsupportedPlatform = errors.New("platform unsupported")
	ErrHashing             = errors.New("encountered error in hashing")
	ErrLocking             = errors.New("error while scanning file")
	ErrFileDeletion        = errors.New("file could not be deleted")
	MarkErrFileDeletion    = errors.New("file could not be marked for deletion")
	ErrDataCopy            = errors.New("error occured while copying")
)

var ErrorMap = map[error]*AppError{

	ErrFileOpen: {
		Code:    100,
		Message: ErrFileOpen.Error(),
	},
	ErrLocking: {
		Code:    101,
		Message: ErrDatabaseDown.Error(),
	},
	ErrFileDeletion: {
		Code:    102,
		Message: ErrFileDeletion.Error(),
	},
	MarkErrFileDeletion: {
		Code:    103,
		Message: MarkErrFileDeletion.Error(),
	},
	ErrFileWatch: {
		Code:    104,
		Message: ErrFileWatch.Error(),
	},
	ErrHashing: {
		Code:    200,
		Message: ErrHashing.Error(),
	},
	ErrDataCopy: {
		Code:    201,
		Message: ErrDataCopy.Error(),
	},
	ErrDatabaseDown: {
		Code:    400,
		Message: ErrDatabaseDown.Error(),
	},
	ErrUnsupportedPlatform: {
		Code:    600,
		Message: ErrUnsupportedPlatform.Error(),
	},
}
