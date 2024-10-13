package valueobject

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	return Email{
		value: value,
	}, nil
}
