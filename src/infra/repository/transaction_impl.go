package infrastructure

import (
	"context"
	"fmt"

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
	// すでにトランザクションが開始されている時はそれを使用する
	if existingTx := ctx.Value(txKey); existingTx != nil {
		return f(ctx)
	}

	// この関数は、トランザクションを開始し、f関数を実行します。
	err := t.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//　トランザクション内でpanicが発生した場合、ロールバックを行う
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				// panicの内容をエラーメッセージとして返す
				err := errors.NewDomainError(errors.TransactionError, "トランザクション中にpanicが発生しました: "+fmt.Sprintf("%v", r))
				// panicの内容をエラーメッセージとして返す
				panic(err)
			}
		}()

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
			return errors.NewDomainError(errors.TransactionError, "トランザクション中にエラーが発生しました: "+err.Error())
		}
		// エラーが発生しなかった場合、トランザクションをコミット
		return nil
	})
	if err != nil {
		// トランザクションのコミットに失敗した場合、エラーメッセージを返す
		return errors.NewDomainError(errors.TransactionError, "トランザクションのコミットに失敗しました: "+err.Error())
	}
	// トランザクションのコミットに成功した場合、nilを返す
	return nil
}
