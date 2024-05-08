package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestConvertStringToINT(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int
		wantErr bool
	}{
		{
			name:    "Success",
			args:    "1",
			want:    1,
			wantErr: false,
		},
		{
			name:    "Failed Convert",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ConvertStringToINT(tt.args); (err != nil) != tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertStringToINT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	timeNow := GetTime()

	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "Success",
			want: timeNow,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeNow; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.GetTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd30Day(t *testing.T) {
	add30Day := Add30Day()

	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "Success",
			want: add30Day,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := add30Day; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("utils.add30Day() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateNewUniqueNumber(t *testing.T) {
	t.Parallel()

	got := GenerateNewUniqueNumber()

	if len(got) != 32 {
		t.Errorf("GenerateNewReferenceNumber() got = %v, wantLen %v", got, 32)

		return
	}
}
