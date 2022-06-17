package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFullyQualifyPackage(t *testing.T) {
	t.Run("valid FQP Format", func(t *testing.T) {
		valid := map[string]FullyQualifyPackage{
			"organization/name":        {"organization", "name", ""},
			"organization/name:v1.0":   {"organization", "name", "v1.0"},
			"organization/name:latest": {"organization", "name", "latest"},
			"name/organization:2.0":    {"name", "organization", "2.0"},
		}

		for fqp, expected := range valid {
			fqpVO, err := NewFullyQualifyPackage(fqp)
			require.Nil(t, err)

			assert.True(t, fqpVO.Equals(expected))
			assert.Equal(t, fqp, fqpVO.String())
		}
	})

	t.Run("Invalid FQP Format", func(t *testing.T) {
		invalid := []string{
			"organization//name",
			"/organization/name",
			"organization/name::",
			"organization/name:^1.0",
			"organization/name:~1.0",
			"organization/name:1-0",
			"organization /name:1.0",
			"organization/name :1.0",
			"organization\\/name:1.0",
		}

		for _, fqp := range invalid {
			_, err := NewFullyQualifyPackage(fqp)
			assert.Equal(t, err, ErrFullyQualifyPackageInvalidFormat)
		}
	})

	t.Run("Mutate FQP", func(t *testing.T) {
		mutations := map[string]string{
			"organization/name":        "organization/newName",
			"organization/name:v1.0":   "organization/name:newVersion",
			"organization/name:latest": "organization/newName:newVersion",
		}

		for fqp, expected := range mutations {
			fqpVO, err := NewFullyQualifyPackage(fqp)
			require.Nil(t, err)

			expectedVO, err := NewFullyQualifyPackage(expected)
			require.Nil(t, err)

			mutation := fqpVO.CopyWithName(expectedVO.Name()).CopyWithVersion(expectedVO.Version())

			assert.True(t, expectedVO.Equals(mutation))
			assert.NotSame(t, expectedVO, mutation)
		}
	})
}
