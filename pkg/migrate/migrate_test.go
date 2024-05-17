package migrate

import (
	"embed"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

//go:embed testdata/**/local.json
var testfolder embed.FS

func prepareUnstructred(input string) unstructured.Unstructured {
	jsonData := make(map[string]interface{})
	_ = json.Unmarshal([]byte(input), &jsonData)
	testObj := unstructured.Unstructured{}
	testObj.SetUnstructuredContent(jsonData)

	return testObj
}

func TestK3sRancherToRke2Rancher(t *testing.T) {
	tests := []struct {
		name     string
		newInput unstructured.Unstructured
		input    string
		expected string
	}{
		{
			name: "Test 1",
			input: `{
  "metadata": {
    "annotations": {
      "management.cattle.io/current-cluster-controllers-version": "v1.27.11+k3s1"
    },
    "labels": {
      "provider.cattle.io": "k3s"
    },
    "managedFields": [
      {
        "apiVersion": "management.cattle.io/v3",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:spec": {
            ".": {},
            "f:k3sConfig": {
              ".": {},
              "f:k3supgradeStrategy": {
                ".": {},
                "f:serverConcurrency": {},
                "f:workerConcurrency": {}
              },
              "f:kubernetesVersion": {}
            }
          }
        },
        "manager": "rancher",
        "operation": "Update",
        "time": "2024-05-02T13:47:19Z"
      }
    ]
  },
  "spec": {
    "k3sConfig": {
      "k3sUpgradeStrategy": {
        "serverConcurrency": 1,
        "workerConcurrency": 1
      },
      "kubernetesVersion": "v1.27.11+k3s1"
    }
  },
  "status": {
    "driver": "k3s",
    "provider": "k3s",
    "version": {
      "buildDate": "2024-02-29T20:10:42Z",
      "compiler": "gc",
      "gitCommit": "06d6bc80b469a61e5e90438b1f2639cd136a89e7",
      "gitTreeState": "clean",
      "gitVersion": "v1.27.11+k3s1",
      "goVersion": "go1.21.7",
      "major": "1",
      "minor": "27",
      "platform": "linux/amd64"
    }
  }
}`,
			expected: `{
  "metadata": {
    "annotations": {
      "management.cattle.io/current-cluster-controllers-version": "v1.27.11+rke2r1"
    },
    "labels": {
      "provider.cattle.io": "rke2"
    },
    "managedFields": [
      {
        "apiVersion": "management.cattle.io/v3",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:spec": {
            ".": {},
            "f:rke2Config": {
              ".": {},
              "f:rke2upgradeStrategy": {
                ".": {},
                "f:serverConcurrency": {},
                "f:workerConcurrency": {}
              },
              "f:kubernetesVersion": {}
            }
          }
        },
        "manager": "rancher",
        "operation": "Update",
        "time": "2024-05-02T13:47:19Z"
      }
    ]
  },
  "spec": {
    "rke2Config": {
      "rke2UpgradeStrategy": {
        "serverConcurrency": 1,
        "workerConcurrency": 1
      },
      "kubernetesVersion": "v1.27.11+rke2r1"
    }
  },
  "status": {
    "driver": "imported",
    "provider": "rke2",
    "version": {
      "platform": "linux/amd64"
    }
  }
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testObj := prepareUnstructred(test.input)
			logrus.Infof("Test obj before: %v", testObj)
			actual := K3sRancherToRke2Rancher(testObj)
			logrus.Infof("Test obj after: %v", actual)
			expected := prepareUnstructred(test.expected)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestRke2RancherToK3sRancher(t *testing.T) {
	tests := []struct {
		name     string
		newInput unstructured.Unstructured
		input    string
		expected string
	}{
		{
			name: "Test 1",
			input: `{
  "metadata": {
    "annotations": {
      "management.cattle.io/current-cluster-controllers-version": "v1.27.11+rke2r1"
    },
    "labels": {
      "provider.cattle.io": "rke2"
    },
    "managedFields": [
      {
        "apiVersion": "management.cattle.io/v3",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:spec": {
            ".": {},
            "f:rke2Config": {
              ".": {},
              "f:rke2upgradeStrategy": {
                ".": {},
                "f:serverConcurrency": {},
                "f:workerConcurrency": {}
              },
              "f:kubernetesVersion": {}
            }
          }
        },
        "manager": "rancher",
        "operation": "Update",
        "time": "2024-05-02T13:47:19Z"
      }
    ]
  },
  "spec": {
    "rke2Config": {
      "rke2UpgradeStrategy": {
        "serverConcurrency": 1,
        "workerConcurrency": 1
      },
      "kubernetesVersion": "v1.27.11+rke2r1"
    }
  },
  "status": {
    "driver": "rke2",
    "provider": "rke2",
    "version": {
      "buildDate": "2024-02-29T20:10:42Z",
      "compiler": "gc",
      "gitCommit": "06d6bc80b469a61e5e90438b1f2639cd136a89e7",
      "gitTreeState": "clean",
      "gitVersion": "v1.27.11+rke2r1",
      "goVersion": "go1.21.7",
      "major": "1",
      "minor": "27",
      "platform": "linux/amd64"
    }
  }
}`,
			expected: `{
  "metadata": {
    "annotations": {
      "management.cattle.io/current-cluster-controllers-version": "v1.27.11+k3s1"
    },
    "labels": {
      "provider.cattle.io": "k3s"
    },
    "managedFields": [
      {
        "apiVersion": "management.cattle.io/v3",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:spec": {
            ".": {},
            "f:k3sConfig": {
              ".": {},
              "f:k3supgradeStrategy": {
                ".": {},
                "f:serverConcurrency": {},
                "f:workerConcurrency": {}
              },
              "f:kubernetesVersion": {}
            }
          }
        },
        "manager": "rancher",
        "operation": "Update",
        "time": "2024-05-02T13:47:19Z"
      }
    ]
  },
  "spec": {
    "k3sConfig": {
      "k3sUpgradeStrategy": {
        "serverConcurrency": 1,
        "workerConcurrency": 1
      },
      "kubernetesVersion": "v1.27.11+k3s1"
    }
  },
  "status": {
    "driver": "imported",
    "provider": "k3s",
    "version": {
      "platform": "linux/amd64"
    }
  }
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testObj := prepareUnstructred(test.input)
			logrus.Infof("Test obj before: %v", testObj)
			actual := Rke2RancherToK3sRancher(testObj)
			logrus.Infof("Test obj after: %v", actual)
			expected := prepareUnstructred(test.expected)
			assert.Equal(t, expected, actual)
		})
	}
}
