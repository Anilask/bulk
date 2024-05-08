package logger

import (
	"bulk/constant"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		service  string
		env      string
		logLevel string
	}
	mockLogger := Init("test", "test", "error")
	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{
			name: "OK",
			args: args{
				service:  "test",
				env:      "test",
				logLevel: "",
			},
			want: mockLogger,
		},
		{
			name: "Wrong log level",
			args: args{
				service:  "test",
				env:      "test",
				logLevel: "-1",
			},
			want: mockLogger,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := Init(tt.args.service, tt.args.env, tt.args.logLevel)
				if got.env != tt.args.env {
					t.Errorf("Init() = %v, want %v", got, tt.want)
				}
				if got.service != tt.args.service {
					t.Errorf("Init() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestLogger_Debug_Dev(t *testing.T) {
	mockLogger := Init("test", constant.EnvDev, "error")
	mockLogger.Debug("test")
}

func TestLogger_Debug_Stage(t *testing.T) {
	mockLogger := Init("test", constant.EnvStage, "error")
	mockLogger.Debug("test")
}

func TestLogger_Debug_Prod(t *testing.T) {
	mockLogger := Init("test", "prod", "error")
	mockLogger.Debug("test")
}

func TestLogger_Debugf_Dev(t *testing.T) {
	mockLogger := Init("test", constant.EnvDev, "error")
	mockLogger.Debugf("%s", "test")
}

func TestLogger_Debugf_Stage(t *testing.T) {
	mockLogger := Init("test", constant.EnvStage, "error")
	mockLogger.Debugf("%s", "test")
}

func TestLogger_Debugf_Prod(t *testing.T) {
	mockLogger := Init("test", "prod", "debug")
	mockLogger.Debugf("%s", "test")
}

func TestLogger_Error(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Error("test")
}

func TestLogger_Errorf(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Errorf("%s", "test")
}

func TestLogger_Info(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Info("test")
}

func TestLogger_Infof(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Infof("%s", "test")
}

func TestLogger_Infow(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Infow(map[string]interface{}{}, "test")
}

func TestLogger_Warn(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Warn("test")
}

func TestLogger_Warnf(t *testing.T) {
	mockLogger := Init("test", "test", "error")
	mockLogger.Warnf("%s", "test")
}

func TestNew(t *testing.T) {
	type args struct {
		serviceName string
		env         string
		logLevel    string
	}
	mockLogger := New("test", "test", "debug")
	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{
			name: "OK",
			args: args{
				serviceName: "test",
				env:         "test",
				logLevel:    "debug",
			},
			want: &Logger{
				service: "test",
				env:     "test",
				log:     mockLogger.log,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := Init(tt.args.serviceName, tt.args.env, tt.args.logLevel)
				if got.env != tt.args.env {
					t.Errorf("New() = %v, want %v", got, tt.want)
				}
				if got.service != tt.args.serviceName {
					t.Errorf("New() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
