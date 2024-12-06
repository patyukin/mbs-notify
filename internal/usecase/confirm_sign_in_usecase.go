package usecase

import (
	"context"

	"github.com/patyukin/mbs-pkg/pkg/model"
)

func (u *UseCase) ConfirmSignIn(ctx context.Context, message model.AuthSignInCode) error {
	return nil
}
