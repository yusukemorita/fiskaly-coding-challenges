package crypto

import "testing"

func TestFindSupportedAlgorithm(t *testing.T) {
	t.Run("returns found: false when algorithm does not exist", func(t *testing.T) {
		_, found := FindSupportedAlgorithm("INVALID")
		if found {
			t.Error("expected found: false")
		}
	})

	t.Run("returns algorithm when algorithm exists", func(t *testing.T) {
		algorithm, found := FindSupportedAlgorithm("RSA")
		if !found {
			t.Error("expected found: true")
		}

		if algorithm.Name() != "RSA" {
			t.Errorf("expected %s, got %s", "RSA", algorithm.Name())
		}
	})
}