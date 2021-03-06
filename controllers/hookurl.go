package controllers

import (
	"strings"

	servinv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"gitlab.com/pongsatt/githook/api/v1alpha1"
)

func getWebhookURL(source *v1alpha1.GitHook, ksvc *servinv1alpha1.Service) string {
	if ksvc.Status.DeprecatedDomain != "" {
		if source.Spec.SslVerify {
			return "https://" + ksvc.Status.DeprecatedDomain
		}
		return "http://" + ksvc.Status.DeprecatedDomain
	}

	webhookURL := ksvc.Status.URL.String()

	if source.Spec.SslVerify {
		webhookURL = strings.Replace(webhookURL, "http://", "https://", 1)
	} else {
		webhookURL = strings.Replace(webhookURL, "https://", "http://", 1)
	}

	return webhookURL
}
