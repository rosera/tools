// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"strings"
	"net/http"
	"os/exec"
	"runtime"
)

// CmdServe is the "claat serve ..." subcommand.
// addr is the hostname and port to bind the web server to.
// It returns a process exit code.
// rosera: Add a parameter containing the directory to serve
func CmdServe(addr string, serveDir string) int {
  if (serveDir == "."){
	  log.Printf("Serving codelabs from %s", serveDir)
    // Serve the current directory
	  http.Handle("/", http.FileServer(http.Dir(".")))
  } else {
	  log.Printf("Serving codelabs from %s", serveDir)
    // Serve the specified directory
	  http.Handle("/", http.FileServer(http.Dir(serveDir)))
  }
	ch := make(chan error, 1)
	go func() {
	  log.Printf("Serving codelabs on %s", addr)
		ch <- http.ListenAndServe(addr, nil)
	}()

  // rosera: Serve from a directory rather than root 
  if ContainsHttp(serveDir) {
    log.Println("The URL includes 'http'")
	  openBrowser(serveDir)
  } else {
    log.Println("The URL does not include 'http'")
	  openBrowser("http://" + addr + "/" + serveDir)
  }
  // rosera: Serve from a directory rather than root 
	// openBrowser("http://" + addr + "/" + serveDir)

  // rosera: Serve from a storage bucket
  //openBrowser("https://storage.googleapis.com/qwiklabs-lab-architect-rosera/labs/index.html")
  //openBrowser("https://drive.google.com/drive/folders/1PU64mu1Yvm023OKefdiEX4Jl5H6V15fp?usp=sharing")
	log.Fatalf("claat serve: %v", <-ch)
	return 0
}

func ContainsHttp(url string) bool {
  return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// openBrowser tries to open the URL in a browser.
func openBrowser(url string) error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	case "linux":
		args = []string{"xdg-open"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start()
}
