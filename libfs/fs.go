package libfs

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

// InitializeFs 以./dist文件夹初始化文件服务器
func InitializeFs(port int) {
	dirFS := os.DirFS("./dist")
	fs := http.FileServerFS(dirFS)
	err := http.ListenAndServe(":"+strconv.Itoa(port), fs)
	if err != nil {
		log.Fatal("Error create file server: ", err)
	}
}
