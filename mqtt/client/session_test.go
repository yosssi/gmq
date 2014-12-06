package client

import "testing"

func TestNewSession(t *testing.T) {
	sess := NewSession(nil)

	if sess.CleanSession != DefaultCleanSession {
		t.Errorf("sess.CleanSession => %t, want => %t", sess.CleanSession, DefaultCleanSession)
		return
	}

	if sess.ClientID != hostname {
		t.Errorf("sess.ClientID => %q, want => %q", sess.ClientID, hostname)
		return
	}
}
