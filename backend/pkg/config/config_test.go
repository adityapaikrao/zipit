package config

import (
	"os"
	"testing"
)

func TestNewDBConfig(t *testing.T) {
	// 1. Define your test cases using a "Table Driven Tests" pattern
	tests := []struct {
		name     string
		envVars  map[string]string
		wantErr  bool
		validate func(*testing.T, *DBConfig)
	}{
		{
			name: "Valid settings with custom values",
			envVars: map[string]string{
				"DB_USER": "testuser",
				"DB_NAME": "testdb",
				"DB_PORT": "5433",
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *DBConfig) {
				if cfg.User != "testuser" {
					t.Errorf("expected User testuser, got %s", cfg.User)
				}
				if cfg.Port != 5433 {
					t.Errorf("expected Port 5433, got %d", cfg.Port)
				}
				if cfg.DbName != "testdb" {
					t.Errorf("expected DbName testdb, got %s", cfg.DbName)
				}
			},
		},
		{
			name: "Missing required DB_USER",
			envVars: map[string]string{
				"DB_NAME": "testdb",
			},
			wantErr: true,
		},
		{
			name: "Invalid port number",
			envVars: map[string]string{
				"DB_USER": "u",
				"DB_NAME": "d",
				"DB_PORT": "not_a_number",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear Env and set test variables
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// 2. Call the function you are testing
			cfg, err := NewDBConfig()

			// 3. Logic to check if tt.wantErr matches the actual error
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDBConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 4. If there is a validate function and no error, call it
			if !tt.wantErr && err == nil && tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// TODO: Implement a simple test for the getEnvOrDefault helper
	tests := []struct {
		name        string
		envVals     [2]string
		defaultVal  string
		expectedVal string
	}{
		{
			name:        "Returns default when value is missing",
			envVals:     [2]string{"SOME_KEY", ""},
			defaultVal:  "default_val",
			expectedVal: "default_val",
		},
		{
			name:        "Returns value when value is present",
			envVals:     [2]string{"EXISTING_KEY", "existing_val"},
			defaultVal:  "default_val",
			expectedVal: "existing_val",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.envVals[1] != "" {
				os.Setenv(tt.envVals[0], tt.envVals[1])
			}
			val := getEnvOrDefault(tt.envVals[0], tt.defaultVal)

			if tt.expectedVal != val {
				t.Errorf("getEnvOrDefault() error, expectedVal = %v actualValue = %v", tt.expectedVal, val)
			}

		})
	}
}
