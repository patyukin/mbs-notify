package usecase

import (
	"context"
	"github.com/patyukin/mbs-pkg/pkg/model"
)

func (u *UseCase) SendAuthSignUpResultMessage(_ context.Context, msg model.AuthSignUpResultMessage) error {

	return nil
}
