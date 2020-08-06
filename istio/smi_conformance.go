package istio

import (
	"context"
	"fmt"

	"github.com/layer5io/learn-layer5/smi-conformance/conformance"
	"github.com/layer5io/meshery-istio/meshes"
	"github.com/sirupsen/logrus"
)

type ConformanceResponse struct {
	Tests    string                       `json:"tests,omitempty"`
	Failures string                       `json:"failures,omitempty"`
	Results  []*SingleConformanceResponse `json:"results,omitempty"`
}

type Failure struct {
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`
}

type SingleConformanceResponse struct {
	Name       string   `json:"name,omitempty"`
	Time       string   `json:"time,omitempty"`
	Assertions string   `json:"assertions,omitempty"`
	Failure    *Failure `json:"failure,omitempty"`
}

func (iClient *Client) runConformanceTest(adaptorname string, arReq *meshes.ApplyRuleRequest) error {
	annotations := make(map[string]string, 0)
	// err := json.Unmarshal([]byte(arReq.CustomBody), &annotations)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return errors.Wrapf(err, "Error unmarshaling annotation body.")
	// }

	go func(name string, a map[string]string, req *meshes.ApplyRuleRequest) {

		cClient, err := conformance.CreateClient(context.TODO(), "127.0.0.1:53110")
		if err != nil {
			logrus.Error(err)
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: req.OperationId,
				EventType:   meshes.EventType_ERROR,
				Summary:     "Error creating a smi conformance tool client.",
				Details:     err.Error(),
			}
			return
		}
		logrus.Debugf("created client for smi conformance tool: %s", adaptorname)

		result, err := cClient.CClient.RunTest(context.TODO(), &conformance.Request{
			Annotations: a,
			Meshname:    name,
		})
		if err != nil {
			logrus.Error(err)
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: req.OperationId,
				EventType:   meshes.EventType_ERROR,
				Summary:     "Test failed",
				Details:     err.Error(),
			}
			return
		}
		logrus.Debugf("Tests ran successfully for smi conformance tool!!")

		response := ConformanceResponse{
			Tests:    result.Tests,
			Failures: result.Failures,
			Results:  make([]*SingleConformanceResponse, 0),
		}

		if result == nil {
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: req.OperationId,
				EventType:   meshes.EventType_ERROR,
				Summary:     "SMI tool connection crashed!",
				Details:     "Smi-conformance tool unreachable",
			}
			return
		}

		for _, res := range result.SingleTestResult {
			response.Results = append(response.Results, &SingleConformanceResponse{
				Name:       res.Name,
				Time:       res.Time,
				Assertions: res.Assertions,
				Failure: &Failure{
					Text:    res.Failure.Test,
					Message: res.Failure.Message,
				},
			})
		}

		logrus.Debugf(fmt.Sprintf("Tests Results: %+v", response))

		iClient.eventChan <- &meshes.EventsResponse{
			OperationId: req.OperationId,
			EventType:   meshes.EventType_INFO,
			Summary:     fmt.Sprintf("Tests Results: %+v", response),
			Details:     "Test completed successfully",
		}

		// _ = cClient.Close()
		return

	}(adaptorname, annotations, arReq)

	return nil
}
