package openpgp

// pgpMessageHeader is the standard armor header for PGP messages.
const pgpMessageHeader = "PGP MESSAGE"

// KeyPair holds the generated armored keys and revocation certificate.
// This is the primary data structure representing a user's PGP identity within the system.
type KeyPair struct {
	PublicKey             string
	PrivateKey            string
	RevocationCertificate string
}
