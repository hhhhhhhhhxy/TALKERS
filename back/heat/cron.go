package heat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"loginTest/config"
	"loginTest/controller"
	"math"
	"time"

	"github.com/robfig/cron/v3"
	"loginTest/common"
	"loginTest/model"
)

// 启动热度计算定时任务
func StartHeat() {
	RefreshHeat()
	// 创建 cron 实例并添加任务
	c := cron.New()

	// 每h刷新
	c.AddFunc("@hourly", RefreshHeat)
	c.Start()

	// 防止程序退出
	select {}
}

// 刷新热度并存入 Redis ZSet
func RefreshHeat() {
	ctx := context.Background()
	db := common.GetDB()

	// 获取当前时间并计算10天前的时间
	now := time.Now()
	tenDaysAgo := now.AddDate(0, 0, -config.HEAT_DAY)

	// 获取10天内的所有帖子
	var posts []model.Post
	db.Where("post_time >= ?", tenDaysAgo).Find(&posts)

	// Redis ZSet的key，用来存储帖子热度
	redisKey := "hot_posts_zset"
	// 清空之前的 ZSet，确保只有当天最新的热度数据
	common.MyRedis.Del(ctx, redisKey)

	// 计算每个帖子的热度并存入 Redis ZSet
	for _, post := range posts {
		heatScore := calculateHeat(post.BrowseNum, post.CommentNum, post.LikeNum, math.Floor(now.Sub(post.PostTime).Hours()/24))

		// 创建包含 PostID 和 Title 的 JSON 字符串
		postInfo := controller.PostInfo{
			PostID:    uint(post.PostID),
			Title:     post.Title,
			BrowseNum: float64(post.BrowseNum),
		}
		postInfoJSON, err := json.Marshal(postInfo)
		if err != nil {
			log.Println("Error marshalling postInfo:", err)
			continue
		}

		// 将热度和帖子信息存入 Redis ZSet
		common.MyRedis.ZAdd(ctx, redisKey, redis.Z{
			Score:  heatScore,
			Member: postInfoJSON, // 存储 JSON 字符串
		})
	}

	// 获取Redis中热度前10的帖子
	topPosts, err := common.MyRedis.ZRevRangeWithScores(ctx, redisKey, 0, 9).Result()
	if err != nil {
		log.Println("获取Redis热度前10帖子失败:", err)
	} else {
		fmt.Println("Redis热度前10帖子已刷新:", topPosts)
	}
}

// 计算帖子热度值的函数
func calculateHeat(browseNum, commentNum, likeNum int, daysSincePost float64) float64 {
	return math.Log10(float64(browseNum)+float64(commentNum)/10+float64(likeNum)/2) - daysSincePost
}
