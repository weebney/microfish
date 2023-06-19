package main

import (
	"embed"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wr "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/net/context"
)

/*
TODO:
- add logging
- add darkmode
- add config in general
- clean everything up
- PDF.js upgrade automation

- fix open file to open a new window instead of a new process
- fix livereloading
*/

//go:embed frontend/dist

var assets embed.FS

var (
	path  string
	tilde string
)

type loggingHandler struct {
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-store")
	h.handler.ServeHTTP(w, r)
}

var fileName string

func downloadRemoteFile(url *url.URL) (string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fileName = filepath.Base(url.Path)
	file, err := os.CreateTemp(os.TempDir(), fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func main() {
	app := NewApp()

	handleConfig()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var titleBar string
	var isUrl bool
	if len(os.Args) == 1 {
		path = ""
		titleBar = "microfish"
	} else {
		path = os.Args[1]
		titleBar = path
		if strings.HasPrefix(path, "~") {
			tilde, err = os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			path = filepath.Join(tilde, path[1:])
		}
		// relative files
		_, err := os.Stat(path)
		if err != nil {
			isUrl = true
		} else {
			path, err = filepath.Abs(path)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if isUrl {
		url, err := url.Parse(path)
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "https") {
			path, err = downloadRemoteFile(url)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// scope of available files from command line input
	var dir string
	if runtime.GOOS == "windows" {
		dir = filepath.VolumeName(pwd)
	} else {
		dir = "/"
	}
	fileServerHandler := http.FileServer(http.Dir(dir))
	handler := loggingHandler{handler: fileServerHandler}

	err = startFilewatcher(app.ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:            titleBar,
		WindowStartState: options.Maximised,
		// Menu:      AppMenu,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: handler,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
}

func startFilewatcher(context context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				wr.EventsEmit(context, "reloadPage")
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	return nil
}
