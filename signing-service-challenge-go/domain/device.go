package domain

type SignatureAlgorithm string

const (
	RSS SignatureAlgorithm = "RSS"
	ECC SignatureAlgorithm = "ECC"
)

type SignatureDevice struct {
	uuid              string
	algorithm         SignatureAlgorithm
	encodedPrivateKey []byte
	// (optional) user provided string to be displayed in the UI
	label string
	// track the last signature created with this device
	lastSignature string
	// track how many signatures have been created with this device
	signatureCounter uint
}
