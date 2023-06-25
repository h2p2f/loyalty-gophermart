package jwt

import "testing"

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name  string
		login string
		want  string
	}{
		{
			name:  "Positive test 1",
			login: "test1",
			want:  "test1",
		},
		{
			name:  "Positive test 2",
			login: "test2",
			want:  "test2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.login, "secretest key")
			if err != nil {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
			login, err := ParseToken(got, "secretest key")
			if err != nil {
				t.Errorf("ParseToken() = %v, want %v", login, tt.want)
			}
			if login != tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
