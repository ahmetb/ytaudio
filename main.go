package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {https://www.youtube.com/watch?v=hQApf_JdxQk
	listenAddr := os.Getenv("LISTEN_ADDR")
	addr := listenAddr + `:` + os.Getenv("PORT")
	http.HandleFunc("/watch", stream)
	http.HandleFunc("/view", viewHandler) //POST Reques
	log.Printf("starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------viewHandler------------------")

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		fmt.Println("-----------POST------------------")
		name := r.FormValue("name")
		name = strings.TrimPrefix(name, "https://www.youtube.com/")
		// /watch?v=dj_rapBEOgE
		http.Redirect(w, r, name, http.StatusFound)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//http://localhost:8080/watch?v=LWEBcN2o7Pc
func stream(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("v")
	if v == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "use format /watch?v=...")
		return
	}
	fmt.Println("v = -----------")
	fmt.Println(v)

	err := downloadVideoAndExtractAudio(v, w)
	if err != nil {
		log.Printf("stream error: %v", err)
		fmt.Fprintf(w, "stream error: %v", err)
		return
	}
}

func downloadVideoAndExtractAudio(id string, out io.Writer) error {
	url := fmt.Sprintf("https://youtube.com/watch?v=" + id)

	r, w := io.Pipe()
	defer r.Close()

	ytdl := exec.Command("youtube-dl", url, "-o-")
	ytdl.Stdout = w         // PIPE INPUT
	ytdl.Stderr = os.Stderr // show progress

	ffmpeg := exec.Command("ffmpeg", "-i", "/dev/stdin", "-f", "mp3",
		"-ab", "96000", "-vn", "-")
	ffmpeg.Stdin = r    // PIPE OUTPUT
	ffmpeg.Stdout = out // PIPE2 INPUT
	ffmpeg.Stderr = os.Stderr
	go func() {
		if err := ytdl.Run(); err != nil {
			log.Printf("WARN: ytdl error: %v", err)
		}
	}()
	err := ffmpeg.Run()
	log.Printf("stream finished")
	return err
}
