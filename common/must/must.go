package must

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2(v interface{}, err error) interface{} {
	if err != nil {
		Must(err)
	}
	return v
}

func Error2(v interface{}, err error) error {
	return err
}
