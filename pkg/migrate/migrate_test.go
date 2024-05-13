package migrate

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

//go:embed testdata/**/local.json
var testfolder embed.FS

func TestK3sRancherToRke2Rancher(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Test 1",
			input: map[string]interface{}{
				"k3s": []interface{}{
					"foo",
					"bar",
				},
				"f:k3s": []interface{}{
					[]interface{}{"baz"},
					[]interface{}{"qux"},
				},
			},
			expected: map[string]interface{}{
				"rke2": []interface{}{
					"foo",
					"bar",
				},
				"f:rke2": []interface{}{
					[]interface{}{"baz"},
					[]interface{}{"qux"},
				},
			},
		},
		{
			name: "Test 2",
			input: map[string]interface{}{
				"k3s": []interface{}{
					[]interface{}{
						map[string]interface{}{
							"f:k3s": "foo",
						},
						map[string]interface{}{
							"k3s":              []interface{}{"bar"},
							"something-else":   "notk3sMATCHasdSOMEHASH",
							"provider.rancher": "k3s",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"rke2": []interface{}{
					[]interface{}{
						map[string]interface{}{
							"f:rke2": "foo",
						},
						map[string]interface{}{
							"rke2":             []interface{}{"bar"},
							"something-else":   "notk3sMATCHasdSOMEHASH",
							"provider.rancher": "rke2",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testObj := unstructured.Unstructured{}
			testObj.SetUnstructuredContent(test.input)
			actual := K3sRancherToRke2Rancher(testObj)
			assert.Equal(t, test.expected, actual.UnstructuredContent())
		})
	}
}
