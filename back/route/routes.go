package route

import (
	"loginTest/api"
	"loginTest/controller"
	"loginTest/middleware"

	"github.com/gin-gonic/gin"
)

// 创建路由

func CollectRoute(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CORSMiddleware())
	// 这里把路由分成了两组，其中auth是需要token验证的，也就是需要用户登录，noauth是不需要token的，也就是不需要用户登录。
	auth := r.Group("")
	noauth := r.Group("")

	auth.Use(middleware.AuthMiddleware())
	noauth.POST("/api/auth/register", controller.Register)
	noauth.POST("/api/auth/login", controller.Login)
	auth.POST("/api/auth/apiTest", api.ApiTest)
	auth.POST("/api/auth/post", controller.Post)
	auth.POST("/api/auth/browse", controller.Browse)
	auth.POST("/api/auth/getPostNum", controller.GetPostNum)
	auth.POST("/api/auth/deleteMe", controller.DeleteMe)
	auth.POST("/api/auth/updateLike", controller.UpdateLike)
	auth.POST("/api/auth/showDetails", controller.ShowDetails)
	auth.POST("api/auth/showPcomments", controller.GetComments)
	auth.POST("api/auth/postPcomment", controller.PostPcomment)
	auth.POST("api/auth/postCcomment", controller.PostCcomment)
	auth.POST("api/auth/updateCcommentLike", controller.UpdateCcommentLike)
	auth.POST("api/auth/updatePcommentLike", controller.UpdatePcommentLike)
	auth.POST("api/auth/updateCcommentDeny", controller.UpdateCcommentDeny)
	auth.POST("api/auth/updatePcommentDeny", controller.UpdatePcommentDeny)
	auth.POST("/api/auth/updateSave", controller.UpdateSave)
	auth.POST("/api/auth/updateBrowseNum", controller.UpdateBrowseNum)
	auth.POST("/api/auth/deletePost", controller.DeletePost)
	auth.POST("/api/auth/deletePcomment", controller.DeletePcomment)
	auth.POST("/api/auth/deleteCcomment", controller.DeleteCcomment)
	auth.POST("/api/auth/submitReport", controller.SubmitReport)
	auth.GET("/api/auth/info", controller.Info)
	noauth.POST("/api/auth/uploadPhotos", controller.UploadPhotos)
	noauth.POST("/api/auth/uploadAvatar", controller.UploadAvatar)
	noauth.POST("/api/auth/updateAvatar", controller.UpdateAvatar)
	auth.POST("/api/auth/getAvatar", controller.GetAvatar)
	auth.POST("/api/auth/getInfo", controller.GetInfo)
	auth.POST("/api/auth/updateUserInfo", controller.UpdateUserInfo)
	noauth.POST("/api/auth/uploadZip", controller.UploadZip)
	auth.POST("/api/auth/submitFeedback", controller.SubmitFeedback)
	auth.GET("/api/auth/getAllFeedback", controller.GetAllFeedback)
	auth.GET("/api/auth/calculateHeat", controller.CalculateHeat)
	auth.GET("/api/auth/getNotice", controller.GetNotice)
	auth.GET("/api/auth/getNoticeNum", controller.GetNoticeNum)
	auth.PATCH("api/auth/readNotice/:noticeID", controller.ReadNotice)
	noauth.POST("/api/auth/modifyPassword", controller.ModifyPassword)
	noauth.POST("/api/auth/validateEmail", controller.ValidateEmail)
	auth.POST("/api/auth/identityValidate", controller.IdentityValidate)
	auth.POST("/api/auth/getFileList", controller.GetFileList)
	auth.POST("/api/auth/fileDelete", controller.DeleteFile)
	auth.POST("/api/auth/folderDelete", controller.DeleteFolder)
	auth.POST("/api/auth/addFolder", controller.AddFolder)
	auth.POST("/api/auth/getObjectUrl", controller.GetObjectUrl)
	auth.GET("/api/auth/getTags", controller.GetTags)
	auth.POST("/api/auth/gettitle", controller.GetTitle)
	//获得所有聊天用户列表，并且建立ws链接
	auth.GET("/websocket/auth/chat", controller.ChatHandler)
	//点击与某人聊天主页的接口，需要监听用户是否在聊天页面
	//auth.GET("/api/auth/intoChatView", controller.IntoChatView)
	//退出与某人聊天主页的接口，需要监听用户是否在聊天页面
	//auth.POST("/api/auth/leaveChatView", controller.LeaveChatView)
	auth.GET("/api/auth/getChatHistory", controller.GetChatHistory)
	auth.GET("/api/auth/getChatNotice", controller.GetChatNotice)
	// 给管理员设置一个新的路由分组
	adminAuth := r.Group("")
	adminAuth.Use(middleware.AuthMiddleware_admin())
	//adminAuth.POST("/api/auth/passUsers", controller.PassUsers)
	adminAuth.POST("/api/auth/addAdmin", controller.AddAdmin)
	adminAuth.POST("/api/auth/changePassword", controller.ChangeAdminPassword)
	adminAuth.POST("/api/auth/deleteUser", controller.DeleteUser)
	adminAuth.POST("/api/auth/deleteAdmin", controller.DeleteAdmin)
	adminAuth.POST("/api/auth/showUsers", controller.ShowFilterUsers)
	r.POST("/api/auth/adminLogin", controller.AdminLogin)
	adminAuth.GET("/api/auth/admininfo", controller.AdminInfo)
	adminAuth.GET("/api/auth/getSues", controller.GetSues)
	adminAuth.POST("/api/auth/noViolation", controller.NoViolation)
	adminAuth.POST("/api/auth/violation", controller.Violation)
	adminAuth.POST("/api/auth/adminPost", controller.AdminPost)
	return r
}
