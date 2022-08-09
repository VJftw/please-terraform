package root_test

import (
	"fmt"
	"testing"

	"github.com/VJftw/please-terraform/pkg/root"
	"github.com/stretchr/testify/assert"
)

func TestAutoTFVarsName(t *testing.T) {
	var tests = []struct {
		inIndex             int
		inVarFile           string
		expectedOutFileName string
		hasError            bool
	}{
		{0, "a.tfvars", "0-a.auto.tfvars", false},
		{0, "b.tfvars.json", "0-b.auto.tfvars.json", false},
		{1, "a.tfvars", "1-a.auto.tfvars", false},
		{1, "b.tfvars.json", "1-b.auto.tfvars.json", false},
		{2, "a.tfvars", "2-a.auto.tfvars", false},
		{2, "b.tfvars.json", "2-b.auto.tfvars.json", false},
		{0, "a", "", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", tt.inIndex, tt.inVarFile), func(t *testing.T) {
			newName, err := root.AutoTFVarsName(tt.inIndex, tt.inVarFile)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutFileName, newName)
		})
	}
}
