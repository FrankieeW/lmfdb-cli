package main

import (
	"crypto/tls"
	"encoding/csv"
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

var baseURL = "https://www.lmfdb.org"

var client = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func main() {
	nfCmd := flag.NewFlagSet("nf", flag.ExitOnError)
	ecCmd := flag.NewFlagSet("ec", flag.ExitOnError)
	
	// nf options
	nfDegree := nfCmd.Int("d", 2, "Number field degree")
	nfDisc := nfCmd.String("disc", "", "Filter by discriminant")
	nfClassNum := nfCmd.String("class", "", "Filter by class number")
	nfSignature := nfCmd.String("sig", "", "Filter by signature (e.g., 0,1)")
	nfLimit := nfCmd.Int("n", 10, "Number of results")
	nfOffset := nfCmd.Int("offset", 0, "Result offset")
	nfSort := nfCmd.String("sort", "", "Sort by field (use -field for descending)")
	nfFields := nfCmd.String("f", "", "Fields to return (comma-separated)")
	nfOutput := nfCmd.String("o", "", "Output file")
	nfFormat := nfCmd.String("fmt", "table", "Output format: table, json, csv")
	nfQuiet := nfCmd.Bool("q", false, "Quiet mode")
	nfID := nfCmd.String("id", "", "Get specific field by label")
	
	// ec options
	ecRank := ecCmd.String("r", "", "Filter by rank")
	ecTorsion := ecCmd.String("t", "", "Filter by torsion")
	ecConductor := ecCmd.String("conductor", "", "Filter by conductor")
	ecLimit := ecCmd.Int("n", 10, "Number of results")
	ecOffset := ecCmd.Int("offset", 0, "Result offset")
	ecSort := ecCmd.String("sort", "", "Sort by field")
	ecFields := ecCmd.String("f", "", "Fields to return (comma-separated)")
	ecOutput := ecCmd.String("o", "", "Output file")
	ecFormat := ecCmd.String("fmt", "table", "Output format: table, json, csv")
	ecQuiet := ecCmd.Bool("q", false, "Quiet mode")

	flag.Usage = printHelp
	flag.Parse()
	
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "nf":
		nfCmd.Parse(os.Args[2:])
		queryNumberFields(NumberFieldOptions{
			Degree:   *nfDegree,
			Disc:     *nfDisc,
			ClassNum: *nfClassNum,
			Signature: *nfSignature,
			Limit:    *nfLimit,
			Offset:   *nfOffset,
			Sort:     *nfSort,
			Fields:   *nfFields,
			Output:   *nfOutput,
			Format:   *nfFormat,
			Quiet:    *nfQuiet,
			ID:       *nfID,
		})
	case "ec":
		ecCmd.Parse(os.Args[2:])
		queryEllipticCurves(EllipticCurveOptions{
			Rank:      *ecRank,
			Torsion:   *ecTorsion,
			Conductor: *ecConductor,
			Limit:     *ecLimit,
			Offset:    *ecOffset,
			Sort:      *ecSort,
			Fields:    *ecFields,
			Output:    *ecOutput,
			Format:    *ecFormat,
			Quiet:     *ecQuiet,
		})
	case "list", "ls":
		listCollections()
	case "version", "v":
		fmt.Println("LMFDB CLI v1.2.0")
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

type NumberFieldOptions struct {
	Degree    int
	Disc      string
	ClassNum  string
	Signature string
	Limit     int
	Offset    int
	Sort      string
	Fields    string
	Output    string
	Format    string
	Quiet     bool
	ID        string
}

type EllipticCurveOptions struct {
	Rank      string
	Torsion   string
	Conductor string
	Limit     int
	Offset    int
	Sort      string
	Fields    string
	Output    string
	Format    string
	Quiet     bool
}

func printHelp() {
	fmt.Println(`LMFDB CLI v1.2.0 - Query LMFDB from command line

Usage:
  lmfdb <command> [options]

Commands:
  nf                  Query Number Fields
  ec                  Query Elliptic Curves  
  list (ls)           List available API collections
  version (v)         Show version information

Number Fields (nf):
  -d, --degree <n>    Number field degree (default: 2)
  -n, --limit <n>     Number of results (default: 10)
  --offset <n>        Result offset for pagination
  --sort <field>      Sort by field (use -field for descending)
  --disc <value>      Filter by discriminant
  --class <n>         Filter by class number
  --sig <r1,r2>       Filter by signature (e.g., "0,1")
  -f, --fields <list> Fields to return (comma-separated)
  -o, --output <file> Output file
  --fmt <format>      Output format: table, json, csv (default: table)
  --id <label>       Get specific field by label (e.g., 2.0.3.1)
  -q, --quiet         Quiet mode

  Available fields: label, degree, disc, class_number, class_group, 
                    signature, cm, coeffs, conductor, ramified_primes

Examples:
  lmfdb nf -d 2 -n 20
  lmfdb nf -d 3 --sort -class_number
  lmfdb nf -d 2 --disc -5
  lmfdb nf -d 2 -f label,degree,disc,class_number -o fields.json --fmt json
  lmfdb nf --id 2.0.3.1
  lmfdb nf -d 2 -n 100 --fmt csv -o data.csv

Elliptic Curves (ec):
  -r, --rank <n>       Filter by rank
  -t, --torsion <n>   Filter by torsion
  --conductor <n>     Filter by conductor
  -n, --limit <n>     Number of results (default: 10)
  --offset <n>        Result offset
  --sort <field>      Sort by field
  -f, --fields <list> Fields to return (comma-separated)
  -o, --output <file> Output file
  --fmt <format>      Output format: table, json, csv

  Available fields: label, conductor, rank, torsion, a1, a2, a3, a4, a6,
                    j_invariant, cid, ciso, abvar_eqn

Examples:
  lmfdb ec -r 0 -n 10
  lmfdb ec -t 5
  lmfdb ec --conductor 11a1
  lmfdb ec -r 0 -t 1 -n 50 --fmt json -o curves.json

Note: Some queries may be blocked by reCAPTCHA.`)
}

func queryNumberFields(opt NumberFieldOptions) {
	if !opt.Quiet {
		fmt.Println("Querying LMFDB...")
	}

	var url string
	
	if opt.ID != "" {
		url = fmt.Sprintf("%s/api/nf_fields/%s/?_format=json", baseURL, opt.ID)
	} else {
		url = fmt.Sprintf("%s/api/nf_fields/?_format=json&_limit=%d&degree=i%d", 
			baseURL, opt.Limit, opt.Degree)
		
		if opt.Offset > 0 {
			url += fmt.Sprintf("&_offset=%d", opt.Offset)
		}
		if opt.Sort != "" {
			url += "&_sort=" + opt.Sort
		}
		if opt.Disc != "" {
			url += "&disc=i" + opt.Disc
		}
		if opt.ClassNum != "" {
			url += "&class_number=i" + opt.ClassNum
		}
		if opt.Signature != "" {
			url += "&signature=li" + strings.Replace(opt.Signature, ",", ";", -1)
		}
		if opt.Fields != "" {
			url += "&_fields=" + opt.Fields
		}
	}

	if !opt.Quiet {
		fmt.Printf("  %s\n", url)
	}

	data := queryAPI(url)
	if data == nil {
		fmt.Println("Error: Could not fetch data from LMFDB")
		fmt.Println("Note: Some queries may be blocked by reCAPTCHA")
		return
	}

	// Handle single record
	if opt.ID != "" && len(data) == 1 {
		if opt.Output != "" {
			writeFormat(data, opt.Output, "json", opt.Quiet)
		} else {
			printRecordDetails(data[0], "Number Field")
		}
		return
	}

	// Output
	if opt.Output != "" {
		writeFormat(data, opt.Output, opt.Format, opt.Quiet)
	} else {
		printTable(data, opt.Format)
	}
}

func queryEllipticCurves(opt EllipticCurveOptions) {
	if !opt.Quiet {
		fmt.Println("Querying LMFDB...")
	}

	url := fmt.Sprintf("%s/api/ec_curvedata/?_format=json&_limit=%d", 
		baseURL, opt.Limit)
	
	if opt.Offset > 0 {
		url += fmt.Sprintf("&_offset=%d", opt.Offset)
	}
	if opt.Sort != "" {
		url += "&_sort=" + opt.Sort
	}
	if opt.Rank != "" {
		url += "&rank=i" + opt.Rank
	}
	if opt.Torsion != "" {
		url += "&torsion=i" + opt.Torsion
	}
	if opt.Conductor != "" {
		url += "&conductor=" + opt.Conductor
	}
	if opt.Fields != "" {
		url += "&_fields=" + opt.Fields
	}

	if !opt.Quiet {
		fmt.Printf("  %s\n", url)
	}

	data := queryAPI(url)
	if data == nil {
		fmt.Println("Error: Could not fetch data from LMFDB")
		return
	}

	if opt.Output != "" {
		writeFormat(data, opt.Output, opt.Format, opt.Quiet)
	} else {
		printTable(data, opt.Format)
	}
}

func listCollections() {
	collections := map[string]string{
		"nf_fields":       "Number fields",
		"ec_curvedata":    "Elliptic curves",
		"ec_classdata":    "Elliptic curve isogeny classes",
		"g2c_curves":     "Genus 2 curves",
		"char_dirichlet": "Dirichlet characters",
		"maass_newforms": "Maass forms",
		"mf_newforms":    "Modular forms",
		"lf_fields":      "Local fields",
		"artin":          "Artin representations",
		"belyi":          "Belyi maps",
	}
	
	fmt.Println("\n📚 Available API Collections:\n")
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

	bodyStr := string(body)
	if strings.Contains(bodyStr, "recaptcha") || strings.Contains(bodyStr, "Checking your browser") {
		fmt.Println("Error: Blocked by reCAPTCHA")
		return nil
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		var single map[string]interface{}
		err = json.Unmarshal(body, &single)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			return nil
		}
		return []map[string]interface{}{single}
	}

	return result.Data
}

func printTable(data []map[string]interface{}, format string) {
	if len(data) == 0 {
		fmt.Println("No results found")
		return
	}

	if format == "json" {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonData))
		return
	}

	if format == "csv" {
		writeCSV(os.Stdout, data)
		return
	}

	// Table format
	keys := make([]string, 0)
	for k := range data[0] {
		keys = append(keys, k)
	}
	// Limit to first 6 keys
	if len(keys) > 6 {
		keys = keys[:6]
	}

	fmt.Printf("\nResults (%d rows)\n\n", len(data))
	
	// Header
	for _, k := range keys {
		val := truncate(k, 14)
		fmt.Printf("%-15s ", val)
	}
	fmt.Println()
	for i := 0; i < len(keys); i++ {
		fmt.Print(strings.Repeat("-", 14) + " ")
	}
	fmt.Println()

	// Rows
	for _, item := range data {
		for _, k := range keys {
			val := truncate(formatValue(item[k]), 14)
			fmt.Printf("%-15s ", val)
		}
		fmt.Println()
	}
	fmt.Println("")
}

func writeFormat(data []map[string]interface{}, filename, format string, quiet bool) {
	var err error
	switch format {
	case "json":
		var jsonData []byte
		jsonData, err = json.MarshalIndent(data, "", "  ")
		if err == nil {
			err = os.WriteFile(filename, jsonData, 0644)
		}
	case "csv":
		var file *os.File
		file, err = os.Create(filename)
		if err == nil {
			writeCSV(file, data)
			file.Close()
		}
	default:
		err = fmt.Errorf("unsupported format: %s", format)
	}
	
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	} else if !quiet {
		fmt.Printf("✓ Results saved to %s\n", filename)
	}
}

func writeCSV(w *os.File, data []map[string]interface{}) {
	if len(data) == 0 {
		return
	}

	keys := make([]string, 0)
	for k := range data[0] {
		keys = append(keys, k)
	}

	csv := csv.NewWriter(w)
	csv.Write(keys)
	
	for _, item := range data {
		row := make([]string, len(keys))
		for i, k := range keys {
			row[i] = formatValue(item[k])
		}
		csv.Write(row)
	}
	csv.Flush()
}

func printRecordDetails(data map[string]interface{}, title string) {
	fmt.Printf("\n=== %s Details ===\n\n", title)
	
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
	if v == nil {
		return "N/A"
	}
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

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-2] + ".."
	}
	return s
}
