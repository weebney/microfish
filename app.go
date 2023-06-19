package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	wr "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Capture anchors and send to system browser
	wr.EventsOn(ctx, "sendToBrowser", func(optionalData ...interface{}) {
		url := fmt.Sprintf("%s", optionalData...)
		wr.BrowserOpenURL(ctx, url)
	})
	wr.EventsOn(ctx, "saveDoc", func(_ ...interface{}) {
		saveDocument(ctx)
	})
	wr.EventsOn(ctx, "openDoc", func(_ ...interface{}) {
		openDocument(ctx)
	})
}

func (a *App) GetPath() string {
	return path
}

func saveDocument(ctx context.Context) {
	opts := wr.SaveDialogOptions{
		DefaultDirectory:     tilde,
		DefaultFilename:      fileName,
		CanCreateDirectories: true,
	}
	savePath, err := wr.SaveFileDialog(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	CopyFile(path, savePath)
}

func openDocument(ctx context.Context) {
	opts := wr.OpenDialogOptions{
		DefaultDirectory: tilde,
		DefaultFilename:  fileName,
		Filters: []wr.FileFilter{
			{
				DisplayName: "Portable Document Files (*.pdf)",
				Pattern:     "*.pdf",
			},
		},
	}
	openPath, err := wr.OpenFileDialog(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if openPath == "" {
		return
	}

	OpenFile(openPath)
}

func CopyFile(originalPath string, destinationPath string) error {
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, originalFile)
	if err != nil {
		return err
	}

	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// TODO: THIS SHOULD BE DONE BY OPENING A NEW WINDOW, NOT FORKING
// HOPEFULLY POSSIBLE? TOO LAZY TO LOOK RN
func OpenFile(path string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(exe, abs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // detach
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	os.Exit(0)

	return nil
}
