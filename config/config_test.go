package config_test

import (
	"testing"

	. "bulk/config"
)

func TestConfigValidate(t *testing.T) { //nolint:tparallel
	t.Parallel()

	type fields struct {
		ENV         string
		LogLevel    string
		GRPCAddress string
		Tracer      Tracer
		GCS         GCS
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Correct All Fields Required already assigned", fields: fields{
				ENV:         "development",
				LogLevel:    "test-logger",
				GRPCAddress: "9000",
				Tracer:      Tracer{TracerName: "", Enable: false},
				GCS: GCS{
					Bucket: "test-bucket",
				},
			}, wantErr: false,
		},
		{
			name: "ENV was missing", fields: fields{
				ENV:         "",
				LogLevel:    "test-logger",
				GRPCAddress: "9000",
				Tracer:      Tracer{TracerName: "", Enable: false},
			}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				cfg := Config{
					ENV:         tt.fields.ENV,
					LogLevel:    tt.fields.LogLevel,
					GRPCAddress: tt.fields.GRPCAddress,
					Tracer:      tt.fields.Tracer,
					GCSBulk:     tt.fields.GCS,
				}
				if err := cfg.Validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
