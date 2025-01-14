package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type getFile struct {
	Key        string `json:"key"`
	Name       string `json:"Name"`
	EditTime   string `json:"EditTime"`
	Type       int    `json:"Type"`
	Size       int64  `json:"Size"`
	SuffixName string `json:"SuffixName"`
}

func getClient() (*cos.Client, string) {
	bucketName := viper.GetString("cos.bucketName")
	appid := viper.GetString("cos.appid")
	args := fmt.Sprintf("https://%s-%s.cos.ap-guangzhou.myqcloud.com", bucketName, appid)
	u, _ := url.Parse(args)
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  viper.GetString("cos.secretId"),
			SecretKey: viper.GetString("cos.secretKey"),
		},
	}), args
}

func UploadImage(filepath string, localpath string) (string, error) {
	client, args := getClient()
	_, _, err := client.Object.Upload(
		context.Background(), filepath, localpath, nil,
	)
	return args + filepath, err
}
func AddFolder(key string) error {
	client, _ := getClient()
	_, err := client.Object.Put(context.Background(), key, strings.NewReader(""), nil)
	return err
}
func UploadFile(key string, file io.Reader) error {
	client, _ := getClient()
	_, err := client.Object.Put(context.Background(), key, file, nil)
	return err
}
func ListObject(path string) (any, error) {
	client, _ := getClient()
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    path, // prefix 表示要查询的文件夹
		Delimiter: "/",  // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys:   1000, // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}
	isTruncated := true
	var files []getFile
	for isTruncated {
		opt.Marker = marker
		v, _, err := client.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			return nil, err
			break
		}
		log.Println("prefix:" + v.Prefix)
		log.Println("parent" + getParentPath(v.Prefix))
		if strings.Count(v.Prefix, "/") >= 3 {
			files = append(files, getFile{
				Key:        getParentPath(v.Prefix),
				Name:       "上一级目录",
				EditTime:   "",
				Type:       1,
				Size:       0,
				SuffixName: "",
			})
		}
		// common prefix 表示表示被 delimiter 截断的路径, 如 delimter 设置为/, common prefix 则表示所有子目录的路径
		for _, commonPrefix := range v.CommonPrefixes {
			fmt.Printf("CommonPrefixes: %v\n", commonPrefix)
			//添加文件夹
			files = append(files, getFile{
				Key:        commonPrefix,
				Name:       filepath.Base(commonPrefix),
				EditTime:   "",
				Type:       1,
				Size:       0,
				SuffixName: "",
			})
		}
		for _, content := range v.Contents {
			log.Println("key:" + content.Key)
			if content.Key == v.Prefix {
				// 创建一个新的 getFile 结构
				//newFile := getFile{
				//	Key:        getParentPath(content.Key),
				//	Name:       "上一级目录",
				//	EditTime:   "",
				//	Type:       1,
				//	Size:       content.Size,
				//	SuffixName: "",
				//}
				// 使用切片切割操作将新文件插入到切片开头
				//files = append([]getFile{newFile}, files...)
			} else {
				ext := filepath.Ext(content.Key)
				cleanedExt := strings.TrimPrefix(ext, ".")
				// cleanedExt 现在包含不带点号的文件扩展名
				files = append(files, getFile{
					Key:        content.Key,
					Name:       filepath.Base(content.Key),
					EditTime:   content.LastModified,
					Type:       0,
					Size:       content.Size / 1024,
					SuffixName: cleanedExt,
				})
			}
		}

		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}
	return files, nil
}
func getParentPath(path string) string {
	// 使用filepath包的Dir函数获取路径的上一级目录
	parentDir := filepath.Dir(path)
	parentDir = filepath.Dir(parentDir)
	if parentDir == "." {
		parentDir = ""
	} else {
		// 在window上会出现\，为了防止在不同平台的差异，采用ToSlash全换为“/”
		parentDir = filepath.ToSlash(parentDir) + "/"
	}
	log.Printf("path:" + path + " dir:" + parentDir)
	return parentDir
}
func FileDelete(filename string) error {
	client, _ := getClient()
	_, err := client.Object.Delete(context.Background(), filename)
	return err
}
func FolderDelete(dir string) error {
	client, _ := getClient()
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:  dir,
		MaxKeys: 1000,
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := client.Bucket.Get(context.Background(), opt)
		if err != nil {
			// Error
			return err
		}
		for _, content := range v.Contents {
			_, err = client.Object.Delete(context.Background(), content.Key)
			if err != nil {
				// Error
				return err
			}
		}
		isTruncated = v.IsTruncated
		marker = v.NextMarker
	}
	return nil
}
func GetUrl(filename string) string {
	client, _ := getClient()
	oUrl := client.Object.GetObjectURL(filename)
	return oUrl.String()
}
