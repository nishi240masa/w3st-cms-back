package infrastructure

import (
	"context"
	"w3st/domain/repositories"
	"w3st/errors"

	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepositoryImpl(db *gorm.DB) repositories.TransactionRepository {
	return &TransactionRepositoryImpl{DB: db}
}

type contextKey string

const txKey contextKey = "tx"

func (t *TransactionRepositoryImpl) Do(ctx context.Context, f func(ctx context.Context) error) error {
	// この関数は、トランザクションを開始し、f関数を実行します。
	return t.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// トランザクションのコンテキストを作成
		txCtx := context.WithValue(ctx, txKey, tx)

		// トランザクションのコンテキストをf関数に渡す
		if err := f(txCtx); err != nil {
			// エラーが発生した場合、トランザクションをロールバック
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				// ロールバックに失敗した場合、エラーメッセージを返す
				return errors.NewDomainError(errors.TransactionError, "トランザクションのロールバックに失敗しました: "+rollbackErr.Error())
			}
			// ロールバックに成功した場合、エラーメッセージを返す
			return err // f関数からのエラーを返す
		}
		// エラーが発生しなかった場合、トランザクションをコミット
		return nil
	})
}