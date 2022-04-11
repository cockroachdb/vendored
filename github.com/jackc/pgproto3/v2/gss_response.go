package pgproto3

import (
	"encoding/json"
	"github.com/jackc/pgio"
)

type GSSResponse struct {
	Token []byte
}

// Frontend identifies this message as sendable by a PostgreSQL frontend.
func (g *GSSResponse) Frontend() {}

func (g *GSSResponse) Decode(data []byte) error {
	g.Token = data
	return nil
}

func (g *GSSResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'p')
	dst = pgio.AppendInt32(dst, int32(4+len(g.Token)))
	dst = append(dst, g.Token...)
	return dst
}

// MarshalJSON implements encoding/json.Marshaler.
func (g *GSSResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Token []byte
	}{
		Type:  "GSSResponse",
		Token: g.Token,
	})
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (g *GSSResponse) UnmarshalJSON(data []byte) error {
	var msg struct {
		Token []byte
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	g.Token = msg.Token
	return nil
}
