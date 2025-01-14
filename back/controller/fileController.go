package controller

import (
	"fmt"
	"loginTest/api"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetFileList(c *gin.Context) {
	//db := common.GetDB()
	//从中间件存入的user中获取userID
	value, exisits := c.Get("user")
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	var user model.User
	user = value.(model.User)
	if user.UserID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	path := c.PostForm("path")
	// 强制规定只能在course文件夹内操作
	if !strings.HasPrefix(path, "course/") {
		// 如果不是以 "course/" 开头，添加它
		path = "course/" + path
	}
	fmt.Println("path:" + path)
	object, err := api.ListObject(path)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "文件列表获取失败")
		return
	}
	c.JSON(200, object)
}

// AddFolder 新增文件夹
func AddFolder(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件夹创建失败")
		return
	}
	_, exist := c.Get("user")
	if !exist {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	// 强制规定只能在course文件夹内操作
	if !strings.HasPrefix(path, "course/") {
		// 如果不是以 "course/" 开头，添加它
		path = "course/" + path
	}
	fmt.Println("path:" + path)
	filename := path
	err := api.AddFolder(filename)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "上传错误"+err.Error())
		return
	}
	// 返回成功
	response.Success(c, nil, "文件夹创建成功")
}
func DeleteFile(c *gin.Context) {
	//从中间件存入的user中获取userID
	_, exists := c.Get("user")
	if !exists {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	path := c.PostForm("path")
	fmt.Println("path:" + path)
	err := api.FileDelete(path)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "删除错误"+err.Error())
		return
	}
	// 返回成功
	response.Success(c, nil, "删除成功")
}

// DeleteFolder 删除文件夹
func DeleteFolder(c *gin.Context) {
	//从中间件存入的user中获取userID
	_, exists := c.Get("user")
	if !exists {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	path := c.PostForm("path")
	fmt.Println("path:" + path)
	err := api.FolderDelete(path)
	if err != nil {
		response.Fail(c, nil, "删除错误"+err.Error())
		return
	}
	response.Success(c, nil, "删除成功！")
}
func GetObjectUrl(c *gin.Context) {
	//从中间件存入的user中获取userID
	_, exisits := c.Get("user")
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	key := c.PostForm("filename")
	url := api.GetUrl(key)
	response.Response(c, http.StatusOK, 200, gin.H{"url": url}, "文件url获取成功")
}
