package usecase

import (
	"context" // context をインポート
	"errors"
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
	"w3st/interfaces/services"
	"w3st/utils"
)

type UserUsecase interface {
	Create(newUser *models.Users, ctx context.Context) (models.Token, *myerrors.DomainError)
	FindByEmail(email string) (models.Token, *myerrors.DomainError)
}

type userUsecase struct {
	userRepo    repositories.UserRepository
	authService services.AuthService
	tx          repositories.TransactionRepository
}

func NewUserUsecase(userRepo repositories.UserRepository, authService services.AuthService, txRepo repositories.TransactionRepository) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		authService: authService,
		tx:          txRepo,
	}
}

func (u *userUsecase) Create(newUser *models.Users, ctx context.Context) (models.Token, *myerrors.DomainError) {
	var generatedToken models.Token // トークンを格納する変数を宣言

	// トランザクションを開始
	err := u.tx.Do(ctx, func(txCtx context.Context) error {
		_, err := u.userRepo.FindByEmail(ctx, newUser.Email)

		// ユーザーが見つからなかった場合は新規作成
		if errors.Is(err, &myerrors.DomainError{ErrType: myerrors.QueryDataNotFoundError}) {
			err = u.userRepo.Create(ctx, newUser) // context を渡す
			if err != nil {
				// リポジトリで技術的なエラーが発生した場合
				return myerrors.NewDomainError(myerrors.QueryError, err.Error())
			}

			// token生成
			stringID := utils.UuidToString(newUser.ID)
			token, err := u.authService.GenerateToken(stringID)
			if err != nil {
				return err
			}
			generatedToken = token // トランザクション内で生成したトークンを格納
			return nil
		} else if err != nil {
			// リポジトリで技術的なエラーが発生した場合
			return myerrors.NewDomainError(myerrors.QueryError, err.Error())
		}
		// ユーザーが見つかった場合はエラーを返す
		return myerrors.NewDomainError(myerrors.AlreadyExist, "このメールアドレスは既に登録されています")
	})

	if err != nil {
		// トランザクション内でエラーが発生した場合
		if err, ok := err.(*myerrors.DomainError); ok {
			// ドメインエラーの場合はそのまま返す
			return "", err
		}
		// その他のエラーの場合は一般的なエラーメッセージを返す
		return "", myerrors.NewDomainError(myerrors.ErrorUnknown, "トランザクション中にエラーが発生しました")
	}
	return generatedToken, nil // トランザクションが成功したら生成されたトークンを返す
}

func (u *userUsecase) FindByEmail(email string) (models.Token, *myerrors.DomainError) {
	user, err := u.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return "", err
	}
	token, err := u.authService.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}
	return token, nil
}
