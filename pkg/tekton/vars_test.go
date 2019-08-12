package tekton

import (
	"testing"
)

func TestReplaceCommitVar(t *testing.T) {
	testcases := []struct {
		commitHash     string
		input          string
		expectedOutput string
	}{
		{
			commitHash:     "034ab39f12bac07af0188cc9fe7b9f18fba8731f",
			input:          "{image: \"test.com/$COMMIT\"}",
			expectedOutput: "{image: \"test.com/034ab39f12\"}",
		},
		{
			commitHash:     "034ab39f12bac07af0188cc9fe7b9f18fba8731f",
			input:          "{image: \"test.com/test\"}",
			expectedOutput: "{image: \"test.com/test\"}",
		},
		{
			commitHash:     "034ab39f12",
			input:          "{image: \"test.com/$COMMIT\"}",
			expectedOutput: "{image: \"test.com/034ab39f12\"}",
		},
		{
			commitHash:     "034",
			input:          "{image: \"test.com/$COMMIT\"}",
			expectedOutput: "{image: \"test.com/034\"}",
		},
		{
			commitHash:     "",
			input:          "{image: \"test.com/$COMMIT\"}",
			expectedOutput: "{image: \"test.com/\"}",
		},
	}

	for _, testcase := range testcases {
		opts := PipelineOptions{
			GitCommit: testcase.commitHash,
		}

		output := replaceVars(testcase.input, opts)

		if output != testcase.expectedOutput {
			t.Fatalf("expected : %s but got %s", testcase.expectedOutput, output)
		}
	}

}
