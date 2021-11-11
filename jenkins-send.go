package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"path/filepath"
	"os"
	"io"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	var socketFile string
	var fileName string
	var commit string
	flag.StringVar(&socketFile, "socket-file", "/tmp/jenkins-notifier.sock", "The socket file")
	flag.StringVar(&commit, "commit", "", "Commit Id")
	flag.StringVar(&fileName, "file-name", "", "Send File to jenkins")
	flag.Parse()

	if socketFile == "" || commit == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketFile)
			},
		},
	}

	infile, err := os.Open(fileName)	
	check(err)
	defer infile.Close()
	var data io.Reader
	data = infile

	var response *http.Response
	url := "http://unix/upload?commit=" + commit + "&fileName=" + filepath.Base(fileName)
	response, err = httpc.Post(url, "application/octet-stream", data)
	check(err)
	io.Copy(os.Stdout, response.Body)
}
