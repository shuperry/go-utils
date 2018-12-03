package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/minio/minio-go"
	logger "github.com/sirupsen/logrus"
)

var (
	client           *minio.Client
	uploadEndpoint   string
	downloadEndpoint string
	accessKey        string
	secretKey        string
)

type UploadUtil struct{}

func (uploadUtil *UploadUtil) InitClient() *minio.Client {
	uploadEndpoint = os.Getenv("MINIO_UPLOAD_HOST")
	downloadEndpoint = os.Getenv("MINIO_DOWNLOAD_HOST")
	accessKey = os.Getenv("MINIO_ACCESS_KEY")
	secretKey = os.Getenv("MINIO_SECRET_KEY")

	minioClient, err := minio.New(uploadEndpoint, accessKey, secretKey, false)
	if err != nil {
		logger.Error("init minio connection with err = ", err)
	}
	client = minioClient
	return minioClient
}

func (uploadUtil *UploadUtil) GetClientInstance() (client *minio.Client) {
	return
}

func (uploadUtil *UploadUtil) UploadFile(filePath string) string {
	//fileName := "质检结果报表数据_20181127220427.xlsx"
	//filePath := "./质检结果报表数据_20181127220427.xlsx"
	//contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	fileName := uploadUtil.getFileNameByPath(filePath)

	objectPrefix := os.Getenv("MINIO_OBJECT_PREFIX")
	objectName := fmt.Sprintf("%s/%s", objectPrefix, fileName)

	// Upload the file with FPutObject.
	size, err := client.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{})

	if err != nil {
		logger.Error("Failed upload ", objectName, " with err = ", err)
	} else {
		logger.Printf("Successfully uploaded %s of size %d\n", objectName, size)
	}

	// 上传成功或失败后均删除文件.
	go os.Remove(filePath)

	// 返回文件下载链接.
	return fmt.Sprintf("http://%s/%s/%s", downloadEndpoint, bucketName, objectName)
}

func (uploadUtil *UploadUtil) getFileNameByPath(filePath string) string {
	strings.LastIndex(filePath, ".")

	f, err := os.Stat(filePath)

	if err != nil {
		logger.Error("File ", filePath, " does not exists.")
	}

	return f.Name()
}
