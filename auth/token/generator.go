package token

import "time"

type tokenGeneratorStatic struct {
	value []byte
}

var _ TokenGenerator = (*tokenGeneratorStatic)(nil)

func NewTokenGeneratorStatic(value []byte) TokenGenerator {
	return &tokenGeneratorStatic{
		value: value,
	}
}

func (g *tokenGeneratorStatic) Generate() ([]byte, error) {
	return g.value, nil
}

func (g *tokenGeneratorStatic) TTL() time.Duration {
	return TTL
}

type tokenGeneratorIAM struct {
	account string
	region  string
}

var _ TokenGenerator = (*tokenGeneratorIAM)(nil)

func NewTokenGeneratorIAM(account, region string) TokenGenerator {
	return &tokenGeneratorIAM{
		account: account,
		region:  region,
	}
}

func (g *tokenGeneratorIAM) Generate() ([]byte, error) {
	return Exchange(g.account, g.region)
}

func (g *tokenGeneratorIAM) TTL() time.Duration {
	return TTL
}
