package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

func ReadFileByPath(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return ioutil.ReadAll(f)
}

func GetFilePath(fileName string, uploadPath string) (filePath string, newFileName string) {
	fileSuffix := path.Ext(fileName)
	rd := RandomInt(100, 999)
	rdStr := strconv.Itoa(rd)
	newFileName = Md5Str(fileName+rdStr) + fileSuffix
	filePath = uploadPath + "/" + newFileName[0:2] + "/" + newFileName[2:4] + "/"
	return
}

func MakeFilePath(path string) (err error) {
	_, err = os.Stat(path)
	if err == nil {
		return
	}
	if os.IsExist(err) {
		return
	}
	err = os.MkdirAll(path, 0777)
	_ = os.Chmod(path, 0777)
	return
}

func RandomInt(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(end - start)
	random = start + random
	return random
}

func ReadFileAndWrite(file io.Writer, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	io.Copy(file, f)
}

type FileItem struct {
	FiledName string
	FiledPath string
	IsDir     bool
}

type ParamItem struct {
	FiledName  string
	FiledValue interface{}
	FiledType  string
}

type HeaderItem struct {
	FiledName  string
	FiledValue string
}

const (
	FiledTypeString  = "string"
	FiledTypeFloat64 = "float64"
	FiledTypeInt     = "int"
	FiledTypeBool    = "bool"
)

func UploadFile(url string, filePaths []FileItem, params []ParamItem, headers []HeaderItem) (res *http.Response, err error) {

	// 创建一个Buffer，用于存储POST请求体
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加文件
	if len(filePaths) > 0 {
		for _, v := range filePaths {
			f, _ := writer.CreateFormFile(v.FiledName, filepath.Base(v.FiledPath))
			ReadFileAndWrite(f, v.FiledPath)
		}
	}

	//添加参数
	if len(params) > 0 {
		for _, v := range params {
			field, _ := writer.CreateFormField(v.FiledName)
			// 默认字符串，其它类型作断言处理
			switch v.FiledType {
			case FiledTypeFloat64:
				field.Write([]byte(fmt.Sprintf("%f", v.FiledValue.(float64))))
			case FiledTypeInt:
				field.Write([]byte(fmt.Sprintf("%d", v.FiledValue.(int))))
			case FiledTypeBool:
				field.Write([]byte(fmt.Sprintf("%t", v.FiledValue.(bool))))
			default:
				field.Write([]byte(v.FiledValue.(string)))
			}
		}
	}

	writer.Close()

	// 发送POST请求到Python接口
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	//添加header 参数
	if (len(headers)) > 0 {
		for _, v := range headers {
			request.Header.Set(v.FiledName, v.FiledValue)
		}
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	//defer response.Body.Close()

	return response, nil
}

func ReadFile(filePath string) (multipart.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	//defer f.Close()

	return f, nil
}

// WriteAudioFile 转存文件
func WriteAudioFile(audioStr, fileDir, suffix string) (filepath string, err error) {
	//创建目录
	err = MakeFilePath(fileDir)
	if err != nil {
		return "", err
	}
	// 将base64字符串解码
	data, err := base64.StdEncoding.DecodeString(audioStr)
	if err != nil {
		return "", err
	}
	// 将解码后的字符串写入文件
	fileName := Krand(10, 1)
	filepath = fmt.Sprintf("%s/%s.%s", fileDir, fileName, suffix)

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return "", err
	}
	return
}

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// WriteFile 文件临时存储
func WriteFile(file multipart.File, handler *multipart.FileHeader, fileDir string) (filePath, fileName string, err error) {
	// 创建新文件并将内容写入其中
	filePath, fileName = GetFilePath(handler.Filename, fileDir)
	err = MakeFilePath(filePath)
	if err != nil {
		return "", "", err
	}
	filePathName := filePath + fileName

	newFile, err := os.Create(filePathName)
	if err != nil {
		return "", "", err
	}
	defer newFile.Close()

	io.Copy(newFile, file)

	return filePathName, fileName, nil
}

// File 实现 multipart.File 接口所需的方法
type File struct {
	*bytes.Reader
}

func (f *File) Close() error {
	return nil // bytes.Reader 不需要关闭资源，所以这里返回 nil 即可
}

// NewFile 创建一个新的 File 实例，该实例满足 multipart.File 接口
func NewFile(data []byte) *File {
	return &File{
		bytes.NewReader(data),
	}
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetHttpFileBody 获取http文件内容
func GetHttpFileBody(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, "http.Get")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "io.ReadAll")
		return
	}
	return
}

type FileInfo struct {
	//os.FileInfo
	ContentType string
	Name        string
	Size        int64     // length in bytes for regular files; system-dependent for others
	ModTime     time.Time // modification time
	IsDir       bool      // abbreviation for Mode().IsDir()
}

// GetFileList 取读目录下的文件列表
func GetFileList(dirPth string) (files []FileInfo, err error) {
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		info, err := fi.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(dirPth, info.Name())
		fileInfo := FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
		// 判断是否为软链，如果是软链获取真实文件路径
		if !info.IsDir() {
			symlink, _ := IsSymlink(fullPath)
			if symlink {
				readlink, rlErr := os.Readlink(fullPath)
				if rlErr != nil {
					continue
				}
				fullPath = path.Join(dirPth, readlink)
				f, err := os.Open(fullPath)
				if err != nil {
					continue
				}
				// 读取文件的前 512 个字节
				buffer := make([]byte, 512)
				if _, rErr := f.Read(buffer); rErr != nil {
					continue
				}
				// 使用 net/http 包的 DetectContentType 函数来获取文件类型
				info, _ = f.Stat()
				_ = f.Close()
				fileInfo.Size = info.Size()
				fileInfo.ModTime = info.ModTime()
				fileInfo.IsDir = info.IsDir()
				fileInfo.ContentType = http.DetectContentType(buffer)
			}
		}
		files = append(files, fileInfo)
	}
	return files, nil
}

// IsSymlink 判断给定路径是否为符号链接
func IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	// 检查文件模式位是否设置了 os.ModeSymlink
	return info.Mode()&os.ModeSymlink != 0, nil
}

// DirectoryExists 判断给定路径是否为存在的目录
func DirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
