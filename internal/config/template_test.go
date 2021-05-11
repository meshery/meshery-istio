package config

import (
	"testing"
)

func TestGenerateImagehubTemplate(t *testing.T) {
	_, err := GenerateImagehubTemplates("stuff")
	if err != nil {
		t.Error(err)
	}
}
