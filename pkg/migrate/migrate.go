package migrate

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"regexp"
)

func K3sRancherToRke2Rancher(k3sLocalCluster unstructured.Unstructured) unstructured.Unstructured {
	// Replace "^k3s" with "^rke2" on both key and values
	res := findAndReplaceValue(k3sLocalCluster, `^k3s`, "rke2")
	res = findAndReplaceKey(res, `^f:k3s`, "f:rke2")
	res = findAndReplaceValue(res, `k3s([0-9]*)$`, "rke2r$1")
	return res
}

func Rke2RancherToK3sRancher(k3sLocalCluster unstructured.Unstructured) unstructured.Unstructured {
	// Replace "^k3s" with "^rke2" on both key and values
	res := findAndReplaceValue(k3sLocalCluster, `^rke2`, "k3s")
	res = findAndReplaceKey(res, `^f:rke2`, "f:k3s")
	res = findAndReplaceValue(res, `rke2r([0-9]*)$`, "k3s$1")
	return res
}

func findAndReplaceValue(data unstructured.Unstructured, targetKey string, replacement string) unstructured.Unstructured {
	matchRegex := regexp.MustCompile(targetKey)

	labels := data.GetLabels()
	if val, ok := labels["provider.cattle.io"]; ok && matchRegex.MatchString(val) {
		labels["provider.cattle.io"] = matchRegex.ReplaceAllString(val, replacement)
	}
	data.SetLabels(labels)

	annotations := data.GetAnnotations()
	if val, ok := labels["provider.cattle.io"]; ok && matchRegex.MatchString(val) {
		annotations["management.cattle.io/current-cluster-controllers-version"] = matchRegex.ReplaceAllString(val, replacement)
	}
	data.SetAnnotations(annotations)

	status, found, _ := unstructured.NestedFieldCopy(data.Object, "status")
	if found {
		status.(map[string]interface{})["driver"] = "imported"
		if matchRegex.MatchString(status.(map[string]string)["provider"]) {
			status.(map[string]interface{})["provider"] = matchRegex.ReplaceAllString(status.(map[string]string)["provider"], replacement)
		}
	}

	return data
}

func findAndReplaceKey(data unstructured.Unstructured, targetKey string, replacement string) unstructured.Unstructured {
	matchRegex := regexp.MustCompile(targetKey)
	spec := data.UnstructuredContent()
	managedFields := spec["managedFields"].([]map[string]interface{})[0]
	for key, value := range managedFields {
		if matchRegex.MatchString(key) {
			newKey := matchRegex.ReplaceAllString(key, replacement)
			delete(managedFields, key)
			managedFields[newKey] = value
		}
	}
	spec["managedFields"] = []interface{}{managedFields}

	for key, value := range spec {
		if matchRegex.MatchString(key) {
			newKey := matchRegex.ReplaceAllString(key, replacement)
			delete(spec, key)
			spec[newKey] = value
		}
	}

	data.SetUnstructuredContent(spec)

	return data
}
