package instanceaddress

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

//	InstanceAddress is a referenced implementation for cryptoScheme PQRINGCT1.0.
//	In PQRINGCT1.0, When forming a TXO, the coin-address directly uses cryptoAddress, rather than using an address derived from cryptoAddress.
//	and it has the following two features:
//	1) the coin-address is a part of the instance-address, and
//	2) the extracting algorithm is deterministic, and the coin-address and instance-address is one-to-one map.
//	In this referenced implementation, 1 byte is used to specify netId.
//	0x00: MainNet; 0x01:RegressionNet; 0x02:TestNet3; 0x03:SimNet

type InstanceAddress struct {
	netID         byte
	cryptoAddress []byte
}

func (instAddr *InstanceAddress) SerializeSize() uint32 {
	//return 1 + abecrypto.GetCryptoAddressSerializeSize(instAdd.cryptoScheme)
	return uint32(1 + len(instAddr.cryptoAddress))
}

func (instAddr *InstanceAddress) Serialize() []byte {
	b := make([]byte, instAddr.SerializeSize())

	b[0] = instAddr.netID
	copy(b[1:], instAddr.cryptoAddress[:])

	return b
}

func (instAddr *InstanceAddress) Deserialize(serializedInstAddr []byte) error {
	if len(serializedInstAddr) <= 1 {
		return errors.New("byte length of serializedInstAddr does not match the design")
	}

	instAddr.netID = serializedInstAddr[0]
	instAddr.cryptoAddress = make([]byte, len(serializedInstAddr)-1)
	copy(instAddr.cryptoAddress, serializedInstAddr[1:])

	return nil
}

func (instAddr *InstanceAddress) Encode() string {
	serialized := instAddr.Serialize()
	checkSum := CheckSum(serialized)

	encodeAddrStr := hex.EncodeToString(serialized)
	encodeAddrStr = encodeAddrStr + hex.EncodeToString(checkSum[:])
	return encodeAddrStr
}

func (instAddr *InstanceAddress) Decode(addrStr string) error {
	addrBytes, err := hex.DecodeString(addrStr)
	if err != nil {
		return err
	}
	checksumLen := CheckSumLength()
	if len(addrBytes) <= 1+checksumLen {
		errStr := fmt.Sprintf("abel-address %v has a wrong length", addrStr)
		return errors.New(errStr)
	}

	serializedInstantAddr := addrBytes[:len(addrBytes)-checksumLen]
	checkSum := addrBytes[len(addrBytes)-checksumLen:]
	checkSumComputed := CheckSum(serializedInstantAddr)
	if bytes.Compare(checkSum, checkSumComputed[:]) != 0 {
		errStr := fmt.Sprintf("abel-address %v has a wrong check sum", addrStr)
		return errors.New(errStr)
	}

	err = instAddr.Deserialize(serializedInstantAddr)
	if err != nil {
		return err
	}

	return nil
}

func (instAddr *InstanceAddress) String() string {
	serialized := instAddr.Serialize()
	return hex.EncodeToString(serialized)
}

func (instAddr *InstanceAddress) CryptoAddress() []byte {
	return instAddr.cryptoAddress
}

func (instAddr *InstanceAddress) IsForNet(netId byte) bool {
	return instAddr.netID == netId
}

func NewInstanceAddress(netId byte, cryptoAddress []byte) *InstanceAddress {
	return &InstanceAddress{netID: netId, cryptoAddress: cryptoAddress}
}

// to be compatible with the previous implementation of abewallet,
// here use DoubleHashH in abec.
func CheckSum(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:]
}
func CheckSumLength() int {
	return 32
}
