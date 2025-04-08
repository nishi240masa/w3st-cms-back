package usecase

import (
	"context" // context をインポート
	"w3st/domain/models"
	"w3st/domain/repositories"
	"w3st/errors"
	"w3st/interfaces/services"
	"w3st/utils"
)

type UserUsecase interface {
	Create(newUser *models.User, ctx context.Context ) (models.Token, *errors.DomainError)
	FindByEmail(email string) (models.Token, *errors.DomainError)
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

func (u *userUsecase) Create(newUser *models.User,ctx context.Context) (models.Token, *errors.DomainError) {
	var generatedToken models.Token // トークンを格納する変数を宣言

	// トランザクションを開始
	err := u.tx.Do(ctx, func(txCtx context.Context) error {
		_, findErr := u.userRepo.FindByEmail(ctx, newUser.Email) 
		if findErr == nil {
			// ユーザーがすでに存在する場合
			return errors.NewDomainError(errors.AlreadyExist, "このメールアドレスは既に登録されています") 
		} else if findErr.ErrType != errors.QueryDataNotFoundError {
			// ユーザーが見つからなかった場合以外のエラー
			return findErr
		}
		createErr := u.userRepo.Create(ctx, newUser) // context を渡す
		if createErr != nil {
			// リポジトリで技術的なエラーが発生した場合
			return errors.NewDomainError(errors.QueryError, createErr.Error())
		}

		// token生成
		stringID := utils.UuidToString(newUser.ID)
		token, genErr := u.authService.GenerateToken(stringID)
		if genErr != nil {
			return genErr
		}
		generatedToken = token // トランザクション内で生成したトークンを格納
		return nil             // エラーがなければ nil を返す
	})

	if err != nil {
		// トランザクション内でエラーが発生した場合
		if domainErr, ok := err.(*errors.DomainError); ok {
			// ドメインエラーの場合はそのまま返す
			return "", domainErr
		}
		// その他のエラーの場合は一般的なエラーメッセージを返す
		return "", errors.NewDomainError(errors.ErrorUnknown, "トランザクション中にエラーが発生しました")
	}
	return generatedToken, nil // トランザクションが成功したら生成されたトークンを返す
}

func (u *userUsecase) FindByEmail(email string) (models.Token, *errors.DomainError) {
	user, err := u.userRepo.FindByEmail(context.Background(), email) 
	if err != nil {
		return "",err
	}
	token, genErr := u.authService.GenerateToken(user.ID.String())
	if genErr != nil {
		return "", genErr
	}
	return token, nil
}