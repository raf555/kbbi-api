package encoding

import (
	"encoding"
	"encoding/hex"
	"fmt"
)

type HexString []byte

var _ encoding.TextUnmarshaler = (*HexString)(nil)

// UnmarshalText implements encoding.TextUnmarshaler.
func (h *HexString) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return fmt.Errorf("hex.DecodeString: %w", err)
	}

	*h = b
	return nil
}

func (h HexString) String() string {
	return fmt.Sprintf("HexString(%s)", hex.EncodeToString(h))
}
