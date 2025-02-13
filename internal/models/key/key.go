package key

import (
	"encoding"
	"encoding/hex"
)

type HexString []byte

var _ encoding.TextUnmarshaler = (*HexString)(nil)

// UnmarshalText implements encoding.TextUnmarshaler.
func (h *HexString) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}

	*h = b
	return nil
}

func (h HexString) String() string {
	return hex.EncodeToString(h)
}
