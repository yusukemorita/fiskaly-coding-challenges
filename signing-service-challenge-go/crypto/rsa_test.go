package crypto

import (
	"crypto/rsa"
	"testing"
)

func TestRSAKeyPair_Sign(t *testing.T) {
	generator := RSAGenerator{}
	keyPair, err := generator.generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signature, err := keyPair.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	digest, err := ComputeHashDigest([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	err = rsa.VerifyPSS(keyPair.Public, HashFunction, digest, signature, nil)
	if err != nil {
		t.Errorf("signature verification failed: %s", err)
	}
}

func TestRSAKeyPair_EncodedPublicKey(t *testing.T) {
	// encodedPublicKey, encodedPrivateKey are a key pair pre-generated
	// for this test
	encodedPublicKey := `-----BEGIN RSA_PUBLIC_KEY-----
MEgCQQCxbNTi8ctLVekhxdanVhtnlW06idj6FZUIheaAmjEa+8LmHLziMVCaqcF9
eN+H7mE1Mj1NSx/WIA+131XJn8wfAgMBAAE=
-----END RSA_PUBLIC_KEY-----
`
	encodedPrivateKey := `-----BEGIN RSA_PRIVATE_KEY-----
MIIBPAIBAAJBALFs1OLxy0tV6SHF1qdWG2eVbTqJ2PoVlQiF5oCaMRr7wuYcvOIx
UJqpwX1434fuYTUyPU1LH9YgD7XfVcmfzB8CAwEAAQJAdMC/HlAqjPqNnRHI/Pim
s/UamajYRUkqdx9V3U6Z/byRCxJNI/el/D6swo4nZE25fnj6eZv/LCFhFkCw+g0U
sQIhAMbMLiCLMxGOSRHPsiSRtsYApdj4AB4Gs/QMVVSyVs1lAiEA5HpQe9gTJEDi
a2gLphBjhc93Pe9eEbTTSDJwSCQ/zTMCIQC+MejL0AG7ASNdfBWWsSZpx4Lk01jh
YU5X5ljZYIp1lQIhALIZHk/LWPBzm4t56Uqjj9CorhybUEqhF+k5WAkEKK+9AiEA
xTSSpbskbeY1FbD4MD8TGz6YpWPkTdgGYo2ZE+f3cr8=
-----END RSA_PRIVATE_KEY-----
`

	marshaler := RSAMarshaler{}
	keyPair, err := marshaler.Unmarshal([]byte(encodedPrivateKey))
	if err != nil {
		t.Fatal(err)
	}

	got, err := keyPair.EncodedPublicKey()
	if err != nil {
		t.Error(err)
	}
	if got != encodedPublicKey {
		t.Errorf("expected:\n%s\ngot:\n%s\n", encodedPublicKey, got)
	}
}
