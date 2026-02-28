package cmd

import (
	"testing"
)

func TestCreditsCommand_Output(t *testing.T) {
	// Test that the command exists and has correct structure
	if creditsCmd.Use != "credits" {
		t.Errorf("expected Use 'credits', got %s", creditsCmd.Use)
	}
}
