package seeder

import (
	"testing"
)

func TestSeedIfEmpty_NilDB(t *testing.T) {
	// Initialize Seeder with nil DB
	s := NewSeeder(nil)

	// Call SeedIfEmpty
	err := s.SeedIfEmpty()

	// Assert error is not nil
	if err == nil {
		t.Error("Expected error when seeding with nil DB, got nil")
	}

	// Assert error message
	expectedErr := "seeder database connection is nil"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
	}
}
