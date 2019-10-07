package istio

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	repoURL     = "https://api.github.com/repos/istio/istio/releases/latest"
	urlSuffix   = "-linux.tar.gz"
	crdPattern  = "crd(.*)yaml"
	cachePeriod = 6 * time.Hour
)

var (
	localByPassFile = "/app/istio.tar.gz"

	localFile                  = path.Join(os.TempDir(), "istio.tar.gz")
	destinationFolder          = path.Join(os.TempDir(), "istio")
	basePath                   = path.Join(destinationFolder, "%s")
	installFile                = path.Join(basePath, "install/kubernetes/istio-demo.yaml")
	installWithmTLSFile        = path.Join(basePath, "install/kubernetes/istio-demo-auth.yaml")
	bookInfoInstallFile        = path.Join(basePath, "samples/bookinfo/platform/kube/bookinfo.yaml")
	bookInfoGatewayInstallFile = path.Join(basePath, "samples/bookinfo/networking/bookinfo-gateway.yaml")
	crdFolder                  = path.Join(basePath, "install/kubernetes/helm/istio-init/files/")

	defaultBookInfoDestRulesFile                 = path.Join(basePath, "samples/bookinfo/networking/destination-rule-all-mtls.yaml")
	bookInfoRouteToV1AllServicesFile             = path.Join(basePath, "samples/bookinfo/networking/virtual-service-all-v1.yaml")
	bookInfoRouteToReviewsV2ForJasonFile         = path.Join(basePath, "samples/bookinfo/networking/virtual-service-reviews-test-v2.yaml")
	bookInfoCanary50pcReviewsV3File              = path.Join(basePath, "samples/bookinfo/networking/virtual-service-reviews-50-v3.yaml")
	bookInfoCanary100pcReviewsV3File             = path.Join(basePath, "samples/bookinfo/networking/virtual-service-reviews-v3.yaml")
	bookInfoInjectDelayForRatingsForJasonFile    = path.Join(basePath, "samples/bookinfo/networking/virtual-service-ratings-test-delay.yaml")
	bookInfoInjectHTTPAbortToRatingsForJasonFile = path.Join(basePath, "samples/bookinfo/networking/virtual-service-ratings-test-abort.yaml")
)

type apiInfo struct {
	TagName    string   `json:"tag_name,omitempty"`
	PreRelease bool     `json:"prerelease,omitempty"`
	Assets     []*asset `json:"assets,omitempty"`
}

type asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

func (iClient *Client) getLatestReleaseURL() error {
	if iClient.istioReleaseDownloadURL == "" || time.Since(iClient.istioReleaseUpdatedAt) > cachePeriod {
		logrus.Debugf("API info url: %s", repoURL)
		resp, err := http.Get(repoURL)
		if err != nil {
			err = errors.Wrapf(err, "error getting latest version info")
			logrus.Error(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("unable to fetch release info due to an unexpected status code: %d", resp.StatusCode)
			logrus.Error(err)
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrapf(err, "error parsing response body")
			logrus.Error(err)
			return err
		}
		// logrus.Debugf("Raw api info: %s", body)
		result := &apiInfo{}
		err = json.Unmarshal(body, result)
		if err != nil {
			err = errors.Wrapf(err, "error unmarshalling response body")
			logrus.Error(err)
			return err
		}
		logrus.Debugf("retrieved api info: %+#v", result)
		if result != nil && result.Assets != nil && len(result.Assets) > 0 {
			for _, asset := range result.Assets {
				if strings.HasSuffix(asset.Name, urlSuffix) {
					iClient.istioReleaseVersion = strings.Replace(asset.Name, urlSuffix, "", -1)
					iClient.istioReleaseDownloadURL = asset.DownloadURL
					iClient.istioReleaseUpdatedAt = time.Now()
					return nil
				}
			}
		}
		err = errors.New("unable to extract the download URL")
		logrus.Error(err)
		return err
	}
	return nil
}

func (iClient *Client) downloadFile(localFile string) error {
	dFile, err := os.Create(localFile)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a file on the filesystem at %s", localFile)
		logrus.Error(err)
		return err
	}
	defer dFile.Close()

	resp, err := http.Get(iClient.istioReleaseDownloadURL)
	if err != nil {
		err = errors.Wrapf(err, "unable to download the file from URL: %s", iClient.istioReleaseDownloadURL)
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to download the file from URL: %s, status: %s", iClient.istioReleaseDownloadURL, resp.Status)
		logrus.Error(err)
		return err
	}

	_, err = io.Copy(dFile, resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "unable to write the downloaded file to the file system at %s", localFile)
		logrus.Error(err)
		return err
	}
	return nil
}

func (iClient *Client) untarPackage(destination, fileToUntar string) error {
	lFile, err := os.Open(fileToUntar)
	if err != nil {
		err = errors.Wrapf(err, "unable to read the local file %s", fileToUntar)
		logrus.Error(err)
		return err
	}

	gzReader, err := gzip.NewReader(lFile)
	if err != nil {
		err = errors.Wrap(err, "unable to load the file into a gz reader")
		logrus.Error(err)
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			err = errors.Wrap(err, "error during untar")
			logrus.Error(err)
			return err
		case header == nil:
			continue
		}

		fileInLoop := filepath.Join(destination, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(fileInLoop); err != nil {
				if err := os.MkdirAll(fileInLoop, 0755); err != nil {
					err = errors.Wrapf(err, "error creating directory %s", fileInLoop)
					logrus.Error(err)
					return err
				}
			}
		case tar.TypeReg:
			fileAtLoc, err := os.OpenFile(fileInLoop, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				err = errors.Wrapf(err, "error opening file %s", fileInLoop)
				logrus.Error(err)
				return err
			}

			if _, err := io.Copy(fileAtLoc, tarReader); err != nil {
				err = errors.Wrapf(err, "error writing file %s", fileInLoop)
				logrus.Error(err)
				return err
			}
			fileAtLoc.Close()
		}
	}
}

func (iClient *Client) downloadIstio() (string, error) {
	var fileName string
	_, err := os.Stat(localByPassFile)
	if err != nil {
		logrus.Debug("preparing to download the latest istio release")
		err := iClient.getLatestReleaseURL()
		if err != nil {
			return "", err
		}
		fileName = iClient.istioReleaseVersion
		downloadURL := iClient.istioReleaseDownloadURL
		logrus.Debugf("retrieved latest file name: %s and download url: %s", fileName, downloadURL)

		proceedWithDownload := true

		lFileStat, err := os.Stat(localFile)
		if err == nil {
			if time.Since(lFileStat.ModTime()) > cachePeriod {
				proceedWithDownload = true
			} else {
				proceedWithDownload = false
			}
		}

		if proceedWithDownload {
			if err = iClient.downloadFile(localFile); err != nil {
				return "", err
			}
			logrus.Debug("package successfully downloaded, now unzipping . . .")
		}
	} else {
		localFile = localByPassFile
		fileName = os.Getenv("ISTIO_VERSION")
		logrus.Debugf("using local bypass file: %s & version name from env: %s", localFile, fileName)
	}
	if err = iClient.untarPackage(destinationFolder, localFile); err != nil {
		return "", err
	}
	logrus.Debug("successfully unzipped")
	return fileName, nil
}

func (iClient *Client) getIstioComponentYAML(fileName string) (string, error) {
	specificVersionName, err := iClient.downloadIstio()
	if err != nil {
		return "", err
	}
	installFileLoc := fmt.Sprintf(fileName, specificVersionName)
	logrus.Debugf("checking if install file exists at path: %s", installFileLoc)
	_, err = os.Stat(installFileLoc)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Error(err)
			return "", err
		}
		err = errors.Wrap(err, "unknown error")
		logrus.Error(err)
		return "", err
	}
	fileContents, err := ioutil.ReadFile(installFileLoc)
	if err != nil {
		err = errors.Wrap(err, "unable to read file")
		logrus.Error(err)
		return "", err
	}
	return string(fileContents), nil
}

func (iClient *Client) getCRDsYAML() ([]string, error) {
	res := []string{}

	rEx, err := regexp.Compile(crdPattern)
	if err != nil {
		err = errors.Wrap(err, "unable to compile crd pattern")
		logrus.Error(err)
		return nil, err
	}

	specificVersionName, err := iClient.downloadIstio()
	if err != nil {
		return nil, err
	}
	startFolder := fmt.Sprintf(crdFolder, specificVersionName)
	err = filepath.Walk(startFolder, func(currentPath string, info os.FileInfo, err error) error {
		if err == nil && rEx.MatchString(info.Name()) {
			contents, err := ioutil.ReadFile(currentPath)
			if err != nil {
				err = errors.Wrap(err, "unable to read file")
				logrus.Error(err)
				return err
			}
			res = append(res, string(contents))
		}
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "unable to read the directory")
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (iClient *Client) getLatestIstioYAML(installmTLS bool) (string, error) {
	if installmTLS {
		return iClient.getIstioComponentYAML(installWithmTLSFile)
	}
	return iClient.getIstioComponentYAML(installFile)
}

func (iClient *Client) getBookInfoAppYAML() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoInstallFile)
}

func (iClient *Client) getBookInfoGatewayYAML() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoGatewayInstallFile)
}

func (iClient *Client) getBookInfoDefaultDesinationRulesYAML() (string, error) {
	return iClient.getIstioComponentYAML(defaultBookInfoDestRulesFile)
}

func (iClient *Client) getBookInfoRouteToV1AllServicesYAML() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoRouteToV1AllServicesFile)
}

func (iClient *Client) getBookInfoRouteToReviewsV2ForJasonFile() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoRouteToReviewsV2ForJasonFile)
}

func (iClient *Client) getBookInfoCanary50pcReviewsV3File() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoCanary50pcReviewsV3File)
}

func (iClient *Client) getBookInfoCanary100pcReviewsV3File() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoCanary100pcReviewsV3File)
}

func (iClient *Client) getBookInfoInjectDelayForRatingsForJasonFile() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoInjectDelayForRatingsForJasonFile)
}

func (iClient *Client) getBookInfoInjectHTTPAbortToRatingsForJasonFile() (string, error) {
	return iClient.getIstioComponentYAML(bookInfoInjectHTTPAbortToRatingsForJasonFile)
}
