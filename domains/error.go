package domains

import (
	"banking-service/enums"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ErrorResp struct {
		Message string `json:"message"`
	}
)

type XError struct {
	Err       error
	ErrorCode enums.ErrorCode
}

func NewXError(err error, errCode enums.ErrorCode) XError {
	return XError{
		Err:       err,
		ErrorCode: errCode,
	}
}

func (xerror XError) Is(errCode enums.ErrorCode) bool {
	return xerror.ErrorCode == errCode
}

func (xerror XError) Error() string {
	return xerror.Err.Error()
}

func (xerror XError) Response(c *gin.Context) {
	switch xerror.ErrorCode {
	case enums.BadRequest:
		c.JSON(http.StatusBadRequest, ErrorResp{
			Message: xerror.Err.Error(),
		})
	case enums.InternalError:
		c.JSON(http.StatusInternalServerError, ErrorResp{
			Message: xerror.Err.Error(),
		})
	}
}
