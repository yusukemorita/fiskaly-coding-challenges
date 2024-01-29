package crypto

import (
	"crypto/ecdsa"
	"testing"
)

func TestECCKeyPair_Sign(t *testing.T) {
	generator := ECCGenerator{}
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

	result := ecdsa.VerifyASN1(keyPair.Public, digest, signature)
	if !result {
		t.Errorf("signature verification failed: %s", err)
	}
}

func TestECCKeyPair_EncodedPublicKey(t *testing.T) {
	// encodedPublicKey, encodedPrivateKey are a key pair pre-generated
	// for this test
	encodedPublicKey := `-----BEGIN PUBLIC_KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE3ZA9cyJ1LUM8APSX8So+Id8fx0PI+u8s
0CP1qUr5FxCwAzYuCauch5k7zS6ikiChiRxlXYE89drr55OfiEHflSq3XXTX6evj
I/dxHZ28t7rbetKSxU64GxuXdT6JytqX
-----END PUBLIC_KEY-----
`
	encodedPrivateKey := `-----BEGIN PRIVATE_KEY-----
MIGkAgEBBDDu03JFZzy/SxN5jOvnoFwiecUUE+eMn43EgUhIcJUhF03gNtBZxhNI
bFDiORLoSX2gBwYFK4EEACKhZANiAATdkD1zInUtQzwA9JfxKj4h3x/HQ8j67yzQ
I/WpSvkXELADNi4Jq5yHmTvNLqKSIKGJHGVdgTz12uvnk5+IQd+VKrdddNfp6+Mj
93Ednby3utt60pLFTrgbG5d1PonK2pc=
-----END PRIVATE_KEY-----
`

	marshaler := ECCMarshaler{}
	keyPair, err := marshaler.Decode([]byte(encodedPrivateKey))
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
