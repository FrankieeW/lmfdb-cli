package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config
var baseURL = "https://www.lmfdb.org"

// HTTP Client
var client = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func main() {
	// Main flags
	nfCmd := flag.NewFlagSet("nf", flag.ExitOnError)
	ecCmd := flag.NewFlagSet("ec", flag.ExitOnError)
	
	// nf options
	nfDegree := nfCmd.Int("d", 2, "Number field degree")
	nfDisc := nfCmd.String("disc", "", "Filter by discriminant")
	nfClassNum := nfCmd.String("class", "", "Filter by class number")
	nfLimit := nfCmd.Int("n", 10, "Number of results")
	nfFields := nfCmd.String("f", "", "Fields to return (comma-separated)")
	nfOutput := nfCmd.String("o", "", "Output file (JSON)")
	nfQuiet := nfCmd.Bool("q", false, "Quiet mode")
	nfID := nfCmd.String("id", "", "Get specific field by label")
	
	// ec options
	ecRank := ecCmd.String("r", "", "Filter by rank")
	ecTorsion := ecCmd.String("t", "", "Filter by torsion")
	ecConductor := ecCmd.String("conductor", "", "Filter by conductor")
	ecLimit := ecCmd.Int("n", 10, "Number of results")
	ecFields := ecCmd.String("f", "", "Fields to return (comma-separated)")
	ecOutput := ecCmd.String("o", "", "Output file (JSON)")
	ecQuiet := ecCmd.Bool("q", false, "Quiet mode")

	// Parse top-level flags first
	flag.Usage = printHelp
	flag.Parse()
	
	// Get command
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "nf":
		nfCmd.Parse(os.Args[2:])
		queryNumberFields(NumberFieldOptions{
			Degree:     *nfDegree,
			Disc:       *nfDisc,
			ClassNum:   *nfClassNum,
			Limit:      *nfLimit,
			Fields:     *nfFields,
			Output:     *nfOutput,
			Quiet:      *nfQuiet,
			ID:         *nfID,
		})
	case "ec":
		ecCmd.Parse(os.Args[2:])
		queryEllipticCurves(EllipticCurveOptions{
			Rank:       *ecRank,
			Torsion:    *ecTorsion,
			Conductor:  *ecConductor,
			Limit:      *ecLimit,
			Fields:     *ecFields,
			Output:     *ecOutput,
			Quiet:      *ecQuiet,
		})
	case "list":
		listCollections()
	case "collections":
		listCollections()
	case "version":
		fmt.Println("LMFDB CLI v1.1.0")
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

type NumberFieldOptions struct {
	Degree  int
	Disc    string
	ClassNum string
	Limit   int
	Fields  string
	Output  string
	Quiet   bool
	ID      string
}

type EllipticCurveOptions struct {
	Rank      string
	Torsion   string
	Conductor string
	Limit     int
	Fields    string
	Output    string
	Quiet     bool
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

Number Fields (nf) Options:
  -d, --degree <n>     Number field degree (default: 2)
  -n, --limit <n>     Number of results (default: 10)
  --disc <value>      Filter by discriminant
  --class <value>     Filter by class number
  -f, --fields <list> Fields to return (comma-separated)
  -o, --output <file> Output to JSON file
  --id <label>       Get specific field by label (e.g., 2.0.3.1)
  -q, --quiet         Quiet mode

Examples:
  lmfdb nf -d 2 -n 10
  lmfdb nf -d 3 --disc -23
  lmfdb nf --id 2.0.3.1
  lmfdb nf -d 2 -n 100 -o fields.json
  lmfdb nf -d 2 -f label,degree,disc

Elliptic Curves (ec) Options:
  -r, --rank <n>      Filter by rank
  -t, --torsion <n>   Filter by torsion
  --conductor <n>      Filter by conductor
  -n, --limit <n>     Number of results (default: 10)
  -f, --fields <list> Fields to return (comma-separated)
  -o, --output <file> Output to JSON file

Examples:
  lmfdb ec -r 2 -n 10
  lmfdb ec -t 5
  lmfdb ec --conductor 11

Note: LMFDB API may require browser verification for some queries.`)
}

func queryNumberFields(opt NumberFieldOptions) {
	if !opt.Quiet {
		fmt.Println("Querying LMFDB...")
	}

	// Build URL
	var url string
	
	// Special case: query by ID/label
	if opt.ID != "" {
		url = fmt.Sprintf("%s/api/nf_fields/%s/?_format=json", baseURL, opt.ID)
	} else {
		url = fmt.Sprintf("%s/api/nf_fields/?_format=json&_limit=%d&degree=i%d", 
			baseURL, opt.Limit, opt.Degree)
		
		if opt.Disc != "" {
			url += "&disc=i" + opt.Disc
		}
		if opt.ClassNum != "" {
			url += "&class_number=i" + opt.ClassNum
		}
		if opt.Fields != "" {
			url += "&_fields=" + opt.Fields
		}
	}

	fmt.Printf("  %s\n", url)

	data := queryAPI(url)
	if data == nil {
		fmt.Println("Error: Could not fetch data from LMFDB")
		fmt.Println("Note: LMFDB API may require browser verification (reCAPTCHA)")
		return
	}

	// Output
	if opt.Output != "" {
		writeJSON(data, opt.Output)
		if !opt.Quiet {
			fmt.Printf("Results saved to %s\n", opt.Output)
		}
	} else {
		printNumberFieldsTable(data)
	}
}

func queryEllipticCurves(opt EllipticCurveOptions) {
	if !opt.Quiet {
		fmt.Println("Querying LMFDB...")
	}

	// Build URL
	url := fmt.Sprintf("%s/api/ec_curvedata/?_format=json&_limit=%d", 
		baseURL, opt.Limit)
	
	if opt.Rank != "" {
		url += "&rank=i" + opt.Rank
	}
	if opt.Torsion != "" {
		url += "&torsion=i" + opt.Torsion
	}
	if opt.Conductor != "" {
		url += "&conductor=i" + opt.Conductor
	}
	if opt.Fields != "" {
		url += "&_fields=" + opt.Fields
	}

	fmt.Printf("  %s\n", url)

	data := queryAPI(url)
	if data == nil {
		fmt.Println("Error: Could not fetch data from LMFDB")
		fmt.Println("Note: LMFDB API may require browser verification (reCAPTCHA)")
		return
	}

	// Output
	if opt.Output != "" {
		writeJSON(data, opt.Output)
		if !opt.Quiet {
			fmt.Printf("Results saved to %s\n", opt.Output)
		}
	} else {
		printEllipticCurvesTable(data)
	}
}

func listCollections() {
	collections := map[string]string{
		"nf_fields":      "Number fields (数域)",
		"ec_curvedata":   "Elliptic curves (椭圆曲线)",
		"ec_classdata":   "Elliptic curve isogeny classes",
		"g2c_curves":    "Genus 2 curves",
		"char_dirichlet": "Dirichlet characters",
		"maass_newforms": "Maass forms",
		"mf_newforms":   "Modular forms",
		"lf_fields":     "Local fields",
		"artin":          "Artin representations",
	}
	
	fmt.Println("\nAvailable API Collections:\n")
	for name, desc := range collections {
		fmt.Printf("  %-20s %s\n", name, desc)
	}
	fmt.Println("")
}

func queryAPI(url string) []map[string]interface{} {
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
		fmt.Println("Error: Blocked by reCAPTCHA")
		return nil
	}

	// Try to parse JSON - handle both formats
	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		// Try alternative format: single object
		var single map[string]interface{}
		err = json.Unmarshal(body, &single)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			return nil
		}
		// Return as single-item array
		return []map[string]interface{}{single}
	}

	return result.Data
}

func printNumberFieldsTable(data []map[string]interface{}) {
	if len(data) == 0 {
		fmt.Println("No results found")
		return
	}

	// Check if single record (--id query)
	if len(data) == 1 {
		printRecordDetails(data[0], "Number Field")
		return
	}

	fmt.Printf("\nNumber Fields (showing %d results)\n\n", len(data))

	// Print header
	headers := []string{"#", "Label", "Degree", "Disc", "Class #", "CM"}
	printTableHeader(headers)
	
	// Print rows
	for i, item := range data {
		row := []string{
			fmt.Sprintf("%d", i+1),
			getString(item["label"]),
			getString(item["degree"]),
			getString(item["disc"]),
			getString(item["class_number"]),
			getString(item["cm"]),
		}
		printTableRow(row)
	}
	fmt.Println("")
}

func printEllipticCurvesTable(data []map[string]interface{}) {
	if len(data) == 0 {
		fmt.Println("No results found")
		return
	}

	// Check if single record
	if len(data) == 1 {
		printRecordDetails(data[0], "Elliptic Curve")
		return
	}

	fmt.Printf("\nElliptic Curves (showing %d results)\n\n", len(data))

	// Print header
	headers := []string{"#", "Label", "Conductor", "Rank", "Torsion"}
	printTableHeader(headers)
	
	// Print rows
	for i, item := range data {
		row := []string{
			fmt.Sprintf("%d", i+1),
			getString(item["label"]),
			getString(item["conductor"]),
			getString(item["rank"]),
			getString(item["torsion"]),
		}
		printTableRow(row)
	}
	fmt.Println("")
}

func printTableHeader(headers []string) {
	for _, h := range headers {
		fmt.Printf("%-15s ", h)
	}
	fmt.Println()
	for i := 0; i < len(headers); i++ {
		fmt.Print(strings.Repeat("-", 14) + " ")
	}
	fmt.Println()
}

func printTableRow(row []string) {
	for _, cell := range row {
		fmt.Printf("%-15s ", cell)
	}
	fmt.Println()
}

func printRecordDetails(data map[string]interface{}, title string) {
	fmt.Printf("\n=== %s Details ===\n\n", title)
	
	// Print all fields
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	
	for _, k := range keys {
		fmt.Printf("%-25s: %v\n", k, formatValue(data[k]))
	}
	fmt.Println("")
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return strconv.FormatFloat(val, 'f', 0, 64)
		}
		return fmt.Sprintf("%v", val)
	case []interface{}:
		return fmt.Sprintf("%v", val)
	case map[string]interface{}:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func writeJSON(data interface{}, filename string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
}

func getString(val interface{}) string {
	if val == nil {
		return "N/A"
	}
	switch v := val.(type) {
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatFloat(v, 'f', 0, 64)
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
