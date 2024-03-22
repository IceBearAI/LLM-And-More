package util

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif")  //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	fileTypeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	fileTypeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")
	fileTypeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3")
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b0304140000000800", "zip")
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")
	fileTypeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")         //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")         //Real Audio (ram)
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// GetFileType 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

func Base64ToMultipartFile(base64Data string) (multipart.File, error) {
	// 解码 Base64 数据
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "base64-")
	if err != nil {
		return nil, err
	}
	defer func(tempFile *os.File) {
		_ = tempFile.Close()
	}(tempFile)

	// 将数据写入临时文件
	if _, err := tempFile.Write(imageBytes); err != nil {
		return nil, err
	}

	// 打开临时文件并返回 *os.File 类型
	multipartFile, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, err
	}

	return multipartFile, nil
}

func IsAudioFile(file multipart.File) (b bool, fileType string) {
	// 读取文件前 512 个字节
	head := make([]byte, 512)
	if _, err := file.Read(head); err != nil {
		return false, ""
	}
	// 将文件的读取位置重置回开始
	file.Seek(0, io.SeekStart)
	// 获取文件类型
	fileType = http.DetectContentType(head)
	return strings.HasPrefix(fileType, "audio"), fileType
}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

// probeOutputDuration 获取音频文件时长
func probeOutputDuration(a string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func GetAudioDuration(file multipart.File) (duration float64, err error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("/tmp", "upload-audio-*.tmp")
	if err != nil {
		return
	}
	defer os.Remove(tempFile.Name())
	defer func(tempFile *os.File) {
		_ = tempFile.Close()
	}(tempFile)
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return
	}
	// 获取音频时长
	var probe string
	probe, err = ffmpeg_go.Probe(tempFile.Name())
	if err != nil {
		return
	}
	duration, err = probeOutputDuration(probe)
	if err != nil {
		return
	}
	return duration, nil
}

func probeDuration(filePath string) (duration float64, err error) {
	// 获取音频时长
	var probe string
	probe, err = ffmpeg_go.Probe(filePath)
	if err != nil {
		return
	}
	duration, err = probeOutputDuration(probe)
	return
}

// ConvertAudioToMp3 将音频文件转换成mp3格式并读取时长
func ConvertAudioToMp3(file multipart.File) (convert multipart.File, duration float64, err error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("/tmp", "upload-audio-*.tmp")
	if err != nil {
		return
	}
	defer func(tempFile *os.File) {
		_ = tempFile.Close()
	}(tempFile)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return
	}
	// 将音频数据写入临时文件
	if _, err = tempFile.Write(fileBytes); err != nil {
		return
	}
	// 转换音频格式
	output := tempFile.Name() + ".mp3"
	if err = ffmpeg_go.Input(tempFile.Name()).
		Output(output, ffmpeg_go.KwArgs{"c:a": "libmp3lame", "b:a": "320k"}).
		OverWriteOutput().
		Run(); err != nil {
		return
	}
	duration, err = probeDuration(output)
	if err != nil {
		return
	}
	// 打开临时文件并返回 *os.File 类型
	convert, err = os.Open(output)
	if err != nil {
		return
	}
	return convert, duration, nil
}

// GetFileTypeName 获取文件类型名
func GetFileTypeName(file *os.File) (res string, err error) {
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Error reading the file")
		return
	}
	var fileTypeName string

	fileType := http.DetectContentType(buffer)
	switch fileType {
	case "application/pdf":
		fileTypeName = "pdf"
		break
	case "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		fileTypeName = "doc"
		break
	case "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		fileTypeName = "xls"
		break
	default:
		err = errors.Wrap(err, "")
	}
	return fileTypeName, nil
}

func Base64ToMultipartFileAndHeader(base64Data string, outFileName string, fileType string) (multipart.File, *multipart.FileHeader, error) {
	// 解码 Base64 数据
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, nil, err
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "base64-")
	if err != nil {
		return nil, nil, err
	}
	defer func(tempFile *os.File) {
		_ = tempFile.Close()
	}(tempFile)

	// 将数据写入临时文件
	if _, err := tempFile.Write(imageBytes); err != nil {
		return nil, nil, err
	}

	// 打开临时文件并返回 *os.File 类型
	multipartFile, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, nil, err
	}

	fileHeader := &multipart.FileHeader{ //构造一个head
		Filename: outFileName, //"out.png"
		Size:     int64(len(imageBytes)),
		Header: map[string][]string{
			"Content-Type": {fileType}, // 假设图片是PNG格式 "image/png"
		},
	}

	return multipartFile, fileHeader, nil
}
