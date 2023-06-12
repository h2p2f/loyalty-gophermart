package luhn

import "testing"

func Test_Validate(t *testing.T) {
	tests := []struct {
		name  string
		order string
		want  bool
	}{
		{
			name:  "Positive test",
			order: "1234567890123452",
			want:  true,
		},
		{
			name:  "Negative test",
			order: "1234567890123456",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.order)

			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
