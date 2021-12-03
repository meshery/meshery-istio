package istio

import (
	"context"
	"reflect"
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	adapterconfig "github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	internalconfig "github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/logger"
)

func TestNew(t *testing.T) {
	type args struct {
		c  adapterconfig.Handler
		l  logger.Handler
		kc adapterconfig.Handler
	}

	type test struct {
		name string
		args args
		want adapter.Handler
	}

	tests := []test{
		{
			name: "no arguments",
			args: args{
				c:  nil,
				l:  nil,
				kc: nil,
			},
			want: &Istio{
				Adapter: adapter.Adapter{},
			},
		},
		// Add more cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.c, tt.args.l, tt.args.kc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v\n, want %+v\n", got, tt.want)
			}
		})
	}
}

func TestIstio_ApplyOperation(t *testing.T) {
	type fields struct {
		Adapter adapter.Adapter
	}
	type args struct {
		ctx   context.Context
		opReq adapter.OperationRequest
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
		wantErr bool
	}{
		//{
		//	name: "Unseeded config operation",
		//	fields: fields{
		//		Adapter: adapter.Adapter{
		//			Log:     getLoggerHandler(t),
		//			Config:  getConfigHandlerUnseeded(t),
		//			Channel: &ch,
		//		},
		//	},
		//	args: args{
		//		ctx: context.TODO(),
		//		opReq: adapter.OperationRequest{
		//			OperationName:     "stale",
		//			Namespace:         "default",
		//			IsDeleteOperation: false,
		//			OperationID:       "test_id",
		//		},
		//	},
		//	wantErr: false,
		//},
		// Tests for stale operation
		{
			name:   "Stale operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     "stale",
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		// Tests for istio operation
		{
			name:   "Istio operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.IstioOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		// Tests for sample apps operation
		{
			name:   "BookInfo operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     common.BookInfoOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "HTTPBin operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     common.HTTPBinOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "ImageHub operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     common.ImageHubOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "EmojiVoto operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     common.EmojiVotoOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		// Tests for validate operation
		{
			name:   "SMI operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     common.SmiConformanceOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		//{
		//	name:   "Istio Vet operation",
		//	fields: fs,
		//	args: args{
		//		ctx: context.TODO(),
		//		opReq: adapter.OperationRequest{
		//			OperationName:     internalconfig.IstioVetOperation,
		//			Namespace:         "default",
		//			IsDeleteOperation: false,
		//			OperationID:       "test_id",
		//		},
		//	},
		//	wantErr: false,
		//},
		// Tests for configure operation
		{
			name:   "Deny All Policy operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.DenyAllPolicyOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Strict MTLS Policy operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.StrictMTLSPolicyOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Mutual MTLS Policy operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.MutualMTLSPolicyOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Disable MTLS Policy operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.DisableMTLSPolicyOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Label Namespace operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.LabelNamespace,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Envoy Filter operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.EnvoyFilterOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		// Tests for custom operation
		{
			name:   "Custom operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     common.CustomOperation,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		// Tests for addon operation
		{
			name:   "Prometheus Addon operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.PrometheusAddon,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Grafana Addon operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.GrafanaAddon,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Kiali Addon operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.KialiAddon,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Jaeger Addon operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.JaegerAddon,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		{
			name:   "Zipkin Addon operation",
			fields: fs,
			args: args{
				ctx: context.TODO(),
				opReq: adapter.OperationRequest{
					OperationName:     internalconfig.ZipkinAddon,
					Namespace:         "default",
					IsDeleteOperation: false,
					OperationID:       "test_id",
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.Adapter.ApplyOperation(tt.args.ctx, tt.args.opReq); (err != nil) != tt.wantErr {
				t.Errorf("Istio.ApplyOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIstio_ProcessOAM(t *testing.T) {
	type fields struct {
		Adapter adapter.Adapter
	}
	type args struct {
		ctx    context.Context
		oamReq adapter.OAMRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			istio := &Istio{
				Adapter: tt.fields.Adapter,
			}
			got, err := istio.ProcessOAM(tt.args.ctx, tt.args.oamReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Istio.ProcessOAM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Istio.ProcessOAM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getConfigHandler(t *testing.T) adapterconfig.Handler {
	h, _ := internalconfig.New(configprovider.ViperKey)
	return h
}

//func getConfigHandlerUnseeded(t *testing.T) adapterconfig.Handler {
//	h, _ := meshkitprovider.NewViper(meshkitprovider.Options{
//		FileName: "istio",
//		FileType: "yaml",
//		FilePath: path.Join(utils.GetHome(), ".meshery"),
//	})
//	return h
//}

func getLoggerHandler(t *testing.T) logger.Handler {
	log, _ := logger.New("istio test", logger.Options{
		Format:     logger.SyslogLogFormat,
		DebugLevel: true,
	})
	return log
}

func getKubeconfigHandler(t *testing.T) adapterconfig.Handler {
	h, _ := internalconfig.NewKubeconfigBuilder(configprovider.ViperKey)
	return h
}
