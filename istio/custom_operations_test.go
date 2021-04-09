// Package istio - Common operations for the adapter
package istio

import (
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
)

func TestIstio_applyCustomOperation(t *testing.T) {
	type fields struct {
		Adapter adapter.Adapter
	}
	type args struct {
		namespace string
		manifest  string
		isDel     bool
	}

	ch := make(chan interface{}, 10)
	fs := fields{
		Adapter: adapter.Adapter{
			Config:            getConfigHandler(t),
			Log:               getLoggerHandler(t),
			KubeconfigHandler: getKubeconfigHandler(t),
			Channel:           &ch,
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "no manifest or empty manifest",
			fields: fs,
			args: args{
				namespace: "default",
				manifest:  "",
				isDel:     false,
			},
			want:    status.Starting,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			istio := &Istio{
				Adapter: tt.fields.Adapter,
			}
			got, err := istio.applyCustomOperation(tt.args.namespace, tt.args.manifest, tt.args.isDel)
			if (err != nil) != tt.wantErr {
				t.Errorf("Istio.applyCustomOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Istio.applyCustomOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
