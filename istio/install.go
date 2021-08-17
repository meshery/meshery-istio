package istio

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/errors"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

const (
	platform = runtime.GOOS
	arch = runtime.GOARCH
)

var (
	downloadLocation = os.TempDir()
)

// installs Istio using either helm charts or istioctl.
// Priority given to helm charts unless useBin set to true
func (istio *Istio) installIstio(del, useBin bool, version, namespace string) (string, error) {
	istio.Log.Debug(fmt.Sprintf("Requested install of version: %s", version))
	istio.Log.Debug(fmt.Sprintf("Requested action is delete: %v", del))
	istio.Log.Debug(fmt.Sprintf("Requested action is in namespace: %s", namespace))

	st := status.Installing

	if del {
		st = status.Removing
	}

	err := istio.Config.GetObject(adapter.MeshSpecKey, istio)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	// Fetch and/or return the path to downloaded and extracted release bundle
	dirName, err := istio.getIstioRelease(version)
	if err != nil {
		// ErrGettingIstioRelease
		return st, err
	}

	// Install using istioctl if explicitly stated
	if useBin {
		err = istio.runIstioCtlCmd(version, del, dirName)
		if err != nil {
			//ErrInstallUsingIstioCtl
			istio.Log.Error(ErrInstallIstio(err))
			return st, ErrInstallIstio(err)
		}
	}

	// Install using Helm Chart and fallback to istioctl
	err = istio.applyHelmChart(del, version, namespace, dirName)
	if err != nil {
		istio.Log.Error(ErrApplyHelmChart(err))

		err = istio.runIstioCtlCmd(version, del, dirName)
		if err != nil {
			//ErrInstallUsingIstioCtl
			istio.Log.Error(ErrInstallIstio(err))
			return st, ErrInstallIstio(err)
		}

		return st, nil
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (istio *Istio) applyHelmChart(del bool, version, namespace, dirName string) error {
	kClient := istio.MesheryKubeclient

	err := kClient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		LocalPath: path.Join(downloadLocation, dirName, "manifests/charts/base"),
		Namespace:       namespace,
		Delete:          del,
		CreateNamespace: true,
	})

	return err
}

// getIstioRelease gets the manifests for latest istio release.
// It first checks if the artifacts exist in OS's temp dir. If they don't,
// it proceeds to fetch the latest one.
func (istio *Istio) getIstioRelease(release string) (string, error) {
	releaseName := fmt.Sprintf("istio-%s", release)

	istio.Log.Info("Looking for artifacts of requested version locally...")
	_, err := os.Stat(path.Join(downloadLocation, releaseName))
	if err == nil {
		//ErrGettingIstioRelease
		return releaseName, nil
	}
	istio.Log.Info("Artifacts not found...")

	istio.Log.Info("Downloading requested istio version artifacts...")
	res, err := downloadTar(releaseName, release)
	if err != nil {
		//ErrGettingIstioRelease
		return "", err
	}

	err = extractTar(res)
	if err != nil {
		//ErrGettingIstioRelease
		return "", err
	}

	return releaseName, nil
}

func downloadTar(releaseName, release string) (*http.Response, error) {
	url := "https://github.com/istio/istio/releases/download"
	switch platform {
	case "darwin":
		url = fmt.Sprintf("%s/%s/%s-osx.tar.gz", url, release, releaseName)
	case "windows":
		url = fmt.Sprintf("%s/%s/%s-win.zip", url, release, releaseName)
	case "linux":
		url = fmt.Sprintf("%s/%s/%s-%s-%s.tar.gz", url, release, releaseName, platform, arch)
	default:
		//ErrUnsupportedPlatfrom
		return nil, errors.NewDefault("platform not supported")
	}

	resp, err := http.Get(url)
	if err != nil {
		// ErrDownloadTar
		return nil, err//ErrDownloadBinary(err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		// ErrDownloadTar
		return nil, err//ErrDownloadBinary(fmt.Errorf("bad status: %s", resp.Status))
	}

	return resp, nil
}

func extractTar(res *http.Response) error {
	// Close the response body
	defer func() {
	if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		if err := tarxzf(downloadLocation, res.Body); err != nil {
			//ErrExtracingFromTar
			return err//ErrInstallBinary(err)
		}
	case "windows":
		if err := unzip(downloadLocation, res.Body); err != nil {
			//ErrExtracingFromTar
			return err//ErrInstallBinary(err)
		}
	}

	return nil
}

// Installs Istio using Istioctl
// TODO: Figure out why this is not working in containers
func (istio *Istio) runIstioCtlCmd(version string, isDel bool, dirName string) error {
	var (
		out bytes.Buffer
		er  bytes.Buffer
	)

	Executable, err := istio.getExecutable(version)
	if err != nil {
		return ErrRunIstioCtlCmd(err, err.Error())
	}
	execCmd := []string{"install", "--set", "profile=demo", "-y"}
	if isDel {
		execCmd = []string{"x", "uninstall", "--purge", "-y"}
	}

	// We need a variable executable here hence using nosec
	// #nosec
	command := exec.Command(Executable, execCmd...)
	command.Stdout = &out
	command.Stderr = &er
	err = command.Run()
	if err != nil {
		return ErrRunIstioCtlCmd(err, er.String())
	}

	return nil
}

func (istio *Istio) applyManifest(contents []byte, isDel bool, namespace string) error {

	err := istio.MesheryKubeclient.ApplyManifest(contents, mesherykube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    isDel,
	})
	if err != nil {
		return err
	}

	return nil
}

// getExecutable looks for the executable in
// 1. $PATH
// 2. Root config path
//
// If it doesn't find the executable in the path then it proceeds
// to download the binary from github releases and installs it
// in the root config path
func (istio *Istio) getExecutable(release string) (string, error) {
	const platform = runtime.GOOS
	binaryName := generatePlatformSpecificBinaryName("istioctl", platform)
	alternateBinaryName := generatePlatformSpecificBinaryName("istioctl-"+release, platform)

	// Look for the executable in the path
	istio.Log.Info("Looking for istio in the path...")
	executable, err := exec.LookPath(binaryName)
	if err == nil {
		return executable, nil
	}
	executable, err = exec.LookPath(alternateBinaryName)
	if err == nil {
		return executable, nil
	}

	binPath := path.Join(config.RootPath(), "bin")

	// Look for config in the root path
	istio.Log.Info("Looking for istio in", binPath, "...")
	executable = path.Join(binPath, alternateBinaryName)
	if _, err := os.Stat(executable); err == nil {
		return executable, nil
	}

	// Proceed to download the binary in the config root path
	istio.Log.Info("istio not found in the path, downloading...")
	res, err := downloadBinary(platform, runtime.GOARCH, release)
	if err != nil {
		return "", err
	}
	// Install the binary
	istio.Log.Info("Installing...")
	if err = installBinary(binPath, platform, binaryName, res); err != nil {
		return "", err
	}
	// Rename the binary
	err = os.Rename(path.Join(binPath, binaryName), path.Join(binPath, alternateBinaryName))
	if err != nil {
		return "", err
	}

	istio.Log.Info("Done")
	return path.Join(binPath, alternateBinaryName), nil
}

func downloadBinary(platform, arch, release string) (*http.Response, error) {
	var url = "https://github.com/istio/istio/releases/download"
	switch platform {
	case "darwin":
		url = fmt.Sprintf("%s/%s/istioctl-%s-osx.tar.gz", url, release, release)
	case "windows":
		url = fmt.Sprintf("%s/%s/istioctl-%s-win.zip", url, release, release)
	case "linux":
		url = fmt.Sprintf("%s/%s/istioctl-%s-%s-%s.tar.gz", url, release, release, platform, arch)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrDownloadBinary(err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, ErrDownloadBinary(fmt.Errorf("bad status: %s", resp.Status))
	}

	return resp, nil
}

func installBinary(location, platform, name string, res *http.Response) error {
	// Close the response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// ~/.meshery/bin
	err := os.MkdirAll(location, 0750)
	if err != nil {
		return err
	}

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		if err := tarxzf(location, res.Body); err != nil {
			return ErrInstallBinary(err)
		}
		// Change permissions, we need the binary to be executable, hence
		// #nosec
		if err = os.Chmod(path.Join(location, name), 0750); err != nil {
			return err
		}
	case "windows":
		if err := unzip(location, res.Body); err != nil {
			return ErrInstallBinary(err)
		}
	}

	return nil
}

func tarxzf(location string, stream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return ErrTarXZF(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// File traversal is required to store the extracted manifests at the right place
			// #nosec
			if err := os.MkdirAll(path.Join(location, header.Name), 0750); err != nil {
				return ErrTarXZF(err)
			}
		case tar.TypeReg:
			// File traversal is required to store the extracted manifests at the right place
			// #nosec
			outFile, err := os.Create(path.Join(location, header.Name))
			if err != nil {
				return ErrTarXZF(err)
			}
			// Trust istioctl tar
			// #nosec
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return ErrTarXZF(err)
			}
			if err = outFile.Close(); err != nil {
				return ErrTarXZF(err)
			}

		default:
			return ErrTarXZF(err)
		}
	}

	return nil
}

func unzip(location string, zippedContent io.Reader) error {
	// Keep file in memory: Approx size ~ 50MB
	// TODO: Find a better approach
	zipped, err := ioutil.ReadAll(zippedContent)
	if err != nil {
		return ErrUnzipFile(err)
	}

	zReader, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		return ErrUnzipFile(err)
	}

	for _, file := range zReader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return ErrUnzipFile(err)
		}
		defer func() {
			if err := zippedFile.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		// need file traversal to place the extracted files at the right place, hence
		// #nosec
		extractedFilePath := path.Join(location, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(extractedFilePath, file.Mode()); err != nil {
				return ErrUnzipFile(err)
			}
		} else {
			// we need a variable path hence,
			// #nosec
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return ErrUnzipFile(err)
			}
			defer func() {
				if err := outputFile.Close(); err != nil {
					fmt.Println(err)
				}
			}()

			// Trust istio zip hence,
			// #nosec
			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return ErrUnzipFile(err)
			}
		}
	}

	return nil
}

func generatePlatformSpecificBinaryName(binName, platform string) string {
	if platform == "windows" && !strings.HasSuffix(binName, ".exe") {
		return binName + ".exe"
	}

	return binName
}
