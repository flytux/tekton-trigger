package template_test

import (
	"testing"

	"gitlab.com/pongsatt/githook/pkg/template"
)

func TestReplace(t *testing.T) {
	input := `"spec": {
        "params": [
            {
                "name": "appName",
                "value": "ch-back"
            },
            {
                "name": "pathToYamlFile",
                "value": "deploy"
            },
            {
                "name": "imageUrl",
                "value": "registry.app.pongzt.com/ch-back"
            },
            {
                "name": "imageTag",
                "value": "$REVISION"
            }
        ],
        "pipelineRef": {
            "name": "build-and-deploy-pipeline"
        },
        "resources": [
            {
                "name": "git-source",
                "resourceRef": {
                    "name": "ch-back-pipeline-git-source-k6lsh"
                }
            }
        ],
        "serviceAccount": "build-bot",
        "timeout": "10m0s"
    },`

	keyVals := make(map[string]string)
	keyVals["REVISION"] = "rev1"

	result := template.Replace(input, keyVals)

	t.Error(result)
}
