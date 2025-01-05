package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	pass1 := "password"
	pass2 := "password2"
	hash1, _ := HashPassword(pass1)
	hash2, _ := HashPassword(pass2)
	tests := []struct {
		name string
		password string
		hash string
		wantErr bool
	}{
		{
			name: "Password is correct",
			password: pass1,
			hash: hash1,
			wantErr: false,
		},
		{
			name: "Password is incorrect",
			password: pass2,
			hash: hash1,
			wantErr: true,
		},
		{
			name: "Different hashes",
			password: pass1,
			hash: hash2,
			wantErr: true,
		},
		{
			name: "Empty password",	
			password: "",
			hash: hash1,
			wantErr: true,
		},
		{
			name: "Invalid hash",
			password: pass1,
			hash: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPasswordHash(tt.password, tt.hash); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}