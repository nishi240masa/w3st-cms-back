package usecase

import (
	"context" // context をインポート
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
	"w3st/interfaces/services"
)

type UserUsecase interface {
	Create(newUser *models.Users, ctx context.Context) (models.Token, error)
	FindByEmail(email string) (models.Token, error)
	FindByID(userID string) (*models.Users, error)
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

func (u *userUsecase) Create(newUser *models.Users, ctx context.Context) (models.Token, error) {
	var generatedToken models.Token // トークンを格納する変数を宣言

	// トランザクションを開始
	err := u.tx.Do(ctx, func(txCtx context.Context) error {
		_, err := u.userRepo.FindByEmail(ctx, newUser.Email)

		// ユーザーが見つからなかった場合は新規作成
		if err != nil && errors.Is(err, &myerrors.DomainError{ErrType: myerrors.QueryDataNotFoundError}) {
			err = u.userRepo.Create(ctx, newUser) // context を渡す
			if err != nil {
				// リポジトリで技術的なエラーが発生した場合
				return myerrors.NewDomainError(myerrors.QueryError, err.Error())
			}

			// token生成
			token, err := u.authService.GenerateToken(newUser.ID)
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
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			// ドメインエラーの場合はそのまま返す
			return "", domainErr
		}
		// その他のエラーの場合は一般的なエラーメッセージを返す
		return "", myerrors.NewDomainError(myerrors.ErrorUnknown, "トランザクション中にエラーが発生しました")
	}
	// トランザクションが成功した場合は、生成されたトークンを返す
	return generatedToken, nil
}

func (u *userUsecase) FindByEmail(email string) (models.Token, error) {
	user, err := u.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return "", err
	}
	token, err := u.authService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) FindByID(userID string) (*models.Users, error) {
	user, err := u.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
