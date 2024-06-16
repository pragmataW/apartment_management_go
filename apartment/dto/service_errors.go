package dto

type PasswordMatchError struct{
	Message string
}

func (e PasswordMatchError) Error() string{
	return e.Message
}

type UserDoesNotExists struct{
	Message string
}

func (e UserDoesNotExists) Error() string{
	return e.Message
}

type PayDayRangeError struct{
	Message string
}

func (e PayDayRangeError) Error() string{
	return e.Message
}

type SendMailError struct{
	Message string
}

func (e SendMailError) Error() string{
	return e.Message
}