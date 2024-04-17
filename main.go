package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/imageseg-20191230/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/joho/godotenv"
)

// 声明一个全局的 logger 变量
var logger *log.Logger

func main() {

	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 创建日志文件
	logFile, err := os.OpenFile("project.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// 初始化全局 logger
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println("Starting application...")

	// 从环境变量中读取 Access Key ID 和 Access Key Secret
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	cred := &credentials.AccessKeyCredential{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}

	config := &openapi.Config{
		Credential: cred,
		RegionId:   String("cn-shanghai"),
		Endpoint:   String("imageseg.cn-shanghai.aliyuncs.com"),
	}

	cli, err := client.NewClient(config)
	if err != nil {
		panic(err)
	}

	inputDir := "input_images"
	outputDir := "output_images"
	ensureDir(inputDir)
	ensureDir(outputDir)

	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(info.Name()) {
			handleImage(cli, path, outputDir)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through input directory:", err)
	}
}

func ensureDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.Mkdir(dirName, 0755); err != nil {
			panic(err)
		}
	}
}

func isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext) == ".png" || strings.ToLower(ext) == ".jpg" || strings.ToLower(ext) == ".jpeg"
}

func handleImage(cli *client.Client, imagePath, outputDir string) {
	// 打开图像文件
	file, err := os.Open(imagePath)
	if err != nil {
		logger.Printf("打开图像文件失败：%s, 错误：%v\n", imagePath, err)
		return
	}
	defer file.Close()

	// 读取图像文件到缓冲区
	buf := new(bytes.Buffer)
	bytesWritten, err := io.Copy(buf, file)
	if err != nil {
		logger.Printf("读取图像到缓冲区失败：%s, 错误：%v\n", imagePath, err)
		return
	}
	if bytesWritten == 0 {
		logger.Printf("图像文件无数据：%s\n", imagePath)
		return
	}

	// 初始化运行时选项
	runtimeObject := new(util.RuntimeOptions).SetAutoretry(false).SetMaxIdleConns(3)

	// 创建请求对象，并设置图像数据
	request := new(client.SegmentBodyAdvanceRequest)
	request.SetImageURLObject(buf)

	// 调用 API 分割图像
	resp, err := cli.SegmentBodyAdvance(request, runtimeObject)
	if err != nil {
		logger.Printf("API 调用失败：%s, 错误：%v\n", imagePath, err)
		return
	}

	// 获取返回的图像 URL
	if resp.Body.Data.ImageURL == nil {
		logger.Printf("API 返回的图像 URL 为空：%s\n", imagePath)
		return
	}

	outputPath := filepath.Join(outputDir, filepath.Base(imagePath))

	// 下载并保存图像
	err = downloadAndSaveImage(*resp.Body.Data.ImageURL, outputPath)
	if err != nil {
		logger.Printf("下载或保存图像失败：%s, 错误：%v\n", imagePath, err)
	}
}

func downloadAndSaveImage(imageURL, outputPath string) error {
	// 发起 HTTP GET 请求下载图像
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("下载图像失败：%s, 错误：%v", imageURL, err)
	}
	defer resp.Body.Close()

	// 读取响应体数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取图像数据失败：%s, 错误：%v", imageURL, err)
	}

	// 将数据写入文件
	err = os.WriteFile(outputPath, imageData, 0644)
	if err != nil {
		return fmt.Errorf("保存图像到文件失败：%s, 错误：%v", imageURL, err)
	}

	return nil
}

func String(s string) *string {
	return &s
}
