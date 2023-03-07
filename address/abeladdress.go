package address

// AbelAddress is the address facing to the users.
// In particular, its EncodeString form is used by users to receive coins.
// AbelAddress vs. CryptoAddress:
// 1. CryptoAddress is generated and used by the underlying crypto scheme, which is provided by sdkapi in abec.
// 2. AbelAddress's string form is used by users to receives coins, such as mining-address or payment-address.
// As a result, AbelAddress instance shall contain (netId, cryptoAddress) and its encodeString form shall contain checksum.
type AbelAddress interface {
	SerializeSize() uint32

	// Serialize() returns the bytes for (netId, cryptoAddress)
	Serialize() []byte

	// Deserialize() parses the bytes for (netId, cryptoAddress) to AbelAddress object.
	Deserialize(serialized []byte) error

	//	Encode() returns hex-codes of (netId, cryptoAddress, checksum) where checksum is the hash of (netId, cryptoAddress)
	Encode() string

	//	Decode() parses the hex-codes of (netId, cryptoAddress, checksum) to AbelAddress object.
	Decode(addrStr string) error

	//	String() returns the hex-codes of (netId, cryptoAddress)
	String() string

	//	CryptoAddress() returns the cryptoAddress
	CryptoAddress() []byte

	IsForNet(netId byte) bool
}
