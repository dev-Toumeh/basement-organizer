package main

import (
	"basement/main/internal/database"
	"basement/main/internal/logg"
	"basement/main/internal/routes"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	LoadConfig()

	db := &database.DB{}
	db.Connect()
	defer db.Sql.Close()

	routes.RegisterRoutes(db)
	if err := templates.InitTemplates("./internal"); err != nil {
		logg.Fatal("Templates failed to initialize", err)
	}

	serverAddr := "localhost:8000"
	serverErrChan := runServer(serverAddr)

	// Allow server to start before launching Chrome
	time.Sleep(3 * time.Second)

	chromeCmd, chromeErrChan, err := launchChromeProcess(serverAddr)
	if err != nil {
		fmt.Println("Failed to launch Chrome:", err)
		os.Exit(1)
	}

	monitorProcesses(serverErrChan, chromeErrChan, chromeCmd)
}

// runServer starts the HTTP server in a goroutine and returns a channel for errors.
func runServer(addr string) chan error {
	errChan := make(chan error, 1)
	go func() {
		fmt.Println("Starting server on", addr)
		errChan <- http.ListenAndServe(addr, nil)
	}()
	return errChan
}

// launchChromeProcess launches Chrome in kiosk mode and returns its command, an error channel, and any startup error.
func launchChromeProcess(url string) (*exec.Cmd, chan error, error) {
	chromeArgs := []string{"--kiosk", "--disable-infobars", "--noerrdialogs", "http://" + url}
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("google-chrome", chromeArgs...)
	case "windows":
		cmd = exec.Command("chrome", chromeArgs...)
	default:
		return nil, nil, fmt.Errorf("unsupported OS")
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Wait()
	}()

	return cmd, errChan, nil
}

// monitorProcesses waits for either the server or Chrome process to exit, then shuts down accordingly.
func monitorProcesses(serverErrChan, chromeErrChan chan error, chromeCmd *exec.Cmd) {
	select {
	case err := <-chromeErrChan:
		if err != nil {
			fmt.Println("Chrome exited with error:", err)
		} else {
			fmt.Println("Chrome closed. Shutting down application.")
		}
		os.Exit(0)
	case err := <-serverErrChan:
		fmt.Println("Server stopped:", err)
		if chromeCmd.Process != nil {
			chromeCmd.Process.Kill()
		}
		os.Exit(0)
	}
}
