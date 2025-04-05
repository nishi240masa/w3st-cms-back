package repositories

import "context"


type TransactionRepository interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}