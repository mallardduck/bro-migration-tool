package migrate

import (
	"fmt"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"regexp"
	"strings"
)

func K3sRancherToRke2Rancher(k3sLocalCluster unstructured.Unstructured) unstructured.Unstructured {
	// Replace "^k3s" with "^rke2" on both key and values
	res := findAndReplaceLabelsValue(k3sLocalCluster, `^k3s$`, "rke2")
	res = findAndReplaceAnnotationsValue(res, `k3s([0-9]*)$`, "rke2r$1")
	res = findAndReplaceStatusValue(k3sLocalCluster, `^k3s$`, "rke2")
	res = findAndReplaceValueOn(res, "spec.k3sConfig", `k3s([0-9]*)$`, "rke2r$1")
	res = findAndReplaceKeyOn(res, "spec.k3sConfig", `^k3s`, "rke2")
	res = findAndReplaceKeyOn(res, "spec", `^k3s`, "rke2")
	res = findAndReplaceManagedFieldsKey(res, `f:k3s`, "f:rke2")

	return res
}

func Rke2RancherToK3sRancher(k3sLocalCluster unstructured.Unstructured) unstructured.Unstructured {
	res := findAndReplaceLabelsValue(k3sLocalCluster, `^rke2$`, "k3s")
	res = findAndReplaceAnnotationsValue(res, `rke2r([0-9]*)$`, "k3s$1")
	res = findAndReplaceStatusValue(k3sLocalCluster, `^rke2$`, "k3s")
	res = findAndReplaceValueOn(res, "spec.rke2Config", `rke2r([0-9]*)$`, "k3s$1")
	res = findAndReplaceKeyOn(res, "spec.rke2Config", `^rke2`, "k3s")
	res = findAndReplaceKeyOn(res, "spec", `^rke2`, "k3s")
	res = findAndReplaceManagedFieldsKey(res, `f:rke2`, "f:k3s")

	return res
}

func findAndReplaceLabelsValue(data unstructured.Unstructured, targetKey string, replacement string) unstructured.Unstructured {
	matchRegex := regexp.MustCompile(targetKey)

	labels := data.GetLabels()
	if val, ok := labels["provider.cattle.io"]; ok && matchRegex.MatchString(val) {
		labels["provider.cattle.io"] = matchRegex.ReplaceAllString(val, replacement)
	}
	data.SetLabels(labels)

	return data
}

func findAndReplaceAnnotationsValue(data unstructured.Unstructured, targetKey string, replacement string) unstructured.Unstructured {
	matchRegex := regexp.MustCompile(targetKey)

	annotations := data.GetAnnotations()
	if val, ok := annotations["management.cattle.io/current-cluster-controllers-version"]; ok && matchRegex.MatchString(val) {
		annotations["management.cattle.io/current-cluster-controllers-version"] = matchRegex.ReplaceAllString(val, replacement)
	}
	data.SetAnnotations(annotations)

	return data
}

func findAndReplaceStatusValue(data unstructured.Unstructured, targetKey string, replacement string) unstructured.Unstructured {
	matchRegex := regexp.MustCompile(targetKey)
	objectData := data.UnstructuredContent()

	status, found := objectData["status"].(map[string]interface{})
	if found {
		status["driver"] = "imported"
		if matchRegex.MatchString(status["provider"].(string)) {
			status["provider"] = matchRegex.ReplaceAllString(status["provider"].(string), replacement)
		}
	}
	version := status["version"].(map[string]interface{})
	currentPlatform := version["platform"]
	version = map[string]interface{}{
		"platform": currentPlatform,
	}
	status["version"] = version
	objectData["status"] = status
	data.SetUnstructuredContent(objectData)

	return data
}

func findAndReplaceKeyOn(localCluster unstructured.Unstructured, target string, targetKey string, replacement string) unstructured.Unstructured {
	targetObject, err := getTargetObject(localCluster.UnstructuredContent(), target)
	if err != nil {
		logrus.Fatalf("Something went wrong getting `%v`; %v", target, err)
	}

	targetObject = findAndReplaceKey(targetObject, targetKey, replacement)

	return localCluster
}

// This is highly opinionated and designed for this tools' implementation.
// DO NOT reuse this without understanding why it was made the way it is.
func getTargetObject(cluster map[string]interface{}, target string) (map[string]any, error) {
	if targetDepth := strings.Count(target, "."); targetDepth >= 1 {
		targetParts := strings.Split(target, ".")
		tempItem := cluster
		for _, part := range targetParts {
			newTempItem, err := getTargetObject(tempItem, part)
			if err != nil {
				return make(map[string]interface{}), err
			}
			tempItem = newTempItem
		}

		return tempItem, nil
	}

	targetObject := cluster[target]
	if targetObject == nil {
		return make(map[string]interface{}), fmt.Errorf("invalid target `%v` on resource", target)
	}

	return targetObject.(map[string]interface{}), nil
}

func findAndReplaceKey(data map[string]interface{}, targetKey string, replacement string) map[string]interface{} {
	matchRegex := regexp.MustCompile(targetKey)

	for key, value := range data {
		if matchRegex.MatchString(key) {
			newKey := matchRegex.ReplaceAllString(key, replacement)
			delete(data, key)
			data[newKey] = value
		}
	}

	return data
}

func findElement(array []metav1.ManagedFieldsEntry, search string) (*metav1.ManagedFieldsEntry, int) {
	for id, element := range array {
		if element.Manager == search { // check the condition if its true return index
			return &element, id
		}
	}

	return nil, 0
}

func findAndReplaceManagedFieldsKey(localCluster unstructured.Unstructured, targetKey string, replacement string) unstructured.Unstructured {
	managedFields := localCluster.GetManagedFields()
	rancherFields, _ := findElement(managedFields, "rancher")
	if rancherFields == nil {
		return localCluster
	}
	matchRegex := regexp.MustCompile(targetKey)

	fields := rancherFields.FieldsV1
	fieldsJson := fields.String()
	fieldsJson = matchRegex.ReplaceAllString(fieldsJson, replacement)
	err := fields.UnmarshalJSON([]byte(fieldsJson))
	if err != nil {
		return localCluster
	}
	localCluster.SetManagedFields(managedFields)

	return localCluster
}

func findAndReplaceValueOn(localCluster unstructured.Unstructured, target string, targetKey string, replacement string) unstructured.Unstructured {
	targetObject, err := getTargetObject(localCluster.UnstructuredContent(), target)
	if err != nil {
		logrus.Fatalf("Something went wrong getting `%v`; %v", target, err)
	}

	logrus.Infof("data out: %v", localCluster)
	updatedObject := findAndReplaceValue(targetObject, targetKey, replacement)
	logrus.Infof("updated out: %v", updatedObject)
	err = setTargetObject(&localCluster, target, updatedObject)
	if err != nil {
		logrus.Fatalf("Something went wrong getting `%v`; %v", target, err)
	}
	logrus.Infof("data out: %v", localCluster)

	return localCluster
}

func setTargetObject(cluster *unstructured.Unstructured, target string, object map[string]interface{}) error {
	unstructuredObj := cluster.UnstructuredContent()
	if targetDepth := strings.Count(target, "."); targetDepth >= 1 {
		targetParts := strings.Split(target, ".")

		// Traverse the map to reach the nested map where the update will happen
		currentMap := unstructuredObj
		for i, key := range targetParts {
			// If it's the last key, update the value
			if i == len(targetParts)-1 {
				currentMap[key] = object
				// Create a new unstructured object with the updated content
				cluster.SetUnstructuredContent(unstructuredObj)
				return nil
			}

			// Otherwise, move deeper into the map
			if nextMap, ok := currentMap[key].(map[string]interface{}); ok {
				currentMap = nextMap
			} else {
				// If the path does not exist, return an error
				return fmt.Errorf("path does not exist: %s", target)
			}
		}
		return nil
	}

	// non-nested update
	unstructuredObj[target] = object
	cluster.SetUnstructuredContent(unstructuredObj)

	return nil
}

func findAndReplaceValue(data map[string]any, targetKey string, replacement string) map[string]interface{} {
	matchRegex := regexp.MustCompile(targetKey)

	for key, value := range data {
		stringVal, ok := value.(string)
		if ok && matchRegex.MatchString(stringVal) {
			newVal := matchRegex.ReplaceAllString(stringVal, replacement)
			data[key] = newVal
		}
	}

	return data
}
