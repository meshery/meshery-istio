package istio

import (
	"context"

	"github.com/layer5io/meshery-istio/meshes"
)

func (iClient *Client) executePrometheusInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getPrometheusYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeKialiInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getKialiYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeGrafanaInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getGrafanaYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeJaegerInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getJaegerYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeZipkinInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getZipkinYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}

func (iClient *Client) executeOperatorInstall(ctx context.Context, arReq *meshes.ApplyRuleRequest) error {
	if !arReq.DeleteOp {
		if err := iClient.labelNamespaceForAutoInjection(ctx, arReq.Namespace); err != nil {
			return err
		}
	}
	yamlFileContents, err := iClient.getOperatorYAML()
	if err != nil {
		return err
	}
	if err := iClient.applyConfigChange(ctx, yamlFileContents, arReq.Namespace, arReq.DeleteOp, false); err != nil {
		return err
	}
	return nil
}
