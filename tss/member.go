package tss

import (
	"github.com/xuyp1991/thresholdSign/operator"
	"encoding/hex"
	"bytes"
	"math/big"
)

type MemberID []byte

// MemberIDFromPublicKey creates a MemberID from a public key.
func MemberIDFromPublicKey(publicKey *operator.PublicKey) MemberID {
	return operator.Marshal(publicKey)
}

// PublicKey returns the MemberID as a public key.
func (id MemberID) PublicKey() (*operator.PublicKey, error) {
	return operator.Unmarshal(id)
}

// MemberIDFromPublicKey creates a MemberID from a string.
func MemberIDFromString(string string) (MemberID, error) {
	return hex.DecodeString(string)
}

// String converts MemberID to string.
func (id MemberID) String() string {
	return hex.EncodeToString(id)
}

// bigInt converts MemberID to big.Int.
func (id MemberID) bigInt() *big.Int {
	return new(big.Int).SetBytes(id)
}

// Equal checks if member IDs are equal.
func (id MemberID) Equal(memberID MemberID) bool {
	return bytes.Equal(id, memberID)
}

// groupInfo holds information about the group selected for protocol execution.
type GroupInfo struct {
	GroupID        string // globally unique group identifier
	MemberID       MemberID//暂时先不用.
	GroupMemberIDs []MemberID
	DishonestThreshold int
}