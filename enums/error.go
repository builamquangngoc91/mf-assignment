package enums

type ErrorCode int8

const (
	NotRowsAffected string = "no rows affected"

	BadRequest ErrorCode = iota + 1
	InternalError
)
