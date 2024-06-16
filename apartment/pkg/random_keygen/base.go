package randomkeygen

type keyGenerator struct {
	digit int
}

func NewKeygen (digit int) keyGenerator {
	return keyGenerator{
		digit: digit,
	}
}