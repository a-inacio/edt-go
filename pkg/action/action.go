package action

import "context"

func FromError(err error) (Result, error) {
	return nil, err
}

func Nothing() (Result, error) {
	return nil, nil
}

func DoNothing(ctx context.Context) (Result, error) {
	return nil, nil
}
