package istio

import (
	"fmt"
	"strings"

	"github.com/aspenmesh/istio-vet/pkg/istioclient"
	"github.com/aspenmesh/istio-vet/pkg/vetter"
	"github.com/aspenmesh/istio-vet/pkg/vetter/applabel"
	"github.com/aspenmesh/istio-vet/pkg/vetter/conflictingvirtualservicehost"
	"github.com/aspenmesh/istio-vet/pkg/vetter/danglingroutedestinationhost"
	"github.com/aspenmesh/istio-vet/pkg/vetter/meshversion"
	"github.com/aspenmesh/istio-vet/pkg/vetter/podsinmesh"
	"github.com/aspenmesh/istio-vet/pkg/vetter/serviceassociation"
	"github.com/aspenmesh/istio-vet/pkg/vetter/serviceportprefix"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	istioinformer "istio.io/client-go/pkg/informers/externalversions"
	"k8s.io/client-go/informers"
)

type metaInformerFactory struct {
	k8s   informers.SharedInformerFactory
	istio istioinformer.SharedInformerFactory
}

func (m *metaInformerFactory) K8s() informers.SharedInformerFactory {
	return m.k8s
}
func (m *metaInformerFactory) Istio() istioinformer.SharedInformerFactory {
	return m.istio
}

// RunVet runs istio-vet
func (istio *Istio) RunVet(ch chan<- *adapter.Event) {
	istioClient, err := istioclient.New(&istio.RestConfig)
	if err != nil {
		e := &adapter.Event{}
		e.EType = int32(meshes.EventType_ERROR)
		e.Details = ErrCreatingIstioClient(err).Error()
		e.Summary = "Unable to create istio client"
		ch <- e
	}

	kubeInformerFactory := informers.NewSharedInformerFactory(istio.KubeClient, 0)
	istioInformerFactory := istioinformer.NewSharedInformerFactory(istioClient, 0)
	informerFactory := &metaInformerFactory{
		k8s:   kubeInformerFactory,
		istio: istioInformerFactory,
	}

	vList := []vetter.Vetter{
		vetter.Vetter(podsinmesh.NewVetter(informerFactory)),
		vetter.Vetter(meshversion.NewVetter(informerFactory)),
		vetter.Vetter(applabel.NewVetter(informerFactory)),
		vetter.Vetter(serviceportprefix.NewVetter(informerFactory)),
		vetter.Vetter(serviceassociation.NewVetter(informerFactory)),
		vetter.Vetter(danglingroutedestinationhost.NewVetter(informerFactory)),
		vetter.Vetter(conflictingvirtualservicehost.NewVetter(informerFactory)),
	}

	stopCh := make(chan struct{})

	kubeInformerFactory.Start(stopCh)
	oks := kubeInformerFactory.WaitForCacheSync(stopCh)
	for inf, ok := range oks {
		if !ok {
			e := &adapter.Event{}
			e.EType = int32(meshes.EventType_ERROR)
			e.Details = ErrIstioVetSync(fmt.Errorf("%s", inf)).Error()
			e.Summary = "Failed to sync"
			ch <- e
			return
		}
	}

	istioInformerFactory.Start(stopCh)
	oks = istioInformerFactory.WaitForCacheSync(stopCh)
	for inf, ok := range oks {
		if !ok {
			e := &adapter.Event{}
			e.EType = int32(meshes.EventType_ERROR)
			e.Details = ErrIstioVetSync(fmt.Errorf("%s", inf)).Error()
			e.Summary = "Failed to sync"
			ch <- e
			return
		}
	}
	close(stopCh)

	for _, v := range vList {
		nList, err := v.Vet()
		if err != nil {
			e := &adapter.Event{}
			e.Summary = fmt.Sprintf("Vetter: %s reported error", v.Info().GetId())
			e.Details = err.Error()
			e.EType = int32(meshes.EventType_ERROR)
			ch <- e
			continue
		}
		if len(nList) > 0 {
			for i := range nList {
				e := &adapter.Event{}

				var ts []string
				for k, v := range nList[i].Attr {
					ts = append(ts, "${"+k+"}", v)
				}
				r := strings.NewReplacer(ts...)
				e.Summary = r.Replace(nList[i].GetSummary())
				e.Details = r.Replace(nList[i].GetMsg())
				switch nList[i].GetLevel().String() {
				case "WARNING":
					e.EType = int32(meshes.EventType_WARN)
				case "ERROR":
					e.EType = int32(meshes.EventType_ERROR)
				default:
					e.EType = int32(meshes.EventType_INFO)
				}
				ch <- e
			}
		} else {
			e := &adapter.Event{}
			istio.Log.Debug(fmt.Sprintf("Vetter %s ran successfully and generated no notes", v.Info().GetId()))
			e.Summary = fmt.Sprintf("Vetter: %s ran successfully", v.Info().GetId())
			e.Details = "No notes generated"
			e.EType = int32(meshes.EventType_INFO)
			ch <- e
		}
	}

	close(ch)
}

// StreamWarn streams a warning message to the channel
func (istio *Istio) StreamWarn(e *adapter.Event, err error) {
	istio.Log.Warn(err)
	e.EType = int32(meshes.EventType_WARN)
	*istio.Channel <- e
}
