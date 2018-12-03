package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/minio/minio-go"
	logger "github.com/sirupsen/logrus"
)

var (
	client    *minio.Client
	endpoint  string
	accessKey string
	secretKey string
)

func InitClient() *minio.Client {
	endpoint = os.Getenv("MINIO_HOST")
	accessKey = os.Getenv("MINIO_ACCESS_KEY")
	secretKey = os.Getenv("MINIO_SECRET_KEY")

	minioClient, err := minio.New(endpoint, accessKey, secretKey, false)
	if err != nil {
		logger.Error("init minio connection with err = ", err)
	}
	client = minioClient
	return minioClient
}

func GetClientInstance() (client *minio.Client) {
	return
}

func UploadFile(filePath string) string {
	//fileName := "质检结果报表数据_20181127220427.xlsx"
	//filePath := "./质检结果报表数据_20181127220427.xlsx"
	//contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	fileName := getFileNameByPath(filePath)

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

	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName)
}

func getFileNameByPath(filePath string) string {
	strings.LastIndex(filePath, ".")

	f, err := os.Stat(filePath)

	if err != nil {
		logger.Error("File ", filePath, " does not exists.")
	}

	return f.Name()
}
