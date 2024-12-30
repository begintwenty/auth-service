package authservice_test

import (
	"testing"
)

type UserService struct {
}

func New(authservice *UserService) *UserService {
	return &UserService{}
}

func TestAuthcheck(t *testing.T) {

}

// func TestService_Store(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		authRepo authservice.AuthRepo
// 		// Named input parameters for target function.
// 		auth    domain.Auth
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := authservice.New(tt.authRepo)
// 			gotErr := s.Store(context.Background(), tt.auth)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("Store() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("Store() succeeded unexpectedly")
// 			}
// 		})
// 	}
// }
