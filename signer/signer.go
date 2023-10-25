package signer

import "github.com/Ja7ad/irys/errors"

type Signer interface {
	Sign(data []byte) (signature []byte, err error)
	Verify(data []byte, signature []byte) (err error)
	GetOwner() ([]byte, error)
	GetType() SignatureType
	GetSignatureLength() int
	GetOwnerLength() int
}

// Signer created ONLY FOR VERIFICATION of the signature.
// Private key is not initialized.
func GetSigner(SignatureType SignatureType, owner []byte) (signer Signer, err error) {
	switch SignatureType {
	case SignatureTypeArweave:
		signer = &ArweaveSigner{
			Owner: owner,
		}
	case SignatureTypeEthereum:
		signer = &EthereumSigner{
			Owner: owner,
		}
	default:
		err = errors.ErrUnsupportedSignatureType
	}
	return
}
