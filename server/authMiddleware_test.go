package server

import (
	"testing"
)

func TestAuthorityCheck(t *testing.T) {

	//checks Admin may play User role
	if !hasAuthority("ROLE_USER", []string{"ROLE_ADMINISTRATOR"}) {
		t.Error("Incorrect user role validation")
	}

	//checks User cannot play admin role
	if hasAuthority("ROLE_ADMINISTRATOR", []string{"ROLE_USER"}) {
		t.Error("Incorrect user role validation")
	}

	//checks nobody can play unknown role
	if hasAuthority("unknown", []string{"ROLE_USER"}) {
		t.Error("Incorrect unknown user role validation")
	}

	//checks nobody can play unknown role
	if hasAuthority("ROLE_USER", []string{"unknown"}) {
		t.Error("Incorrect unknown user role validation")
	}

}
