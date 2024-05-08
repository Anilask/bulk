package datetime_test

import (
	"testing"

	"bulk/pkg/datetime"
)

func TestGetTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{"Correct time", datetime.GetTime()},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(
			testCase.name, func(t *testing.T) {
				t.Parallel()
				if got := datetime.GetTime(); got != testCase.want {
					t.Errorf("GetTime() = %v, want %v", got, testCase.want)
				}
			},
		)
	}
}
