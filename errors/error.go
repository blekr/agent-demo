package errors

type AppError struct {
	Code string
	Message string
}

func (err *AppError) Error() string {
	return err.Code + ":" + err.Message
}
