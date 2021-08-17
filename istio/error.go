// Package istio - Error codes for the adapter
package istio

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallIstioCode represents the errors which are generated
	// during istio service mesh install process
	ErrInstallIstioCode = "istio_4"

	// ErrUnzipFileCode represents the errors which are generated
	// during unzip process
	ErrUnzipFileCode = "istio_5"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "istio_6"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "istio_7"

	// ErrRunIstioCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunIstioCtlCmdCode = "istio_8"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during binary download process
	ErrDownloadBinaryCode = "istio_9"

	// ErrInstallBinaryCode represents the errors which are generated
	// during binary installation process
	ErrInstallBinaryCode = "istio_10"

	// ErrSampleAppCode represents the errors which are generated
	// duing sample app installation
	ErrSampleAppCode = "istio_11"

	// ErrEnvoyFilterCode represents the errors which are generated
	// duing envoy filter patching
	ErrEnvoyFilterCode = "istio_12"

	// ErrApplyPolicyCode represents the errors which are generated
	// duing policy apply operation
	ErrApplyPolicyCode = "istio_13"

	// ErrCustomOperationCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrCustomOperationCode = "istio_14"

	// ErrAddonFromTemplateCode represents the errors which are generated
	// during addon deployment process
	ErrAddonFromTemplateCode = "istio_15"

	// ErrAddonInvalidConfigCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrAddonInvalidConfigCode = "istio_16"

	// ErrCreatingIstioClientCode represents the errors which are generated
	// during creating istio client process
	ErrCreatingIstioClientCode = "istio_17"

	// ErrIstioVetSyncCode represents the errors which are generated
	// during istio-vet sync process
	ErrIstioVetSyncCode = "istio_18"

	// ErrIstioVetCode represents the errors which are generated
	// during istio-vet process
	ErrIstioVetCode = "istio_19"

	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "istio_20"

	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "istio_21"

	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "istio_22"

	// ErrParseIstioCoreComponentCode represents the error code which is
	// generated when istio core component manifest parsing fails
	ErrParseIstioCoreComponentCode = "istio_23"

	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "istio_24"

	// ErrOpInvalidCode represents the error code which is
	// generated when an invalid operation is requested
	ErrOpInvalidCode = "istio_25"

	// ErrIstioCoreComponentFailCode represents the error code which is
	// generated when an istio core operations fails
	ErrIstioCoreComponentFailCode = "istio_26"

	// ErrProcessOAMCode represents the error code which is
	// generated when an OAM operations fails
	ErrProcessOAMCode = "istio_27"

	// ErrOpInvalid represents the errors which are generated
	// when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{"Istio adapter recived an invalid operation from the meshey server"}, []string{"The operation is not supported by the adapter", "Invalid operation name"}, []string{"Check if the operation name is valid and supported by the adapter"})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occured while prasing application component in the OAM request made"}, []string{"Invalid OAM component passed in OAM request"}, []string{"Check if your request has vaild OAM components"})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occured while prasing component config in the OAM request made"}, []string{"Invalid OAM config passed in OAM request"}, []string{"Check if your request has vaild OAM config"})

	// ErrNilClient represents the error which is
	// generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"kubernetes client not initialized"}, []string{"Reconnect the adaptor to Meshery server"})
)

// ErrInstallIstio is the error for install mesh
func ErrInstallIstio(err error) error {
	return errors.New(ErrInstallIstioCode, errors.Alert, []string{"Error with istio operation"}, []string{"Error occured while installing istio mesh through istioctl", err.Error()}, []string{}, []string{})
}

// ErrUnzipFile is the error for unzipping the file
func ErrUnzipFile(err error) error {
	return errors.New(ErrUnzipFileCode, errors.Alert, []string{"Error while unzipping"}, []string{err.Error()}, []string{"File might be corrupt"}, []string{})
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.New(ErrTarXZFCode, errors.Alert, []string{"Error while extracting file"}, []string{err.Error()}, []string{"/The gzip might be corrupt"}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error(), "Error getting MeshSpecKey config from in-memory configuration"}, []string{}, []string{})
}

// ErrRunIstioCtlCmd is the error for mesh port forward
func ErrRunIstioCtlCmd(err error, des string) error {
	return errors.New(ErrRunIstioCtlCmdCode, errors.Alert, []string{"Error running istioctl command"}, []string{err.Error()}, []string{"Corrupted istioctl binary", "Command might be invalid"}, []string{})
}

// ErrDownloadBinary is the error while downloading istio binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading istio binary"}, []string{err.Error(), "Error occured while download istio binary from its github release"}, []string{"Checkout https://docs.github.com/en/rest/reference/repos#releases for more details"}, []string{})
}

// ErrInstallBinary is the error while downloading istio binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing istio binary"}, []string{err.Error()}, []string{"Corrupted istio release binary", "Invalid installation location"}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error(), "Error occured while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to meshery server to refresh the kubeclient"})
}

// ErrEnvoyFilter is the error for streaming event
func ErrEnvoyFilter(err error) error {
	return errors.New(ErrEnvoyFilterCode, errors.Alert, []string{"Error with envoy filter operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyPolicy is the error for streaming event
func ErrApplyPolicy(err error) error {
	return errors.New(ErrApplyPolicyCode, errors.Alert, []string{"Error with apply policy operation"}, []string{err.Error(), "Error occured while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to meshery server to refresh the kubeclient"})
}

// ErrAddonFromTemplate is the error for streaming event
func ErrAddonFromTemplate(err error) error {
	return errors.New(ErrAddonFromTemplateCode, errors.Alert, []string{"Error with addon install operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation"}, []string{"Error occured while applying custom manifest to the cluster", err.Error()}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{})
}

// ErrCreatingIstioClient is the error for streaming event
func ErrCreatingIstioClient(err error) error {
	return errors.New(ErrCreatingIstioClientCode, errors.Alert, []string{"Unable to create a new istio client"}, []string{err.Error()}, []string{}, []string{})
}

// ErrIstioVetSync is the error for streaming event
func ErrIstioVetSync(err error) error {
	return errors.New(ErrIstioVetSyncCode, errors.Alert, []string{"Failed to sync"}, []string{err.Error()}, []string{}, []string{})
}

// ErrIstioVet is the error for streaming event
func ErrIstioVet(err error) error {
	return errors.New(ErrIstioVetCode, errors.Alert, []string{"error while running Istio vet command"}, []string{err.Error()}, []string{}, []string{})
}

// ErrParseIstioCoreComponent is the error when istio core component manifest parsing fails
func ErrParseIstioCoreComponent(err error) error {
	return errors.New(ErrParseIstioCoreComponentCode, errors.Alert, []string{"istio core component manifest parsing failing"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInvalidOAMComponentType is the error when the OAM component name is not valid
func ErrInvalidOAMComponentType(compName string) error {
	return errors.New(ErrInvalidOAMComponentTypeCode, errors.Alert, []string{"invalid OAM component name: ", compName}, []string{}, []string{}, []string{})
}

// ErrIstioCoreComponentFail is the error when core Istion component processing fails
func ErrIstioCoreComponentFail(err error) error {
	return errors.New(ErrIstioCoreComponentFailCode, errors.Alert, []string{"error in istio core component"}, []string{err.Error()}, []string{}, []string{})
}

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}
