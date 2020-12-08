package istio

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	adapterconfig "github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/logger"
)

// Istio represents the istio adapter and embeds adapter.Adapter
type Istio struct {
	adapter.Adapter // Type Embedded
}

// New initializes istio handler.
func New(c adapterconfig.Handler, l logger.Handler, kc adapterconfig.Handler) adapter.Handler {
	return &Istio{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
		},
	}
}

// ApplyOperation applies the operation on istio
func (istio *Istio) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {
	operations := make(adapter.Operations)
	err := istio.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	e := &adapter.Event{
		Operationid: opReq.OperationID,
		Summary:     status.Deploying,
		Details:     "Operation is not supported",
	}

	switch opReq.OperationName {
	case internalconfig.IstioOperation:
		go func(hh *Istio, ee *adapter.Event) {
			version := string(operations[opReq.OperationName].Versions[0])
			stat, err := hh.installIstio(opReq.IsDeleteOperation, version, opReq.Namespace)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s Istio service mesh", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Istio service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("The Istio service mesh is now %s.", stat)
			hh.StreamInfo(e)
		}(istio, e)
	case common.BookInfoOperation, common.HTTPBinOperation, common.ImageHubOperation, common.EmojiVotoOperation:
		go func(hh *Istio, ee *adapter.Event) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(opReq.Namespace, opReq.IsDeleteOperation, operations[opReq.OperationName].Templates)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s %s application", stat, appName)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(e)
		}(istio, e)
	case common.SmiConformanceOperation:
		go func(hh *Istio, ee *adapter.Event) {
			name := operations[opReq.OperationName].Description
			err := hh.ValidateSMIConformance(&adapter.SmiTestOptions{
				Ctx:  context.TODO(),
				OpID: ee.Operationid,
			})
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s %s test", status.Running, name)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s test %s successfully", name, status.Completed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(istio, e)
	case common.CustomOperation:
		go func(hh *Istio, ee *adapter.Event) {
			stat, err := hh.applyCustomOperation(opReq.Namespace, opReq.CustomBody, opReq.IsDeleteOperation)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s custom operation", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Manifest %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(istio, e)
	case internalconfig.LabelNamespace:
		go func(hh *Istio, ee *adapter.Event) {
			err := hh.LoadNamespaceToMesh(opReq.Namespace, opReq.IsDeleteOperation)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while labelling %s", opReq.Namespace)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = "Labelling successful"
			ee.Details = ""
			hh.StreamInfo(e)
		}(istio, e)
	case internalconfig.PrometheusAddon, internalconfig.GrafanaAddon, internalconfig.KialiAddon, internalconfig.JaegerAddon, internalconfig.ZipkinAddon:
		go func(hh *Istio, ee *adapter.Event) {
			_, err := hh.InstallAddon(opReq.Namespace, opReq.IsDeleteOperation, opReq.OperationName)
			operation := "install"
			if opReq.IsDeleteOperation {
				operation = "uninstall"
			}

			if err != nil {
				e.Summary = fmt.Sprintf("Error while %sing %s", operation, opReq.OperationName)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Succesfully %sed %s", operation, opReq.OperationName)
			ee.Details = fmt.Sprintf("Succesfully %sed %s from the %s namespace", operation, opReq.OperationName, opReq.Namespace)
			hh.StreamInfo(e)
		}(istio, e)
	case internalconfig.IstioVetOpertation:
		go func(hh *Istio, ee *adapter.Event) {
			responseChan := make(chan *adapter.Event, 1)

			go hh.RunVet(responseChan)

			for msg := range responseChan {
				switch msg.EType {
				case int32(meshes.EventType_ERROR):
					istio.StreamErr(msg, ErrIstioVet(fmt.Errorf(msg.Details)))
				case int32(meshes.EventType_WARN):
					istio.StreamWarn(msg, ErrIstioVet(fmt.Errorf(msg.Details)))
				default:
					istio.StreamInfo(msg)
				}
			}

			istio.Log.Info("Done")
		}(istio, e)
	default:
		istio.StreamErr(e, ErrOpInvalid)
	}

	return nil
}
