package main

import (
	"fmt"
	"log"
	"loginTest/common"
	"loginTest/config"
	"loginTest/controller"
	"loginTest/heat"
	"loginTest/route"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"github.com/spf13/viper"

	// "loginTest/heat"
	"net/http"
	"os/exec"
	"time"
)

func Copy() {
	// 数据库连接信息
	dbHost := viper.GetString("datasource.host")
	dbPort := viper.GetInt("datasource.port")
	dbUser := viper.GetString("datasource.username")
	dbPassword := viper.GetString("datasource.password")
	dbName := viper.GetString("datasource.database")

	// 备份目录
	backupDir := "/app/database"

	c := cron.New()
	c.AddFunc("@every 12h", func() {
		backupFile := fmt.Sprintf("%s/backup_%s.sql", backupDir, time.Now().Format("2006-01-02 15:04:05"))
		cmd := exec.Command("mysqldump", fmt.Sprintf("-h%s", dbHost), fmt.Sprintf("-P%d", dbPort), fmt.Sprintf("-u%s", dbUser), fmt.Sprintf("-p%s", dbPassword), dbName, "--result-file="+backupFile)
		err := cmd.Run()
		if err != nil {
			log.Println("备份失败:", err)
			return
		}
		log.Println("备份成功:", backupFile)
	})
	c.AddFunc("@every 1h", controller.CalculateAndSaveScores)
	c.Start()
}

var r *gin.Engine

func main() {
	config.InitConfig()
	go Copy()
	db := common.InitDB()
	rds := common.RedisInit()
	common.InitJWTkey()
	defer rds.Close()
	defer db.Close()
	r = gin.Default()
	// 使用 http.FileServer 文件服务器处理 "/uploads/" 开头的请求，
	// 文件服务器获取文件的位置在 "./public" 文件夹下。

	r.StaticFS("/uploads", http.Dir("./public/uploads"))
	r.StaticFS("/resized", http.Dir("./public/resized"))

	// heat.StartHeat()

	route.CollectRoute(r)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// crtFile := "./ssl/1_ssemarket.cn_bundle.crt"
	// keyFile := "./ssl/2_ssemarket.cn.key"
	go heat.StartHeat()
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// if err := srv.ListenAndServeTLS(crtFile, keyFile); err != nil {
	// 	log.Fatal("ListenAndServeTLS: ", err)
	// }

	// if err :=  ; err != nil {
	//     log.Fatal("ListenAndServeTLS: ", err)
	// }

	log.Printf("Server started on port 8080")
	select {}
}

// package main

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	password := "123456" // 明文密码

// 	// 生成哈希密码
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Println("生成哈希密码失败:", err)
// 		return
// 	}

//		fmt.Println("哈希密码:", string(hashedPassword))
//	}
// package main

// import (
// 	"bytes"
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"encoding/base64"
// 	"fmt"
// 	"io"
// )

// // AES 加密
// func AESEncrypt(plaintext string, key []byte) (string, error) {
// 	// 创建 AES 加密块
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	// 对明文进行 PKCS7 填充
// 	plaintextBytes := []byte(plaintext)
// 	plaintextBytes = PKCS7Padding(plaintextBytes, aes.BlockSize)

// 	// 创建 CBC 加密模式
// 	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
// 	iv := ciphertext[:aes.BlockSize] // 初始化向量
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return "", err
// 	}

// 	mode := cipher.NewCBCEncrypter(block, iv)
// 	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintextBytes)

// 	// 返回 Base64 编码的密文
// 	return base64.StdEncoding.EncodeToString(ciphertext), nil
// }

// // PKCS7 填充
// func PKCS7Padding(data []byte, blockSize int) []byte {
// 	padding := blockSize - len(data)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(data, padtext...)
// }

// func main() {
// 	key := []byte("0123456789abcdef") // 16 字节密钥
// 	plaintext := "123456"

// 	ciphertext, err := AESEncrypt(plaintext, key)
// 	if err != nil {
// 		fmt.Println("加密失败:", err)
// 		return
// 	}

// 	fmt.Println("密文:", ciphertext)
// }
