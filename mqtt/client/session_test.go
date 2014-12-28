package client

import "testing"

func Test_newSession(t *testing.T) {
	cleanSession := true
	clientIDStr := "clientID"

	sess := newSession(cleanSession, []byte(clientIDStr))

	if sess.cleanSession != cleanSession {
		t.Errorf("sess.cleanSession => %t, want => %t", sess.cleanSession, cleanSession)
		return
	}

	if string(sess.clientID) != clientIDStr {
		t.Errorf("string(sess.clientID) => %s, want => %s", string(sess.clientID), clientIDStr)
	}
}
