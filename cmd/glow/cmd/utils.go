package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/luaxlou/glow/pkg/api"
)

func request(method, path string, bodyData any, result any) error {
	var bodyReader io.Reader
	if bodyData != nil {
		b, _ := json.Marshal(bodyData)
		bodyReader = bytes.NewReader(b)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := strings.TrimSuffix(serverURL, "/") + path
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiResp api.Response
		// Try to decode error message
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err == nil && apiResp.Message != "" {
			return fmt.Errorf("server error: %s", apiResp.Message)
		}
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func printTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}
