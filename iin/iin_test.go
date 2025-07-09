package iin

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		iin         string
		wantValid   bool
		wantSex     string
		wantDate    string
		wantCentury int
		wantError   bool
	}{
		{
			name:        "Valid IIN - Male 2003",
			iin:         "031231500126",
			wantValid:   true,
			wantSex:     "male",
			wantDate:    "31.12.2003",
			wantCentury: 21, // XXI век (2000-2099)
			wantError:   false,
		},
		{
			name:      "Valid IIN - Female 1970",
			iin:       "700701400456",
			wantValid: false, // Будет false так как контрольная сумма не совпадает
			wantError: true,
		},
		{
			name:      "Invalid length - too short",
			iin:       "12345678901",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid length - too long",
			iin:       "1234567890123",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid characters",
			iin:       "12345678901a",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid month",
			iin:       "031331500123",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid day",
			iin:       "033231500123",
			wantValid: false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := Validate(tt.iin)

			if tt.wantError {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
				return
			}

			if info.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", info.Valid, tt.wantValid)
			}

			if tt.wantValid {
				if info.Sex != tt.wantSex {
					t.Errorf("Validate() sex = %v, want %v", info.Sex, tt.wantSex)
				}

				if info.DateOfBirth != tt.wantDate {
					t.Errorf("Validate() dateOfBirth = %v, want %v", info.DateOfBirth, tt.wantDate)
				}

				if info.Century != tt.wantCentury {
					t.Errorf("Validate() century = %v, want %v", info.Century, tt.wantCentury)
				}
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		iin  string
		want bool
	}{
		{"Valid IIN", "031231500126", true},
		{"Invalid length", "12345678901", false},
		{"Invalid characters", "12345678901a", false},
		{"Invalid checksum", "123456789012", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.iin); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSex(t *testing.T) {
	tests := []struct {
		name    string
		iin     string
		want    string
		wantErr bool
	}{
		{"Male - century 20", "031231500126", "male", false},
		{"Female - century 20", "031231600123", "female", false},
		{"Male - century 19", "701231300123", "male", false},
		{"Female - century 19", "701231400123", "female", false},
		{"Invalid length", "12345678901", "", true},
		{"Invalid sex digit", "031231700123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractSex(tt.iin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractSex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractSex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractDateOfBirth(t *testing.T) {
	tests := []struct {
		name    string
		iin     string
		want    string
		wantErr bool
	}{
		{"Valid date 2003", "031231500126", "31.12.2003", false},
		{"Valid date 1970", "701231300123", "31.12.1970", false},
		{"Invalid month", "031331500123", "", true},
		{"Invalid day", "033231500123", "", true},
		{"Invalid length", "12345678901", "", true},
		{"February 29 leap year", "000229500123", "29.02.2000", false},
		{"February 29 non-leap year", "010229500123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractDateOfBirth(tt.iin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractDateOfBirth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractDateOfBirth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAndExtract(t *testing.T) {
	tests := []struct {
		name      string
		iin       string
		wantValid bool
		wantSex   string
		wantDate  string
		wantErr   bool
	}{
		{
			name:      "Valid IIN",
			iin:       "031231500126",
			wantValid: true,
			wantSex:   "male",
			wantDate:  "31.12.2003",
			wantErr:   false,
		},
		{
			name:      "Invalid IIN",
			iin:       "123456789012",
			wantValid: false,
			wantSex:   "",
			wantDate:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotSex, gotDate, err := ValidateAndExtract(tt.iin)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndExtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotValid != tt.wantValid {
				t.Errorf("ValidateAndExtract() valid = %v, want %v", gotValid, tt.wantValid)
			}

			if gotSex != tt.wantSex {
				t.Errorf("ValidateAndExtract() sex = %v, want %v", gotSex, tt.wantSex)
			}

			if gotDate != tt.wantDate {
				t.Errorf("ValidateAndExtract() date = %v, want %v", gotDate, tt.wantDate)
			}
		})
	}
}

// Бенчмарки для проверки производительности
func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("031231500126")
	}
}

func BenchmarkIsValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsValid("031231500126")
	}
}

func BenchmarkExtractSex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractSex("031231500126")
	}
}

func BenchmarkExtractDateOfBirth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractDateOfBirth("031231500126")
	}
}

func BenchmarkValidateAndExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateAndExtract("031231500126")
	}
}
