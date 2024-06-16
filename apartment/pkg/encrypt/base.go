package encrypt

type encryptor struct {
	Key []byte
}

func NewEncryptor(key string) *encryptor {
	return &encryptor{Key: []byte(key)}
}
