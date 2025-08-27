package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// PackageConfig represents the configuration for each package
type PackageConfig struct {
	Name             string `json:"name"`
	TrafficCardCount int    `json:"traffice_card_count"`
	TrendCardCount   int    `json:"trend_card_count"`
}

const (
	data = `H4sIAAAAAAAA/7yRS07DMBCG9z2F5XWBOI7jJFtuACwrRX7MVBapbfWBWlW9O3JQhZUWJKAiS/+PzMx3nBFCCH0KYfWyVojOPKq1pR35EEbRhJ3f0o6w+eebVYcN7UiRPQEimG2PLsVpdOaBQWsbQIG6RBAVL7jQwIVl2Iq6UUDz9D4qPyaPC7oBb/vtIcKCdsUpt/2yfRmC3fRerSClp9teGNOvaUcEv6r0g9IwpKLnCMap4aJkHXZxcjDnLewnFxvgbezJfVGZ1/6aOQlqCefZykzaeZemYVLUZSE5Z7Kp7or7YvrRMXKa58zB2xsSN7KsjATUtba6EdgK1oKW2gosayj/SPwn7deIn3f9knd1e97sf3iLb3jPTu8BAAD//8OVyVPmAwAA`
)

func main() {
	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "decompress":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run main.go decompress <base64_data>")
				fmt.Println("Example: go run main.go decompress", data)
				return
			}

			// Decompress the provided base64 data
			compressedData := os.Args[2]
			result, err := UnCompress(compressedData)
			if err != nil {
				fmt.Printf("Error decompressing: %v\n", err)
				return
			}
			fmt.Println("Decompressed JSON:")
			fmt.Println(result)
			return

		default:
			fmt.Println("Unknown command. Available commands:")
			fmt.Println("  decompress <base64_data> - Decompress base64 encoded data")
			fmt.Println("  demo                    - Decompress the demo data")
			return
		}
	}

	// Default behavior: Generate SQL statements
	// Define the package configurations from the table
	packageConfigs := []PackageConfig{
		{Name: "agency-task-recruit_host-1host", TrafficCardCount: 1, TrendCardCount: 1},
		{Name: "agency-task-recruit_host-11host", TrafficCardCount: 2, TrendCardCount: 2},
		{Name: "agency-task-recruit_host-31host", TrafficCardCount: 2, TrendCardCount: 2},
		{Name: "agency-task-recruit_host-51host", TrafficCardCount: 5, TrendCardCount: 5},
		{Name: "agency-task-recruit_host-81host", TrafficCardCount: 5, TrendCardCount: 5},
	}

	// Debug: Print the first entry's JSON
	fmt.Println("=== First entry JSON ===")
	firstJSON := generatePackageJSON(packageConfigs[0])
	fmt.Println(firstJSON)
	fmt.Println()

	// Generate INSERT statements
	fmt.Println("-- INSERT statements")
	fmt.Println("INSERT INTO audio_reward.reward_package_cfg_info (`region`,`package_type`,`compressed`,`last_operator_id`,`package_name`,`package_json`) VALUES")

	for i, config := range packageConfigs {
		// Generate package JSON based on the configuration
		packageJSON := generatePackageJSON(config)

		// Compress the JSON
		compressedJSON, err := Compress(packageJSON)
		if err != nil {
			fmt.Printf("Error compressing JSON for %s: %v\n", config.Name, err)
			continue
		}

		// Format the SQL statement
		comma := ","
		if i == len(packageConfigs)-1 {
			comma = ";"
		}

		fmt.Printf("('TR',2,1,289,'%s', '%s')%s\n", config.Name, compressedJSON, comma)
	}

	fmt.Println()

	// Generate UPDATE statements
	fmt.Println("-- UPDATE statements")
	for _, config := range packageConfigs {
		// Generate package JSON based on the configuration
		packageJSON := generatePackageJSON(config)

		// Compress the JSON
		compressedJSON, err := Compress(packageJSON)
		if err != nil {
			fmt.Printf("Error compressing JSON for %s: %v\n", config.Name, err)
			continue
		}

		// Generate UPDATE statement
		fmt.Printf("UPDATE audio_reward.reward_package_cfg_info SET `package_json` = '%s' WHERE `package_name` = '%s';\n", compressedJSON, config.Name)
	}
}

// generatePackageJSON creates the JSON structure based on the package configuration
func generatePackageJSON(config PackageConfig) string {
	// Get current millisecond timestamp
	now := time.Now()
	milliseconds := now.UnixMilli()

	// Create uni strings with millisecond timestamp + 19 zeros
	trafficUni := fmt.Sprintf("%d-0.00000000000000000", milliseconds)
	trendUni := fmt.Sprintf("%d-0.00000000000000000", milliseconds+1) // Add 1 to make them different

	// Create a map with two items: RoomTrafficCard and RoomTrendCard
	goods := []interface{}{
		map[string]interface{}{
			"uni":              trafficUni,
			"goods_type":       53,
			"count":            config.TrafficCardCount,
			"group":            1,
			"level":            1,
			"days":             0,
			"package_type":     2,
			"index":            0,
			"pack_index":       0,
			"goods_name":       "RoomTrafficCard",
			"goods_type_label": "SpecialCard",
			"fid":              "pic/1e9d8ef5fb2fe543035be35d1f9568ae",
			"effect_fid":       "pic/1e9d8ef5fb2fe543035be35d1f9568ae",
			"expand":           "{\"send_type\":0}",
		},
		map[string]interface{}{
			"uni":              trendUni,
			"goods_type":       54,
			"count":            config.TrendCardCount,
			"group":            1,
			"level":            1,
			"days":             0,
			"package_type":     2,
			"index":            1,
			"pack_index":       0,
			"goods_name":       "RoomTrendCard",
			"goods_type_label": "SpecialCard",
			"fid":              "pic/c724c7efb6bdb85f9519eb7bd5f26e2e",
			"effect_fid":       "pic/c724c7efb6bdb85f9519eb7bd5f26e2e",
			"expand":           "{\"send_type\":0}",
		},
	}

	// Marshal to JSON with proper formatting
	jsonBytes, err := json.Marshal(goods)
	if err != nil {
		return "{}"
	}

	return string(jsonBytes)
}

func Compress(content string) (string, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(content))
	if err != nil {
		return content, err
	}

	if err = zw.Close(); err != nil {
		return content, err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func UnCompress(content string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return content, err
	}

	var buf = bytes.NewBuffer(b)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return content, err
	}

	var out bytes.Buffer
	if _, err = io.Copy(&out, zr); err != nil {
		return content, err
	}

	if err = zr.Close(); err != nil {
		return content, err
	}

	return out.String(), nil
}
