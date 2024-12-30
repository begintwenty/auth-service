package domain_test

import (
	"testing"

	"github.com/dev-mantas/authservice/domain"
)

func TestNewUserID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want domain.UserID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.NewUserID()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("NewUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissions_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var p domain.Permissions
			got, gotErr := p.MarshalJSON()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MarshalJSON() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MarshalJSON() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
