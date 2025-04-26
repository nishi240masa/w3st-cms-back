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

	err := t.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(errors.NewDomainErrorWithMessage(
					errors.TransactionError,
					fmt.Sprintf("トランザクション中にpanicが発生しました: %v", r),
				))
			}
		}()

		txCtx := context.WithValue(ctx, txKey, tx)

		if err := f(txCtx); err != nil {
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				return errors.NewDomainErrorWithMessage(
					errors.TransactionError,
					fmt.Sprintf("トランザクションのロールバックに失敗しました: %v", rollbackErr),
				)
			}
			return errors.NewDomainErrorWithMessage(
				errors.TransactionError,
				fmt.Sprintf("トランザクション中にエラーが発生しました: %v", err),
			)
		}

		return nil
	})
	if err != nil {
		return errors.NewDomainErrorWithMessage(
			errors.TransactionError,
			fmt.Sprintf("トランザクションのコミットに失敗しました: %v", err),
		)
	}

	return nil
}
