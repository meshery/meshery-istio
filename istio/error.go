// Package istio - Error codes for the adapter
package istio

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	// Errror code for failed service mesh installation
	ErrInstallIstioCode       = "istio_test_code"
	ErrUnzipFileCode          = "istio_test_code"
	ErrTarXZFCode             = "istio_test_code"
	ErrMeshConfigCode         = "istio_test_code"
	ErrFetchManifestCode      = "istio_test_code"
	ErrDownloadBinaryCode     = "istio_test_code"
	ErrInstallBinaryCode      = "istio_test_code"
	ErrClientConfigCode       = "istio_test_code"
	ErrClientSetCode          = "istio_test_code"
	ErrStreamEventCode        = "istio_test_code"
	ErrSampleAppCode          = "istio_test_code"
	ErrCustomOperationCode    = "istio_test_code"
	ErrAddonFromTemplateCode  = "istio_test_code"
	ErrAddonInvalidConfigCode = "istio_test_code"

	ErrOpInvalid = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")
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

// ErrFetchManifest is the error for mesh port forward
func ErrFetchManifest(err error, des string) error {
	return errors.NewDefault(ErrFetchManifestCode, fmt.Sprintf("Error fetching mesh manifest: %s", des))
}

// ErrDownloadBinary is the error while downloading istio binary
func ErrDownloadBinary(err error) error {
	return errors.NewDefault(ErrDownloadBinaryCode, fmt.Sprintf("Error downloading istio binary: %s", err.Error()))
}

// ErrInstallBinary is the error while downloading istio binary
func ErrInstallBinary(err error) error {
	return errors.NewDefault(ErrInstallBinaryCode, fmt.Sprintf("Error installing istio binary: %s", err.Error()))
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.NewDefault(ErrClientConfigCode, fmt.Sprintf("Error setting client config: %s", err.Error()))
}

// ErrClientSet is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.NewDefault(ErrClientSetCode, fmt.Sprintf("Error setting clientset: %s", err.Error()))
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.NewDefault(ErrStreamEventCode, fmt.Sprintf("Error streaming event: %s", err.Error()))
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.NewDefault(ErrSampleAppCode, fmt.Sprintf("Error with sample app operation: %s", err.Error()))
}

// ErrAddonFromTemplate is the error for streaming event
func ErrAddonFromTemplate(err error) error {
	return errors.NewDefault(ErrAddonFromTemplateCode, fmt.Sprintf("Error with addon install operation: %s", err.Error()))
}

// ErrInvalidConfig is the error for streaming event
func ErrAddonInvalidConfig(err error) error {
	return errors.NewDefault(ErrAddonInvalidConfigCode, fmt.Sprintf("Invalid addon: %s", err.Error()))
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.NewDefault(ErrCustomOperationCode, fmt.Sprintf("Error with custom operation: %s", err.Error()))
}
