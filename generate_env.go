package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func loadJSONEnv(varName string) map[string]interface{} {
	val := os.Getenv(varName)
	if val == "" {
		val = "{}"
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		fmt.Printf("Warning: Invalid JSON in %s. Defaulting to empty dict. Error: %v\n", varName, err)
		return make(map[string]interface{})
	}
	return data
}

func getMergedEnvData() map[string]interface{} {
	secrets := loadJSONEnv("SECRETS_CONTEXT")
	vars := loadJSONEnv("VARS_CONTEXT")

	envData := make(map[string]interface{})
	for k, v := range vars {
		envData[k] = v
	}
	for k, v := range secrets {
		envData[k] = v
	}
	return envData
}

func getExcludeKeys() map[string]bool {
	keys := map[string]bool{
		"github_token":              true,
		"DOCKER_PASSWORD":           true,
		"DOCKER_USERNAME":           true,
		"DOCKER_CONTAINER_REGISTRY": true,
		"DOCKER_REPOSITORY_NAME":    true,
		"VM_PUBLIC_IP":              true,
		"VM_USERNAME":               true,
		"SSH_KEY":                   true,
		"SSH_PORT":                  true,
		"HOST_PORT":                 true,
		"CONTAINER_RUNNING_PORT":    true,
		"DOCKER_NETWORK":            true,
		"JWT_PRIVATE_KEY":           true,
		"JWT_PUBLIC_KEY":            true,
	}

	additionalKeysStr := os.Getenv("EXCLUDE_KEYS")
	if strings.TrimSpace(additionalKeysStr) != "" {
		for _, k := range strings.Split(additionalKeysStr, ",") {
			keys[strings.TrimSpace(k)] = true
		}
	}
	return keys
}

func readLines(path string) []string {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		return lines
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}
	return lines
}

func formatEnvValue(val interface{}) string {
	valStr := fmt.Sprintf("%v", val)
	valStr = strings.ReplaceAll(valStr, "\"", "\\\"")
	valStr = strings.ReplaceAll(valStr, "\n", "\\n")
	return valStr
}

func main() {
	envData := getMergedEnvData()
	excludeKeys := getExcludeKeys()
	lines := readLines(".env.example")

	var newLines []string
	writtenKeys := make(map[string]bool)
	re := regexp.MustCompile(`^[#\s]*([A-Za-z0-9_]+)=`)

	// Update existing keys
	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			key := match[1]
			if val, exists := envData[key]; exists && !excludeKeys[key] {
				newLines = append(newLines, fmt.Sprintf("%s=\"%s\"\n", key, formatEnvValue(val)))
				writtenKeys[key] = true
				continue
			}
		}
		newLines = append(newLines, line)
	}

	// Append remaining new keys
	for key, val := range envData {
		if !writtenKeys[key] && !excludeKeys[key] {
			newLines = append(newLines, fmt.Sprintf("%s=\"%s\"\n", key, formatEnvValue(val)))
		}
	}

	outFile, err := os.Create(".env")
	if err != nil {
		fmt.Printf("Error creating .env file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	for _, line := range newLines {
		outFile.WriteString(line)
	}
}
