package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Config
var baseURL = "https://www.lmfdb.org"

// Collection info
var collections = map[string]string{
	"nf_fields":     "Number fields (数域)",
	"ec_curvedata":  "Elliptic curves (椭圆曲线)",
	"ec_classdata":  "Elliptic curve isogeny classes",
	"g2c_curves":    "Genus 2 curves",
	"char_dirichlet": "Dirichlet characters",
	"maass_newforms": "Maass forms",
	"mf_newforms":   "Modular forms",
	"lf_fields":     "Local fields",
	"artin":         "Artin representations",
}

// HTTP Client
var client = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// LMFDB API response
type APIResponse struct {
	Data []map[string]interface{} `json:"data"`
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "nf":
		queryNumberFields()
	case "ec":
		queryEllipticCurves()
	case "list":
		listCollections()
	case "version":
		fmt.Println("LMFDB CLI v1.0.0")
	case "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`LMFDB CLI - Query LMFDB from command line

Usage:
  lmfdb <command> [options]

Commands:
  nf                  Query Number Fields
  ec                  Query Elliptic Curves  
  list                List available API collections
  version             Show version information

Examples:
  lmfdb nf -d 2 -n 10
  lmfdb ec -r 2 -n 5
  lmfdb list

Note: LMFDB API requires browser verification.
      For production use, consider running a local LMFDB instance.`)
}

func queryNumberFields() {
	degree := "2"
	limit := "10"
	disc := ""

	// Parse args
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if i+1 < len(os.Args) {
			next := os.Args[i+1]
			if arg == "-d" || arg == "--degree" {
				degree = next
				i++
			} else if arg == "-n" || arg == "--limit" {
				limit = next
				i++
			} else if arg == "--disc" {
				disc = next
				i++
			}
		}
	}

	// Build URL
	url := fmt.Sprintf("%s/api/nf_fields/?_format=json&_limit=%s&degree=i%s", baseURL, limit, degree)
	if disc != "" {
		url += "&disc=i" + disc
	}

	data := queryAPI(url)
	if data == nil {
		fmt.Println("Error: Could not fetch data from LMFDB")
		fmt.Println("Note: LMFDB API may require browser verification (reCAPTCHA)")
		return
	}

	printNumberFields(data)
}

func queryEllipticCurves() {
	rank := ""
	torsion := ""
	limit := "10"

	// Parse args
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if i+1 < len(os.Args) {
			next := os.Args[i+1]
			if arg == "-r" || arg == "--rank" {
				rank = next
				i++
			} else if arg == "-t" || arg == "--torsion" {
				torsion = next
				i++
			} else if arg == "-n" || arg == "--limit" {
				limit = next
				i++
			}
		}
	}

	// Build URL
	url := fmt.Sprintf("%s/api/ec_curvedata/?_format=json&_limit=%s", baseURL, limit)
	if rank != "" {
		url += "&rank=i" + rank
	}
	if torsion != "" {
		url += "&torsion=i" + torsion
	}

	data := queryAPI(url)
	if data == nil {
		fmt.Println("Error: Could not fetch data from LMFDB")
		fmt.Println("Note: LMFDB API may require browser verification (reCAPTCHA)")
		return
	}

	printEllipticCurves(data)
}

func listCollections() {
	fmt.Println("\nAvailable API Collections:\n")
	for name, desc := range collections {
		fmt.Printf("  %-20s %s\n", name, desc)
	}
	fmt.Println("")
}

func queryAPI(url string) []map[string]interface{} {
	fmt.Printf("Querying: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return nil
	}

	// Check if we got HTML (reCAPTCHA) instead of JSON
	bodyStr := string(body)
	if strings.Contains(bodyStr, "recaptcha") || strings.Contains(bodyStr, "Checking your browser") {
		return nil
	}

	// Try to parse JSON
	var result APIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		// Try alternative format
		var altResult map[string]interface{}
		err = json.Unmarshal(body, &altResult)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			return nil
		}
		// Try to extract data
		if data, ok := altResult["data"].([]map[string]interface{}); ok {
			return data
		}
		return nil
	}

	return result.Data
}

func printNumberFields(data []map[string]interface{}) {
	if len(data) == 0 {
		fmt.Println("No results found")
		return
	}

	fmt.Printf("\nNumber Fields (showing %d results)\n\n", len(data))

	// Print header
	fmt.Printf("%-10s %-8s %-12s %-8s\n", "ID", "Degree", "Class #", "CM")
	fmt.Println(strings.Repeat("-", 45))

	// Print rows
	for i, item := range data {
		id := fmt.Sprintf("%d", i+1)
		degree := getString(item["degree"])
		classNum := getString(item["class_number"])
		cm := getString(item["cm"])
		
		fmt.Printf("%-10s %-8s %-12s %-8s\n", id, degree, classNum, cm)
	}
	fmt.Println("")
}

func printEllipticCurves(data []map[string]interface{}) {
	if len(data) == 0 {
		fmt.Println("No results found")
		return
	}

	fmt.Printf("\nElliptic Curves (showing %d results)\n\n", len(data))

	// Print header
	fmt.Printf("%-10s %-8s %-8s %-10s\n", "ID", "Rank", "Torsion", "Conductor")
	fmt.Println(strings.Repeat("-", 45))

	// Print rows
	for i, item := range data {
		id := fmt.Sprintf("%d", i+1)
		rank := getString(item["rank"])
		torsion := getString(item["torsion"])
		conductor := getString(item["conductor"])
		
		fmt.Printf("%-10s %-8s %-8s %-10s\n", id, rank, torsion, conductor)
	}
	fmt.Println("")
}

func getString(val interface{}) string {
	if val == nil {
		return "N/A"
	}
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
