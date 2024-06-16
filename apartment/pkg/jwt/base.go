package jwt

type JwtClaim struct {
	FlatNo int
	Role   string
	Exp    int64
	Email  string
}

type jwtGenerator struct {
	JwtKey string
	Claim  JwtClaim
}

func NewJwtGenerator(claim JwtClaim, key string) jwtGenerator {
	return jwtGenerator{
		JwtKey: key,
		Claim:  claim,
	}
}
