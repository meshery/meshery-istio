// Package istio - Error codes for the adapter
package istio

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallIstioCode represents the errors which are generated
	// during istio service mesh install process
	ErrInstallIstioCode = "istio_test_code"

	// ErrUnzipFileCode represents the errors which are generated
	// during unzip process
	ErrUnzipFileCode = "istio_test_code"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "istio_test_code"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "istio_test_code"

	// ErrRunIstioCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunIstioCtlCmdCode = "istio_test_code"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during binary download process
	ErrDownloadBinaryCode = "istio_test_code"

	// ErrInstallBinaryCode represents the errors which are generated
	// during binary installation process
	ErrInstallBinaryCode = "istio_test_code"

	// ErrSampleAppCode represents the errors which are generated
	// duing sample app installation
	ErrSampleAppCode = "istio_test_code"

	// ErrEnvoyFilterCode represents the errors which are generated
	// duing envoy filter patching
	ErrEnvoyFilterCode = "istio_test_code"

	// ErrApplyPolicyCode represents the errors which are generated
	// duing policy apply operation
	ErrApplyPolicyCode = "istio_test_code"

	// ErrCustomOperationCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrCustomOperationCode = "istio_test_code"

	// ErrAddonFromTemplateCode represents the errors which are generated
	// during addon deployment process
	ErrAddonFromTemplateCode = "istio_test_code"

	// ErrAddonInvalidConfigCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrAddonInvalidConfigCode = "istio_test_code"

	// ErrCreatingIstioClientCode represents the errors which are generated
	// during creating istio client process
	ErrCreatingIstioClientCode = "istio_test_code"

	// ErrIstioVetSyncCode represents the errors which are generated
	// during istio-vet sync process
	ErrIstioVetSyncCode = "istio_test_code"

	// ErrIstioVetCode represents the errors which are generated
	// during istio-vet process
	ErrIstioVetCode = "istio_test_code"

	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "istio_test_code"

	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "istio_test_code"

	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "istio_test_code"

	// ErrParseIstioCoreComponentCode represents the error code which is
	// generated when istio core component manifest parsing fails
	ErrParseIstioCoreComponentCode = "istio_test_code"

	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "istio_test_code"

	// ErrOpInvalid represents the errors which are generated
	// when an invalid operation is requested
	ErrOpInvalid = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.NewDefault(ErrParseOAMComponentCode, "error parsing the component")

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.NewDefault(ErrParseOAMConfigCode, "error parsing the configuration")

	// ErrNilClient represents the error which is
	// generated when kubernetes client is nil
	ErrNilClient = errors.NewDefault(ErrNilClientCode, "kubernetes client not initialized")
)

// ErrInstallIstio is the error for install mesh
func ErrInstallIstio(err error) error {
	return errors.NewDefault(ErrInstallIstioCode, fmt.Sprintf("Error with istio operation: %s", err.Error()))
}

// ErrUnzipFile is the error for unzipping the file
func ErrUnzipFile(err error) error {
	return errors.NewDefault(ErrUnzipFileCode, fmt.Sprintf("Error while unzipping: %s", err.Error()))
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.NewDefault(ErrTarXZFCode, fmt.Sprintf("Error while extracting file: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.NewDefault(ErrMeshConfigCode, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrRunIstioCtlCmd is the error for mesh port forward
func ErrRunIstioCtlCmd(err error, des string) error {
	return errors.NewDefault(ErrRunIstioCtlCmdCode, fmt.Sprintf("Error running istioctl command: %s", des))
}

// ErrDownloadBinary is the error while downloading istio binary
func ErrDownloadBinary(err error) error {
	return errors.NewDefault(ErrDownloadBinaryCode, fmt.Sprintf("Error downloading istio binary: %s", err.Error()))
}

// ErrInstallBinary is the error while downloading istio binary
func ErrInstallBinary(err error) error {
	return errors.NewDefault(ErrInstallBinaryCode, fmt.Sprintf("Error installing istio binary: %s", err.Error()))
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.NewDefault(ErrSampleAppCode, fmt.Sprintf("Error with sample app operation: %s", err.Error()))
}

// ErrEnvoyFilter is the error for streaming event
func ErrEnvoyFilter(err error) error {
	return errors.NewDefault(ErrEnvoyFilterCode, fmt.Sprintf("Error with envoy filter operation: %s", err.Error()))
}

// ErrApplyPolicy is the error for streaming event
func ErrApplyPolicy(err error) error {
	return errors.NewDefault(ErrApplyPolicyCode, fmt.Sprintf("Error with apply policy operation: %s", err.Error()))
}

// ErrAddonFromTemplate is the error for streaming event
func ErrAddonFromTemplate(err error) error {
	return errors.NewDefault(ErrAddonFromTemplateCode, fmt.Sprintf("Error with addon install operation: %s", err.Error()))
}

// ErrAddonInvalidConfig is the error for streaming event
func ErrAddonInvalidConfig(err error) error {
	return errors.NewDefault(ErrAddonInvalidConfigCode, fmt.Sprintf("Invalid addon: %s", err.Error()))
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.NewDefault(ErrCustomOperationCode, fmt.Sprintf("Error with custom operation: %s", err.Error()))
}

// ErrCreatingIstioClient is the error for streaming event
func ErrCreatingIstioClient(err error) error {
	return errors.NewDefault(ErrCreatingIstioClientCode, fmt.Sprintf("Unable to create a new istio client %s", err.Error()))
}

// ErrIstioVetSync is the error for streaming event
func ErrIstioVetSync(err error) error {
	return errors.NewDefault(ErrIstioVetSyncCode, fmt.Sprintf("Failed to sync %s", err.Error()))
}

// ErrIstioVet is the error for streaming event
func ErrIstioVet(err error) error {
	return errors.NewDefault(ErrIstioVetCode, err.Error())
}

// ErrParseIstioCoreComponent is the error when istio core component manifest parsing fails
func ErrParseIstioCoreComponent(err error) error {
	return errors.NewDefault(ErrParseIstioCoreComponentCode, err.Error())
}

// ErrInvalidOAMComponentType is the error when the OAM component name is not valid
func ErrInvalidOAMComponentType(compName string) error {
	return errors.NewDefault(ErrInvalidOAMComponentTypeCode, "invalid OAM component name: ", compName)
}
