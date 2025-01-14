package controller

import (
	"context"
	"fmt"
	"image"
	"log"
	"loginTest/api"
	"loginTest/common"
	"loginTest/dto"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	VALIDATESUCCESS = 0
	EMAILERROR      = 1
	EMAILNOTEXIST   = 2
	VALIDATEFAILED  = 3
)

func GetTokenUserID(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || len(tokenString) <= 7 || !strings.HasPrefix(tokenString, "Bearer ") {
		return 0
	}
	tokenString = tokenString[7:]
	token, claims, err := common.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return 0
	}

	// 获取token中的用户标识符
	tokenUserID := claims.UserID
	return tokenUserID
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("phone = ?", telephone).First(&user)
	if user.UserID != 0 {
		return true
	}
	return false
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.UserID != 0 {
		return true
	}
	return false
}

//func isNumExist(db *gorm.DB, num int) bool {
//	var user model.User
//	db.Where("num = ?", num).First(&user)
//	if user.UserID != 0 {
//		return true
//	}
//	return false
//}

type registerUser struct {
	Name      string `gorm:"type:varchar(20);not null"`
	Phone     string `gorm:"type:varchar(11);not null"`
	Password  string `gorm:"size:255;not null"`
	Password2 string `gorm:"size:255;not null"`
	Email     string `gorm:"type:varchar(11);not null"`
	//Num       string `gorm:"type:varchar(8);not null"`
	ValiCode string `gorm:"type:varchar(10);not null"`
	CDKey    string `json:"CDKey"`
}

type identity struct {
	Email    string `gorm:"type:varchar(50);not null"`
	ValiCode string `gorm:"type:varchar(10);not null"`
}

type ModifyUser struct {
	Email     string `gorm:"type:varchar(50);not null"`
	Phone     string `gorm:"type:varchar(11);not null"`
	Password  string `gorm:"size:255;not null"`
	Password2 string `gorm:"size:255;not null"`
	ValiCode  string
}

type requestEmail struct {
	Email string `gorm:"type:varchar(50);not null"`
	Mode  int
}

func DeleteMe(c *gin.Context) {
	db := common.GetDB()

	var requestUser = registerUser{}
	c.Bind(&requestUser)
	phone := requestUser.Phone
	email := requestUser.Email
	valiCode := requestUser.ValiCode

	var checkUser model.User
	db.Where("phone = ?", phone).First(&checkUser)
	userID := checkUser.UserID
	if checkUser.Email != email {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "请用注册时使用的email完成注销")
		return
	}
	if userID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "未找到该用户")
		return
	}

	resCode := IdentityValidateService(email, valiCode)
	if resCode == VALIDATEFAILED {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	} else if resCode == EMAILERROR {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	} else if resCode == EMAILNOTEXIST {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
		return
	}

	checkUser.Name = "用户已注销"
	checkUser.Phone = "0"
	checkUser.Email = ""
	checkUser.AvatarURL = ""

	db.Save(&checkUser)

	response.Response(c, http.StatusOK, 200, nil, "成功删除该用户")
}

//	func Register(c *gin.Context) {
//		// 连接数据库
//		db := common.GetDB()
//		// 从前端获取信息
//		var requestUser = registerUser{}
//		c.Bind(&requestUser)
//		//获取参数
//		name := requestUser.Name
//		telephone := requestUser.Phone
//		password1 := requestUser.Password
//		password2 := requestUser.Password2
//		password1 = util.Decrypt(password1)
//		password2 = util.Decrypt(password2)
//		email := strings.TrimSpace(requestUser.Email)
//		print(email)
//		//num := requestUser.Num
//		valiCode := requestUser.ValiCode
//		cdkey := requestUser.CDKey
//		//若使用postman等工具，写法如下：
//		// name := c.PostForm("name")
//		// telephone := c.PostForm("telephone")
//		// password := c.PostForm("password")
//		// PostForm中的参数与在postman中发送信息的名字要一致
//		if len(telephone) == 0 {
//			for {
//				telephone = util.GenerateRandomDigits(11)
//				// 如果存在，继续生成新的随机字符串
//				// 如果不存在，退出循环
//				if !isTelephoneExist(db, telephone) {
//					break
//				}
//			}
//		}
//		//验证数据
//		if len(telephone) != 11 {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "手机号必须为11位!!!")
//			println(telephone)
//			return
//		}
//		if len(password1) < 6 {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
//			return
//		}
//		if len(password2) < 6 {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
//			return
//		}
//		if password1 != password2 {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "两次密码不相同，请重新输入")
//			return
//		}
//		//如果名称没有传，给一个10位的随机字符串
//		if len(name) == 0 {
//			name = gofakeit.Name()
//		}
//		// 判断手机号是否存在
//		if isTelephoneExist(db, telephone) {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该手机号已存在")
//			return
//		}
//		check := false
//		for i := 0; i < len(email); i++ {
//			if email[i] == '@' {
//				check = true
//			}
//		}
//		if !check {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
//			return
//		}
//		if isEmailExist(db, email) {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱已存在")
//			return
//		}
//		rds := common.MyRedis
//		ctx := context.Background()
//		correctValiCode, err := rds.Get(ctx, email).Result()
//		if err != nil {
//			response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
//			return
//		}
//		if correctValiCode != valiCode {
//			response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
//			return
//		}
//		var key model.CDKey
//		db.Where("content =? AND used = 0", cdkey).First(&key)
//		if key.CDKeyID == 0 {
//			response.Response(c, http.StatusBadRequest, 400, nil, "激活码错误")
//			return
//		}
//		// 创建用户
//		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
//		if err != nil {
//			response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
//			return
//		}
//		//numInt, err := strconv.Atoi(num)
//		//fmt.Println(num, numInt)
//		//if err != nil {
//		//	response.Response(c, http.StatusInternalServerError, 400, nil, "学号异常")
//		//	return
//		//}
//		//if isNumExist(db, numInt) {
//		//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该学号已存在")
//		//	return
//		//}
//		// 创建新用户结构体
//		newUser := model.User{
//			Name:     name,
//			Phone:    telephone,
//			Password: string(hasedPassword),
//			Banend:   time.Now(),
//			Email:    email,
//			//Num:      numInt,
//		}
//		// 将结构体传进Create函数即可在数据库中添加一条记录
//		// 其他的增删查改功能参见postController里的updateLike函数
//		db.Create(&newUser)
//		if newUser.UserID == 0 {
//			response.Response(c, http.StatusInternalServerError, 400, nil, "userID为0")
//			return
//		}
//		// 修改激活码
//		// 警告：当使用 struct 更新时，GORM只会更新那些非零值的字段
//		db.Model(&key).Update(model.CDKey{
//			Used:     true,
//			UsedTime: time.Now(),
//		})
//		//发放token
//		token, err := common.ReleaseToken(newUser)
//		if err != nil {
//			response.Response(c, http.StatusInternalServerError, 400, nil, "账号已注册，但无法返回token")
//			log.Printf("token generate error: %v", err)
//			return
//		}
//		//返回结果
//		response.Success(c, gin.H{"token": token}, "注册成功")
//		rds.Del(ctx, email)
//	}
func Register(c *gin.Context) {
	// 连接数据库
	db := common.GetDB()
	// 从前端获取信息
	var requestUser = registerUser{}
	c.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Phone
	password1 := requestUser.Password
	password2 := requestUser.Password2
	email := strings.TrimSpace(requestUser.Email)
	valiCode := requestUser.ValiCode
	cdkey := requestUser.CDKey

	// 验证数据
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "手机号必须为11位!!!")
		return
	}
	if len(password1) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	if len(password2) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	if password1 != password2 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "两次密码不相同，请重新输入")
		return
	}
	//如果名称没有传，给一个10位的随机字符串
	if len(name) == 0 {
		name = gofakeit.Name()
	}
	// 判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该手机号已存在")
		return
	}
	check := false
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			check = true
		}
	}
	if !check {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	}
	if isEmailExist(db, email) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱已存在")
		return
	}
	rds := common.MyRedis
	ctx := context.Background()
	correctValiCode, err := rds.Get(ctx, email).Result()
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}
	if correctValiCode != valiCode {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}
	var key model.CDKey
	db.Where("content =? AND used = 0", cdkey).First(&key)
	if key.CDKeyID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "激活码错误")
		return
	}
	// 创建用户
	newUser := model.User{
		Name:     name,
		Phone:    telephone,
		Password: password1, // 直接存储明文密码
		Banend:   time.Now(),
		Email:    email,
	}
	db.Create(&newUser)
	if newUser.UserID == 0 {
		response.Response(c, http.StatusInternalServerError, 400, nil, "userID为0")
		return
	}
	// 修改激活码
	db.Model(&key).Update(model.CDKey{
		Used:     true,
		UsedTime: time.Now(),
	})
	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "账号已注册，但无法返回token")
		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "注册成功")
	rds.Del(ctx, email)
}
func IdentityValidate(c *gin.Context) {
	//db := common.GetDB()
	// 从前端获取信息
	var requestUser = identity{}
	c.Bind(&requestUser)
	//获取参数

	email := requestUser.Email
	valiCode := requestUser.ValiCode
	resCode := IdentityValidateService(email, valiCode)
	if resCode == VALIDATEFAILED {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
	} else if resCode == EMAILERROR {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
	} else if resCode == EMAILNOTEXIST {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
	} else {
		response.Response(c, http.StatusOK, 200, nil, "身份验证成功")
	}
}

// 验证码校验服务
func IdentityValidateService(email string, valiCode string) int {
	db := common.GetDB()
	check := false
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			check = true
		}
	}

	if !check {
		return EMAILERROR
	}
	if !isEmailExist(db, email) {
		return EMAILNOTEXIST
	}
	rds := common.MyRedis
	ctx := context.Background()
	correctValiCode, err := rds.Get(ctx, email).Result()
	if err != nil {
		return VALIDATEFAILED
	}
	if correctValiCode != valiCode {
		return VALIDATEFAILED
	}
	//返回结果
	rds.Del(ctx, email)
	return VALIDATESUCCESS
}

type loginuser struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	//Phone    string `gorm:"type:varchar(11);not null"`
	Email    string `json:"email"`
	Password string `gorm:"size:255;not null"`
}

// func Login(c *gin.Context) {
// 	// 邮箱登录
// 	db := common.GetDB()
// 	var requestUser = loginuser{}
// 	c.Bind(&requestUser)
// 	//获取参数
// 	//telephone := requestUser.Phone
// 	email := requestUser.Email
// 	password := requestUser.Password
// 	password = util.Decrypt(password)
// 	//数据验证
// 	//if len(telephone) == 0 {
// 	//	msg := "手机号为空！" + telephone
// 	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, msg)
// 	//	return
// 	//}
// 	//if len(telephone) != 11 {
// 	//	msg := "手机号必须为11位" + telephone
// 	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, msg)
// 	//	return
// 	//}
// 	if len(email) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱号为空！"+email)
// 		return
// 	}
// 	if len(password) < 6 {
// 		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
// 		return
// 	}
// 	//判断手机号是否存在
// 	var user model.User
// 	//db.Where("phone = ?", telephone).First(&user)
// 	db.Where("email = ?", email).First(&user)
// 	if user.UserID == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
// 		return
// 	}
// 	//判断密码是否正确
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
// 		return
// 	}
// 	//if user.IDpass == false {
// 	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户注册尚未通过审核，请耐心等待")
// 	//	return
// 	//}
// 	//发放token
// 	token, err := common.ReleaseToken(user)
// 	if err != nil {
// 		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")

//			log.Printf("token generate error: %v", err)
//			return
//		}
//		//返回结果
//		response.Success(c, gin.H{"token": token}, "登录成功")
//	}
func Login(c *gin.Context) {
	// 邮箱登录
	db := common.GetDB()
	var requestUser = loginuser{}
	c.Bind(&requestUser)
	//获取参数
	email := requestUser.Email
	password := requestUser.Password

	//数据验证
	if len(email) == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱号为空！"+email)
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	//判断邮箱是否存在
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if user.Password != password { // 直接比较明文密码
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")
}

// func ModifyPassword(c *gin.Context) {
// 	db := common.GetDB()
// 	var user model.User
// 	var inputUser ModifyUser
// 	c.Bind(&inputUser)
// 	email := inputUser.Email
// 	password := inputUser.Password
// 	password = util.Decrypt(password)
// 	password2 := inputUser.Password2
// 	password2 = util.Decrypt(password2)
// 	valiCode := inputUser.ValiCode
// 	if password != password2 {
// 		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码输入不一致")
// 		return
// 	}
// 	if !isEmailExist(db, email) {
// 		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "未找到邮箱")
// 		return
// 	}

// 	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
// 		return
// 	}
// 	db.Where("email = ?", email).First(&user)
// 	//fmt.Println(phone)
// 	//fmt.Println(password)
// 	if user.UserID == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
// 		return
// 	}

//		//校验验证码是否正确
//		resCode := IdentityValidateService(email, valiCode)
//		if resCode == VALIDATEFAILED {
//			response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
//			return
//		} else if resCode == EMAILERROR {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
//			return
//		} else if resCode == EMAILNOTEXIST {
//			response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
//			return
//		} else {
//			user.Password = string(hasedPassword)
//			db.Save(&user)
//			response.Success(c, gin.H{"data": user}, "修改密码成功")
//			return
//		}
//	}
func ModifyPassword(c *gin.Context) {
	db := common.GetDB()
	var user model.User
	var inputUser ModifyUser
	c.Bind(&inputUser)
	email := inputUser.Email
	password := inputUser.Password
	password2 := inputUser.Password2
	valiCode := inputUser.ValiCode

	if password != password2 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码输入不一致")
		return
	}
	if !isEmailExist(db, email) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "未找到邮箱")
		return
	}
	db.Where("email = ?", email).First(&user)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}

	//校验验证码是否正确
	resCode := IdentityValidateService(email, valiCode)
	if resCode == VALIDATEFAILED {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	} else if resCode == EMAILERROR {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	} else if resCode == EMAILNOTEXIST {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
		return
	} else {
		user.Password = password // 直接存储明文密码
		db.Save(&user)
		response.Success(c, gin.H{"data": user}, "修改密码成功")
		return
	}
}
func ValidateEmail(c *gin.Context) {
	var request requestEmail
	c.Bind(&request)
	email := request.Email
	mode := request.Mode
	fmt.Println("email = ", email)
	fmt.Println("mode = ", mode)

	db := common.DB
	if mode == 1 && !isEmailExist(db, email) {
		//fmt.Println(1)
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email, "mode": mode}, "该邮箱未注册，请先完成注册操作")
		return
	}
	if mode == 0 && isEmailExist(db, email) {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email, "mode": mode}, "邮箱已注册")
		return
	}
	if email == "" {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email}, "邮箱获取错误")
		return
	}
	err := api.SendEmail(email)
	if err != nil {
		//fmt.Println(2)
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email}, "发送邮箱错误")
		return
	}
	response.Success(c, gin.H{"data": email}, "邮箱发送成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件上传失败")
		return
	}

	// db := common.GetDB()
	// var user model.User

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := "public/uploads/" + filename

	// 保存文件到本地
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件保存失败")
		return
	}
	// 压缩头像
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败"})
		return
	}
	defer src.Close()
	// 读取上传文件的内容，并解码为 image.Image 对象
	img, _, err := image.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件解码失败"})
		return
	}

	// 对图像进行压缩处理
	resizedImage := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	// 这里似乎只支持绝对路径，这里要根据服务器修改
	resizedPath := "public/resized/" + filename
	// 服务器本地备份
	err = imaging.Save(resizedImage, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	fileURL := ""
	fileURL, err = api.UploadImage("/src/images/resized/"+filename, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 返回成功
	response.Success(c, gin.H{"avatar_url": fileURL}, "上传成功")
}

func UpdateAvatar(c *gin.Context) {
	// 用于移动端
	phone := c.PostForm("phone")
	if len(phone) != 11 {
		response.Response(c, http.StatusBadRequest, 400, nil, "Invalid phone number")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件上传失败")
		return
	}

	// Add a check for the file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" { // Add more formats if needed
		response.Response(c, http.StatusBadRequest, 400, nil, "Invalid file format")
		return
	}

	db := common.GetDB()
	var user model.User

	if db.Where("phone = ?", phone).First(&user).RecordNotFound() {
		response.Response(c, http.StatusNotFound, 404, nil, "User not found")
		return
	}
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename)) // Use base to avoid path traversal
	filepath := "public/uploads/" + filename

	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件保存失败")
		return
	}
	// 压缩头像
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败"})
		return
	}
	defer src.Close()
	// 读取上传文件的内容，并解码为 image.Image 对象
	img, _, err := image.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件解码失败"})
		return
	}

	// 对图像进行压缩处理
	resizedImage := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	// 这里似乎只支持绝对路径，这里要根据服务器修改
	resizedPath := "public/resized/" + filename
	// 服务器本地备份
	err = imaging.Save(resizedImage, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	fileURL := ""
	fileURL, err = api.UploadImage("/src/images/resized/"+filename, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 更新 User 的 AvatarURL 字段
	user.AvatarURL = fileURL
	if err := db.Save(&user).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Database update failed")
		return
	}

	response.Success(c, gin.H{"phone": user.Phone, "avatar_url": user.AvatarURL}, "上传成功")
}

// GetAvatar 用于处理获取用户头像的请求
func GetAvatar(c *gin.Context) {
	// 连接数据库
	db := common.GetDB()

	// 从前端获取电话号码
	phone := c.PostForm("phone")
	// 在数据库中查找用户
	var user model.User
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		} else {
			response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		}
		return
	}
	// 返回用户的头像URL
	response.Success(c, gin.H{"phone": user.Phone, "AvatarURL": user.AvatarURL}, "获取成功")
}
func GetInfo(c *gin.Context) {
	db := common.GetDB()
	// 从前端获取电话号码
	var requestUser = model.User{}
	c.Bind(&requestUser)
	phone := requestUser.Phone
	userID := requestUser.UserID
	// 在数据库中查找用户
	var user model.User
	if err := db.Where("phone = ? or userID = ?", phone, userID).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		} else {
			response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		}
		return
	}
	// 返回用户的所有信息
	c.JSON(http.StatusOK, gin.H{
		"userID": user.UserID,
		"phone":  user.Phone,
		"email":  user.Email,
		"name":   user.Name,
		//"num":       user.Num,
		"intro":     user.Intro,
		"score":     user.Score,
		"ban":       user.Banend,
		"punishnum": user.Punishnum,
		"avatarURL": user.AvatarURL,
	})
}

type updateUser struct {
	UserID    int    `json:"userID"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Num       int    `json:"num"`
	Intro     string `json:"intro"`
	AvatarURL string `json:"avatarURL"`
}

func UpdateUserInfo(c *gin.Context) {
	db := common.GetDB()

	// 解析请求参数
	var userInfo updateUser
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数解析错误")
		return
	}

	// 根据用户ID查找用户
	var user model.User
	if err := db.Where("userID = ?", userInfo.UserID).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Response(c, http.StatusNotFound, 404, nil, "用户不存在")
		} else {
			response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		}
		return
	}
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}

	// 更新用户信息
	user.Name = userInfo.Name
	//user.Num = userInfo.Num
	user.Intro = userInfo.Intro
	user.AvatarURL = userInfo.AvatarURL

	if err := db.Save(&user).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		return
	}
	response.Response(c, http.StatusOK, 200, nil, "用户信息更新成功")
}

func CalculateAndSaveScores() {
	log.Println("CalculateAndSaveScores")
	db := common.GetDB()
	var users []model.User
	db.Find(&users)
	for _, user := range users { //对每个用户进行积分统计
		if user.UserID == 0 {
			continue
		}
		// 查询用户的post ID
		var postIDs []int
		db.Model(&model.Post{}).Where("userID = ?", user.UserID).Pluck("postID", &postIDs)
		// 统计post被点赞个数
		totalLikeNum := 0
		for _, postID := range postIDs {
			var likeNum int
			db.Model(&model.Post{}).Where("postID = ?", postID).Select("like_num").Row().Scan(&likeNum)
			totalLikeNum += likeNum
		}
		// 查询用户的 pcomment ID
		var pcommentIDs []int
		db.Model(&model.Pcomment{}).Where("userID = ?", user.UserID).Pluck("pcommentID", &pcommentIDs)
		// 统计pcomment被点赞个数
		for _, pcommentID := range pcommentIDs {
			var likeNum int
			db.Model(&model.Pcomment{}).Where("pcommentID = ?", pcommentID).Select("like_num").Row().Scan(&likeNum)
			totalLikeNum += likeNum
		}
		// 统计被收藏个数
		var psaveCount int
		db.Model(&model.Psave{}).Where("ptargetID IN (?)", postIDs).Count(&psaveCount)
		// 统计帖子被评论个数
		var commentedCount int
		db.Model(&model.Pcomment{}).Where("ptargetID IN (?)", postIDs).Not("userID = ?", user.UserID).Count(&commentedCount)
		// 统计评论被回复个数
		var repliedCount int
		db.Model(&model.Ccomment{}).Where("ctargetID IN (?)", pcommentIDs).Not("userID = ?", user.UserID).Count(&repliedCount)
		// 查询用户的 ccomment ID
		var ccommentIDs []int
		db.Model(&model.Ccomment{}).Where("userID = ?", user.UserID).Pluck("ccommentID", &ccommentIDs)
		// 统计ccomment被点赞个数
		for _, ccommentID := range ccommentIDs {
			var likeNum int
			db.Model(&model.Ccomment{}).Where("ccommentID = ?", ccommentID).Select("like_num").Row().Scan(&likeNum)
			totalLikeNum += likeNum
		}
		// 统计成功举报个数
		var sueCount int
		db.Model(&model.Sue{}).Where("status = ? AND finish = ? AND userID = ?", " ok", 1, user.UserID).
			Select("DISTINCT targettype, targetID").
			Count(&sueCount)
		// 计算总积分并保存
		totalScore := len(postIDs)*10 + len(pcommentIDs)*5 + len(ccommentIDs)*5 + totalLikeNum + psaveCount*3 +
			(commentedCount+repliedCount)*2 + sueCount*20 - (user.Punishnum)*20
		db.Model(&user).Update("Score", totalScore)
	}
	log.Println("Score updated")
}
