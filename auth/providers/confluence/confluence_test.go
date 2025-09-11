package confluence

import (
	"testing"
	"time"
)

func TestRefreshToken(t *testing.T) {
	Init()
	token, err := p.RefreshToken("ey")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("new token: %+v", token)
}

func TestProvider_FetchUser(t *testing.T) {
	Init()
	session := &Session{
		AuthURL:      "",
		AccessToken:  "eyJ",
		RefreshToken: "",
		ExpiresAt:    time.Time{},
		ExpiresIn:    0,
	}
	u, err := p.FetchUser(session)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user: %+v", u)
}
