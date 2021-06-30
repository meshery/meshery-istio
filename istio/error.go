// Package istio - Error codes for the adapter
package istio

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallIstioCode represents the errors which are generated
	// during istio service mesh install process
	ErrInstallIstioCode = "1002"

	// ErrUnzipFileCode represents the errors which are generated
	// during unzip process
	ErrUnzipFileCode = "1003"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "1004"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "1005"

	// ErrRunIstioCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunIstioCtlCmdCode = "1006"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during binary download process
	ErrDownloadBinaryCode = "1007"

	// ErrInstallBinaryCode represents the errors which are generated
	// during binary installation process
	ErrInstallBinaryCode = "1008"

	// ErrSampleAppCode represents the errors which are generated
	// duing sample app installation
	ErrSampleAppCode = "1009"

	// ErrEnvoyFilterCode represents the errors which are generated
	// duing envoy filter patching
	ErrEnvoyFilterCode = "1010"

	// ErrApplyPolicyCode represents the errors which are generated
	// duing policy apply operation
	ErrApplyPolicyCode = "1011"

	// ErrCustomOperationCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrCustomOperationCode = "1012"

	// ErrAddonFromTemplateCode represents the errors which are generated
	// during addon deployment process
	ErrAddonFromTemplateCode = "1013"

	// ErrAddonInvalidConfigCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrAddonInvalidConfigCode = "1014"

	// ErrCreatingIstioClientCode represents the errors which are generated
	// during creating istio client process
	ErrCreatingIstioClientCode = "1015"

	// ErrIstioVetSyncCode represents the errors which are generated
	// during istio-vet sync process
	ErrIstioVetSyncCode = "1016"

	// ErrIstioVetCode represents the errors which are generated
	// during istio-vet process
	ErrIstioVetCode = "1017"

	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "1018"

	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "1019"

	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "1020"

	// ErrParseIstioCoreComponentCode represents the error code which is
	// generated when istio core component manifest parsing fails
	ErrParseIstioCoreComponentCode = "1021"

	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "1022"

	// ErrOpInvalidCode represents the error code which is
	// generated when an invalid operation is requested
	ErrOpInvalidCode = "1023"

	// ErrOpInvalid represents the errors which are generated
	// when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{}, []string{}, []string{})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{}, []string{}, []string{})

	// ErrNilClient represents the error which is
	// generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{}, []string{}, []string{})
)

// ErrInstallIstio is the error for install mesh
func ErrInstallIstio(err error) error {
	return errors.New(ErrInstallIstioCode, errors.Alert, []string{"Error with istio operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrUnzipFile is the error for unzipping the file
func ErrUnzipFile(err error) error {
	return errors.New(ErrUnzipFileCode, errors.Alert, []string{"Error while unzipping"}, []string{err.Error()}, []string{}, []string{})
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.New(ErrTarXZFCode, errors.Alert, []string{"Error while extracting file"}, []string{err.Error()}, []string{}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error()}, []string{}, []string{})
}

// ErrRunIstioCtlCmd is the error for mesh port forward
func ErrRunIstioCtlCmd(err error, des string) error {
	return errors.New(ErrRunIstioCtlCmdCode, errors.Alert, []string{"Error running istioctl command"}, []string{err.Error()}, []string{}, []string{})
}

// ErrDownloadBinary is the error while downloading istio binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading istio binary"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInstallBinary is the error while downloading istio binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing istio binary"}, []string{err.Error()}, []string{}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrEnvoyFilter is the error for streaming event
func ErrEnvoyFilter(err error) error {
	return errors.New(ErrEnvoyFilterCode, errors.Alert, []string{"Error with envoy filter operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyPolicy is the error for streaming event
func ErrApplyPolicy(err error) error {
	return errors.New(ErrApplyPolicyCode, errors.Alert, []string{"Error with apply policy operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrAddonFromTemplate is the error for streaming event
func ErrAddonFromTemplate(err error) error {
	return errors.New(ErrAddonFromTemplateCode, errors.Alert, []string{"Error with addon install operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrAddonInvalidConfig is the error for streaming event
func ErrAddonInvalidConfig(err error) error {
	return errors.New(ErrAddonInvalidConfigCode, errors.Alert, []string{"Invalid addon"}, []string{err.Error()}, []string{}, []string{})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation"}, []string{err.Error()}, []string{}, []string{})
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
