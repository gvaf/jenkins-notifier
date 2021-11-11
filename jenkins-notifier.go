package main
import (
  "log"
  "fmt"
  "flag"
  "os"
  "net"
  "strings"
  "io/ioutil"
  "path/filepath"
  "net/http"
)
 
var socketFile string
var outputDir string

func check(e error) {
    if e != nil {
	log.Fatal(e)
        panic(e)
    }
}


const validChars = "0123456789abcdef"

func isValidCommit(s string) bool {
   for _, char := range s {  
      if !strings.Contains(validChars, strings.ToLower(string(char))) {
         return false
      }
   }
   return true
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	check(err)
	commit := r.URL.Query().Get("commit")
	fileName := r.URL.Query().Get("fileName")
	fmt.Printf("commit = %s\n", commit)

	if isValidCommit(commit) {
		commitDir := filepath.Join(outputDir, commit)
		err = os.MkdirAll(commitDir, os.ModePerm)
		check(err)
		err = os.WriteFile(filepath.Join(commitDir, fileName), reqBody, 0644)
		check(err)
		w.Write([]byte(fmt.Sprintf("Received a POST request for commit = %s\n", commit)))
	} else {
		w.Write([]byte(fmt.Sprintf("Invalid commit = %s\n", commit)))
	}
}

func main() {
	flag.StringVar(&socketFile, "socket-file", "/tmp/jenkins-notifier.sock", "The socket file")
	flag.StringVar(&outputDir, "output-dir", "/tmp/", "The output directory")
	flag.Parse()

	if socketFile == "" || outputDir == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	os.Remove(socketFile)

        listener, err := net.Listen("unix", socketFile)
        if err != nil {
		log.Fatalf("Could not listen on %s: %v", socketFile, err)
		return
	}
        defer listener.Close()

	http.HandleFunc("/upload", uploadFile)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(fmt.Sprintf("Jenkins Notifier listening on %s", socketFile)))
	})
	fmt.Print("Jenkins Notifier listening from unix socket file '", socketFile, "'")

        if err = http.Serve(listener, nil); err != nil {
		log.Fatalf("Could not start Jenkins Nofifier server: %v", err)
	}

}
