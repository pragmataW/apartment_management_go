package randomkeygen

import "math/rand"

func (k keyGenerator) GenerateRandomKey() string {
	charset := "abcdefghijklmnoprstuvyzwqxABCDEFGHIJKLMNOPRSTUVYZWQZ1234567890"

	ret := ""
	for i := 0; i < k.digit; i++ {
		randIndex := rand.Intn(len(charset))
		ret += string(charset[randIndex])
	}
	return ret
}
