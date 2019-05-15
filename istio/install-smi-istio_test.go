package istio

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_getSMIYamls(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getSMIYamls()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSMIYamls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if got != tt.want {
			// 	t.Errorf("getSMIYamls() = %v, want %v", got, tt.want)
			// }
		})
	}
}
