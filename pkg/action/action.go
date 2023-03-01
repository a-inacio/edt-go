package action

func FromError(err error) (Result, error) {
	return nil, err
}

func Nothing() (Result, error) {
	return nil, nil
}
