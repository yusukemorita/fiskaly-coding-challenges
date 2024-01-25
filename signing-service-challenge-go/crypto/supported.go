package crypto

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

var SupportedAlgorithms = []domain.SignatureAlgorithm{
	ECCAlgorithm{},
	RSAAlgorithm{},
}
