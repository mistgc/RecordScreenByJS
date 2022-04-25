package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"io"
	"io/ioutil"
)

func receive(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("content-type")
	contentLen := req.ContentLength

	fmt.Println(req.Context())

	fmt.Printf("[Server] Receving content-type:%s, content-length:%d\n", contentType, contentLen)
	if !strings.Contains(contentType, "multipart/form-data") {
		w.Write([]byte("content-type must be multipart/form-data"))
		return
	}

	if contentLen >= 8 * 1024 * 1024 {
		w.Write([]byte("file to large, limit 8MB"))
		fmt.Println("[Error] file to large, limit 4MB")
		return
	}

	// Parse request body into a FormData
	err := req.ParseMultipartForm(4*1024*1024)
	if err != nil {
		w.Write([]byte("ParseMultipartForm error:" + err.Error()))
		fmt.Println("[Error] ParseMultipartForm error:" + err.Error())
		return
	}

	if len(req.MultipartForm.File) <= 0 {
		w.Write([]byte("Not have any file."))
		fmt.Println("[Error] Not have any file.")
		return
	}

	fileName := req.Form.Get("filename")
	for name, files := range req.MultipartForm.File {
		fmt.Printf("[Server] req.MultipartForm.File, name=%s\n", name)

		if len(files) != 1 {
			w.Write([]byte("Too many files"))
			fmt.Println("[Server] Too many files")
			return
		}

		if name == "" {
			w.Write([]byte("Is not FileData"))
			fmt.Println("[Server] Is not FileData")
			return
		}

		for _, f := range files {
			handle, err := f.Open()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("unknown error, filename=%s, fileSize=%d, err:%s", f.Filename, f.Size, err.Error())))
				fmt.Printf("[Error] unknown error, filename=%s, fileSize=%d, err:%s\n", f.Filename, f.Size, err.Error())
				return
			}

			path := "./" + fileName
			_, err = os.Stat(path)
			if err == nil {
				fmt.Println("[Server] " + path + " is Exist.")
				//dst, _ := os.Open(path)
				dst, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend|os.ModePerm)
				//dst.Seek(0, io.SeekEnd)
				buf := make([]byte, 8 * 1024 * 1024)
				n, err := handle.Read(buf)
				fmt.Printf("[Server] Readed %d bytes.\n", n)
				if err != nil {
					fmt.Println("[Error] " + err.Error())
				}
				writenBytes, _ := dst.Write(buf[:n])
				fmt.Printf("[Server] Writen %d bytes.\n", writenBytes)
				dst.Close()
				fmt.Printf("[Server] Successful received, filename=%s, fileSize=%.2f MB, savePath=%s \n", f.Filename, float64(contentLen)/1024/1024, path)
			} else if os.IsNotExist(err){
				dst, _ := os.Create(path)
				io.Copy(dst, handle)
				fmt.Println("[Server] " + "Create " + path)
				//dst, _ := os.OpenFile(path, 0, 0666)
				//buf := make([]byte, 4 * 1024 *1024)
				//n, _ := handle.Read(buf)
				dst.Close()
				fmt.Printf("[Server] Successful received, filename=%s, fileSize=%.2f MB, savePath=%s \n", f.Filename, float64(contentLen)/1024/1024, path)
			}
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("../client/client.html")
	if err != nil {
		fmt.Printf("[Server] err:%s", err.Error())
		return
	}
	w.Write(f)
}

func main() {
	fmt.Println("[Server] Starting...")
	fmt.Println("[Server] Linstening on :8080 ")
	http.HandleFunc("/", index)
	http.HandleFunc("/vedioupload", receive)
	http.ListenAndServe(":8080", nil)
}
