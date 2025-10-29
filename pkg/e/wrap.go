package e

import (
	"context"
	"fmt"
	"time"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%v, %w", msg, err)
}

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}

func Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}
