package istio

import (
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
)

func TestIstio_installAddon(t *testing.T) {
	type fields struct {
		Adapter adapter.Adapter
	}
	type args struct {
		namespace string
		del       bool
		service   string
		patches   []string
		templates []adapter.Template
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
			name:   "no patches",
			fields: fs,
			args: args{
				namespace: "default",
				del:       false,
				service:   "test",
				patches:   nil,
				templates: []adapter.Template{
					"https://raw.githubusercontent.com/istio/istio/master/samples/addons/jaeger.yaml",
				},
			},
			want:    status.Installing,
			wantErr: true,
		},
		{
			name:   "no templates",
			fields: fs,
			args: args{
				namespace: "default",
				del:       false,
				service:   "test",
				patches:   nil,
				templates: nil,
			},
			want:    status.Installed,
			wantErr: false,
		},
		{
			name:   "delete operation",
			fields: fs,
			args: args{
				namespace: "default",
				del:       true,
				service:   "test",
				patches:   nil,
				templates: nil,
			},
			want:    status.Installed,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			istio := &Istio{
				Adapter: tt.fields.Adapter,
			}
			got, err := istio.installAddon(tt.args.namespace, tt.args.del, tt.args.service, tt.args.patches, tt.args.templates)
			if (err != nil) != tt.wantErr {
				t.Errorf("Istio.installAddon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Istio.installAddon() = %v, want %v", got, tt.want)
			}
		})
	}
}
