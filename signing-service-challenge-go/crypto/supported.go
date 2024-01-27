package crypto

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

var SupportedGenerators = []domain.KeyPairGenerator{
	ECCGenerator{},
	RSAGenerator{},
}

func FindKeyPairGenerator(algorithmName string) (domain.KeyPairGenerator, bool) {
	for _, generator := range SupportedGenerators {
		if generator.AlgorithmName() == algorithmName {
			return generator, true
		}
	}

	return nil, false
}
