package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	ConfigSource string
	WebPort      int
	WebPath      string
	AuthUser     string
	AuthPass     string
)

func init() {
	flag.StringVar(&ConfigSource, "c", "default", "config source default or env.")
	flag.IntVar(&WebPort, "port", 8080, "web port.")
	flag.StringVar(&WebPath, "path", "./www", "web path.")
	flag.StringVar(&AuthUser, "user", "", "authentication username")
	flag.StringVar(&AuthPass, "pass", "", "authentication password")
	flag.Parse()

	if ConfigSource == "env" {
		WebPort = getEnvInt("WEB_PORT", WebPort)
		WebPath = getEnvString("WEB_PATH", WebPath)
		AuthUser = getEnvString("AUTH_USER", AuthUser)
		AuthPass = getEnvString("AUTH_PASS", AuthPass)
	}
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 如果没有设置认证信息，则直接通过
		if AuthUser == "" || AuthPass == "" {
			next(w, r)
			return
		}

		// 从请求中获取认证信息
		user, pass, ok := r.BasicAuth()
		if !ok || user != AuthUser || pass != AuthPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="GoFileserver"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 认证通过，继续处理请求
		next(w, r)
	}
}

func main() {

	fmt.Printf("GoFileserver Port:%d Path:%s\n", WebPort, WebPath)

	// 确保 www 目录存在
	if err := os.MkdirAll(WebPath, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
	}

	// 注册处理函数，添加认证中间件
	http.HandleFunc("/", basicAuth(handleRequest))
	err := http.ListenAndServe(fmt.Sprintf(":%d", WebPort), nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// 处理文件上传
		handleUpload(w, r)
		return
	}
	// 处理文件下载（支持参数）
	handleDownload(w, r)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// 解析多部分表单
	r.ParseMultipartForm(10 << 20) // 10MB 限制

	// 获取上传的文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error retrieving file: %v\n", err)
		return
	}
	defer file.Close()

	// 获取额外参数
	dir := r.FormValue("dir")
	token := r.FormValue("token")

	// 构建保存路径
	savePath := WebPath
	if dir != "" {
		savePath = filepath.Join(savePath, dir)
		// 创建目录
		if err := os.MkdirAll(savePath, 0755); err != nil {
			fmt.Fprintf(w, "Error creating directory: %v\n", err)
			return
		}
	}
	savePath = filepath.Join(savePath, handler.Filename)
	
	// 创建目标文件
	dst, err := os.Create(savePath)
	if err != nil {
		fmt.Fprintf(w, "Error creating file: %v\n", err)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, file); err != nil {
		fmt.Fprintf(w, "Error saving file: %v\n", err)
		return
	}

	// 输出上传结果和参数
	fmt.Fprintf(w, "GoFileserver uploaded successfully: %s\n", handler.Filename)
	if dir != "" {
		fmt.Fprintf(w, "Directory: %s\n", dir)
	}
	if token != "" {
		fmt.Fprintf(w, "Token: %s\n", token)
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	token := r.URL.Query().Get("token")
	version := r.URL.Query().Get("version")

	// 提取文件路径（去掉查询参数）
	filePath := r.URL.Path
	if filePath == "/" {
		// 根路径，显示目录列表
		http.FileServer(http.Dir(WebPath)).ServeHTTP(w, r)
		return
	}

	// 构建完整文件路径
	fullPath := filepath.Join(WebPath, filePath[1:]) // 去掉开头的 /

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 输出下载信息和参数
	fmt.Printf("GoFileserver download request: %s\n", filePath)
	if token != "" {
		fmt.Printf("Token: %s\n", token)
	}
	if version != "" {
		fmt.Printf("Version: %s\n", version)
	}

	// 提供文件下载
	http.ServeFile(w, r, fullPath)
}

func getEnvString(name string, value string) string {
	ret := os.Getenv(name)
	if ret == "" {
		return value
	} else {
		return ret
	}
}

func getEnvInt(name string, value int) int {
	env := os.Getenv(name)
	if ret, err := strconv.Atoi(env); env == "" || err != nil {
		return value
	} else {
		return ret
	}
}
