package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	myerrors "w3st/errors"
	"w3st/infra/logger"

	"github.com/bufbuild/connect-go"
)

func ErrorHandle(domainErr *myerrors.DomainError) *connect.Error {
	switch domainErr.ErrType {
	// 技術的なエラー
	case myerrors.InvalidParameter:
		return connect.NewError(connect.CodeInvalidArgument, domainErr)
		// ビジネスロジックエラー
	case myerrors.UnPermittedOperation:
		return connect.NewError(connect.CodePermissionDenied, domainErr)
		// 既に存在するエラー
	case myerrors.AlreadyExist:
		return connect.NewError(connect.CodeAlreadyExists, domainErr)
		// リポジトリで技術的なエラーが発生した場合
	case myerrors.RepositoryError, myerrors.QueryError:
		logger.Error(domainErr.Error())
		return connect.NewError(connect.CodeInternal, domainErr)
		// ユーザーが見つからなかった場合
	case myerrors.QueryDataNotFoundError:
		logger.Error(domainErr.Error())
		return connect.NewError(connect.CodeNotFound, domainErr)
		// トランザクションエラー
	case myerrors.TransactionError:
		logger.Error(domainErr.Error())
		return connect.NewError(connect.CodeInternal, domainErr)
		// その他のエラー
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

// handlerのエラーハンドリング
func ErrorHandler(c *gin.Context, err error) {
	// エラーが connect.Error 型かどうかを確認
	var domainErr *myerrors.DomainError
	if errors.As(err, &domainErr) {
		// connect.Error に変換
		connectErr := ErrorHandle(domainErr)

		// HTTP ステータスコードを取得
		httpStatusCode := HttpStatusCodeFromConnectCode(connectErr.Code())

		// レスポンスを返す
		c.JSON(httpStatusCode, gin.H{"error": connectErr.Message()})
	} else {
		// その他のエラーは500 Internal Server Errorとして処理
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
}
