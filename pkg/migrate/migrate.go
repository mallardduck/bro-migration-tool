package migrate

import (
	"regexp"
)

func K3sRancherToRke2Rancher(k3sLocalCluster map[string]interface{}) map[string]interface{} {
	// Replace "^k3s" with "^rke2" on both key and values
	res := findAndReplace(k3sLocalCluster, `^k3s`, "rke2")
	res = findAndReplace(res, `^f:k3s`, "f:rke2")
	res = findAndReplace(res, `k3s([0-9]*)$`, "rke2r$1")
	return res.(map[string]interface{})
}

func Rke2RancherToK3sRancher(k3sLocalCluster map[string]interface{}) map[string]interface{} {
	// Replace "^k3s" with "^rke2" on both key and values
	res := findAndReplace(k3sLocalCluster, `^rke2`, "k3s")
	res = findAndReplace(res, `^f:rke2`, "f:k3s")
	res = findAndReplace(res, `rke2r([0-9]*)$`, "k3s$1")
	return res.(map[string]interface{})
}

func findAndReplace(data interface{}, targetKey, replacement string) interface{} {
	matchK3sRegex := regexp.MustCompile(targetKey)
	switch val := data.(type) {
	case string:
		stringData := data.(string)
		if matchK3sRegex.MatchString(stringData) {
			stringData = matchK3sRegex.ReplaceAllString(stringData, replacement)
		}
		return stringData
	case map[string]interface{}:
		// TODO: implement
		for key, value := range val {
			var updatedKey string
			if matchK3sRegex.MatchString(key) {
				updatedKey = findAndReplace(key, targetKey, replacement).(string)
			} else {
				updatedKey = key
			}
			updatedValue := findAndReplace(value, targetKey, replacement)
			if key != updatedKey {
				delete(val, key)
			}
			val[updatedKey] = updatedValue
		}
		return val
	case []interface{}:
		for key, value := range val {
			updatedValue := findAndReplace(value, targetKey, replacement)
			val[key] = updatedValue
		}
		return val
	}

	return nil
}
