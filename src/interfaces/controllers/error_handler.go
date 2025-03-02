package controllers

import (
	"net/http"
	"w3st/errors"
	"w3st/infra/logger"

	"github.com/bufbuild/connect-go"
)


func ErrorHandle(domainErr *errors.DomainError) *connect.Error {

	switch domainErr.ErrType {
	case errors.InvalidParameter:
		return connect.NewError(connect.CodeInvalidArgument, domainErr)
	case errors.UnPemitedOperation:
		return connect.NewError(connect.CodePermissionDenied, domainErr)
	case errors.AlreadyExist:
		return connect.NewError(connect.CodeAlreadyExists, domainErr)
	case errors.RepositoryError, errors.QueryError:
		logger.Error(domainErr.Error())
		return connect.NewError(connect.CodeInternal, domainErr)
	case errors.QueryDataNotFoundError:
		logger.Error(domainErr.Error())
		return connect.NewError(connect.CodeNotFound, domainErr)
	default:
		logger.Error(domainErr.Error())
		return connect.NewError(connect.CodeUnknown, domainErr)
	}
}

// 補助関数例: connect.Error のコードをHTTPステータスコードに変換する
func HttpStatusCodeFromConnectCode(code connect.Code) int {
    switch code {
    case connect.CodeInvalidArgument:
        return http.StatusBadRequest
    case connect.CodePermissionDenied:
        return http.StatusForbidden
    case connect.CodeAlreadyExists:
        return http.StatusConflict
    case connect.CodeInternal:
        return http.StatusInternalServerError
    case connect.CodeNotFound:
        return http.StatusNotFound
    default:
        return http.StatusInternalServerError
    }
}