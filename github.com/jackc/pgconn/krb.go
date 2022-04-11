package pgconn

import (
	"errors"
	"github.com/jackc/pgproto3/v2"
)

// NewGSSFunc creates a GSS authentication provider, for use with
// RegisterGSSProvider.
type NewGSSFunc func() (GSS, error)

var newGSS NewGSSFunc

// RegisterGSSProvider registers a GSS authentication provider. For example, if
// you need to use Kerberos to authenticate with your server, add this to your
// main package:
//
//	import "github.com/jackc/pgconn/auth/krb"
//
//	func init() {
//		pgconn.RegisterGSSProvider(func() (pgconn.GSS, error) { return krb.NewGSS() })
//	}
func RegisterGSSProvider(newGSSArg NewGSSFunc) {
	newGSS = newGSSArg
}

// GSS provides GSSAPI authentication (e.g., Kerberos).
type GSS interface {
	GetInitToken(host string, service string) ([]byte, error)
	GetInitTokenFromSPN(spn string) ([]byte, error)
	Continue(inToken []byte) (done bool, outToken []byte, err error)
}

func (c *PgConn) gssAuth() error {
	if newGSS == nil {
		return errors.New("kerberos error: no GSSAPI provider registered")
	}
	cli, err := newGSS()
	if err != nil {
		return err
	}

	var nextToken []byte
	if spn, ok := c.config.RuntimeParams["krbspn"]; ok {
		// Use the supplied SPN if provided.
		nextToken, err = cli.GetInitTokenFromSPN(spn)
	} else {
		// Allow the kerberos service name to be overridden
		service := "postgres"
		if val, ok := c.config.RuntimeParams["krbsrvname"]; ok {
			service = val
		}
		nextToken, err = cli.GetInitToken(c.config.Host, service)
	}
	if err != nil {
		return err
	}

	for {
		gssResponse := &pgproto3.GSSResponse{
			Token: nextToken,
		}
		_, err = c.conn.Write(gssResponse.Encode(nil))
		if err != nil {
			return err
		}
		resp, err := c.rxGSSContinue()
		if err != nil {
			return err
		}
		var done bool
		done, nextToken, err = cli.Continue(resp.Token)
		if err != nil {
			return err
		}
		if done {
			break
		}
	}
	return nil
}

func (c *PgConn) rxGSSContinue() (*pgproto3.AuthenticationGSSContinue, error) {
	msg, err := c.receiveMessage()
	if err != nil {
		return nil, err
	}
	gssContinue, ok := msg.(*pgproto3.AuthenticationGSSContinue)
	if ok {
		return gssContinue, nil
	}

	return nil, errors.New("expected AuthenticationGSSContinue message but received unexpected message")
}
