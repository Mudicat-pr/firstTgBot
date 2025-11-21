package e

import (
	"context"
	"log/slog"
	"time"
)

func Wrap(msg string, err error) error {
	slog.Error(msg, "Error", err)
	return err
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
