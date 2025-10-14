package jobs

import (
	"darkroom/pkg/config"
	"darkroom/pkg/utils"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	webview "github.com/webview/webview_go"
)

//go:embed static/submit/*
var embeddedFiles embed.FS

type FormData struct {
	Name    string `json:"name"`
	Script  string `json:"script"`
	CPU     int    `json:"cpu"`
	Memory  int    `json:"memory"`
	Image   string `json:"image"`
	JobType string `json:"jobtype"`
	Workers string `json:"workers"`
}

func SubmitGUI(cfg *config.Config, jobName, image, script, cpu, memory, jobType, workers string) (*FormData, error) {
	resultChan := make(chan *FormData, 1) // buffer 1 to avoid blocking

	defaults := map[string]interface{}{
		"image":   image,
		"jobtype": jobType,
		"workers": workers,
	}
	defaultsJSON, _ := json.Marshal(defaults)

	// Serve embedded static files
	fsys, err := fs.Sub(embeddedFiles, "static/submit")
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(fsys)))

	// Pick a free port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// Handle /submit but don’t send to channel here
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		var data FormData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		// Basic validation is performed already in the form before submission
		validName := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
		if !validName.MatchString(strings.ReplaceAll(data.Name, " ", "-")) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": fmt.Sprintf("invalid characters in name %q", data.Name),
			})
			return
		}

		// Just return OK to the browser
		w.WriteHeader(http.StatusOK)
	})

	// Start server in background
	go func() {
		http.Serve(listener, mux)
	}()

	// Open WebView window
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Submit job")
	w.SetSize(800, 680, webview.HintNone)

	url := fmt.Sprintf("http://localhost:%d/form.html?port=%d", port, port)
	w.Navigate(url)

	// Handle JS messages
	type Message struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data,omitempty"`
	}

	w.Bind("sendMessage", func(jsonMsg string) {
		var msg Message
		if err := json.Unmarshal([]byte(jsonMsg), &msg); err != nil {
			fmt.Println("Error parsing message:", err)
			return
		}

		switch msg.Type {
		case "submit":
			var data FormData
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				fmt.Println("Error parsing form data:", err)
				return
			}

			jobName = strings.ReplaceAll(data.Name, " ", "-")
			script = data.Script
			cpu = strconv.Itoa(data.CPU)
			memory = strconv.Itoa(data.Memory) + "Gi"

			err = utils.ValidateJobInputs(cfg, jobName, image, script, cpu, memory, jobType, workers)
			if err != nil {
				w.Dispatch(func() {
					msg := fmt.Sprintf(`showErrors({error: "%s"})`, err)
					w.Eval(msg)
				})
				return
			}

			name, err := SubmitJob(cfg, jobName, image, script, cpu, memory, jobType, workers)
			if err != nil {
				errJSON, _ := json.Marshal(err.Error())
				w.Dispatch(func() {
					msg := fmt.Sprintf(`showErrors({error: %s})`, errJSON)
					w.Eval(msg)
				})
				return
			}

			w.Dispatch(func() {
				msg := fmt.Sprintf(`showSuccessScreen("✅ Job '%s' submitted successfully!", "Script: %s | CPU: %d | Memory: %dGi")`,
					name, data.Script, data.CPU, data.Memory)
				w.Eval(msg)
			})

			select {
			case resultChan <- &data:
			default:
			}

		case "cancel", "close":
			select {
			case resultChan <- nil:
			default:
			}
			w.Terminate()

		case "onload":
			w.Dispatch(func() {
				w.Eval(fmt.Sprintf(`setDefaults(%s);`, defaultsJSON))
			})
		}
	})

	w.Run()

	data := <-resultChan
	return data, nil
}
