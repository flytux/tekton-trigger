package controllers

import (
	"testing"

	"github.com/knative/pkg/apis"
	servinv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"gitlab.com/pongsatt/githook/api/v1alpha1"
)

func TestGetWebhookURL(t *testing.T) {
	tests := []struct {
		domainURL   string
		url         string
		verifySSL   bool
		expectedURL string
	}{
		{
			domainURL:   "testdomain.com",
			url:         "http://test.com",
			verifySSL:   false,
			expectedURL: "http://testdomain.com",
		},
		{
			domainURL:   "testdomain.com",
			url:         "http://test.com",
			verifySSL:   true,
			expectedURL: "https://testdomain.com",
		},
		{
			domainURL:   "testdomain.com",
			url:         "",
			verifySSL:   true,
			expectedURL: "https://testdomain.com",
		},
		{
			domainURL:   "",
			url:         "http://test.com",
			verifySSL:   true,
			expectedURL: "https://test.com",
		},
		{
			domainURL:   "",
			url:         "https://test.com",
			verifySSL:   false,
			expectedURL: "http://test.com",
		},
	}

	for _, test := range tests {
		source := &v1alpha1.GitHook{
			Spec: v1alpha1.GitHookSpec{
				SslVerify: test.verifySSL,
			},
		}

		url, _ := apis.ParseURL(test.url)

		ksvc := &servinv1alpha1.Service{
			Status: servinv1alpha1.ServiceStatus{
				RouteStatusFields: servinv1alpha1.RouteStatusFields{
					URL:              url,
					DeprecatedDomain: test.domainURL,
				},
			},
		}

		resultURL := getWebhookURL(source, ksvc)

		if resultURL != test.expectedURL {
			t.Errorf("expected %s but got %s", test.expectedURL, resultURL)
		}
	}

}
