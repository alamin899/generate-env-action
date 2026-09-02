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
	keys := make(map[string]bool)

	additionalKeysStr := os.Getenv("EXCLUDE_KEYS")
	if strings.TrimSpace(additionalKeysStr) != "" {
		for _, k := range strings.Split(additionalKeysStr, ",") {
			trimmedKey := strings.TrimSpace(k)
			if trimmedKey != "" {
				keys[trimmedKey] = true
			}
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

	// Determine output files
	envFilesStr := os.Getenv("GENERATE_ENV_FILES")
	if strings.TrimSpace(envFilesStr) != "" {
		files := strings.Split(envFilesStr, ",")
		for _, file := range files {
			fileName := strings.TrimSpace(file)
			if fileName == "" {
				continue
			}
			outFile, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Error creating %s file: %v\n", fileName, err)
				continue
			}
			for _, line := range newLines {
				outFile.WriteString(line)
			}
			outFile.Close()
		}
	}

	// Write to GITHUB_ENV if requested
	isSetProcessEnv := os.Getenv("IS_SET_PROCESS_ENV")
	if strings.ToLower(strings.TrimSpace(isSetProcessEnv)) == "true" {
		githubEnvFile := os.Getenv("GITHUB_ENV")
		if githubEnvFile != "" {
			f, err := os.OpenFile(githubEnvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening GITHUB_ENV file: %v\n", err)
			} else {
				defer f.Close()
				for key, val := range envData {
					if !excludeKeys[key] {
						valStr := fmt.Sprintf("%v", val)
						f.WriteString(fmt.Sprintf("%s<<EOF_ENV_GENERATOR\n%s\nEOF_ENV_GENERATOR\n", key, valStr))
					}
				}
			}
		}
	}
}
