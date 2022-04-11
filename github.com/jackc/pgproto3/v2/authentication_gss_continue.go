package pgproto3

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/jackc/pgio"
)

type AuthenticationGSSContinue struct {
	Token []byte
}

func (a *AuthenticationGSSContinue) Backend() {}

func (a *AuthenticationGSSContinue) AuthenticationResponse() {}

func (a *AuthenticationGSSContinue) Decode(src []byte) error {
	if len(src) < 4 {
		return errors.New("authentication message too short")
	}

	authType := binary.BigEndian.Uint32(src)

	if authType != AuthTypeGSSCont {
		return errors.New("bad auth type")
	}

	copy(a.Token, src)
	return nil
}

func (a *AuthenticationGSSContinue) Encode(dst []byte) []byte {
	dst = append(dst, 'R')
	dst = pgio.AppendInt32(dst, int32(len(a.Token))+8)
	dst = pgio.AppendUint32(dst, AuthTypeGSSCont)
	dst = append(dst, a.Token...)
	return dst
}

func (a *AuthenticationGSSContinue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Token []byte
	}{
		Type:  "AuthenticationGSSContinue",
		Token: a.Token,
	})
}

func (a *AuthenticationGSSContinue) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	var msg struct {
		Type  string
		Token []byte
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	a.Token = msg.Token
	return nil
}
