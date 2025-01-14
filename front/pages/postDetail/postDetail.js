Page({
  data: {
    post: {}, // 帖子详情
    isPostLoaded: false, // 帖子是否加载完成
    comments: [], // 评论列表
    postCommentID: -1, // 当前展开的评论ID
    sortType: 'time', // 排序类型：time（时间） / likes（热度）
  },

  onLoad(options) {
    const postID = options.id; // 从路由参数中获取帖子ID
    this.getPostDetail(postID); // 获取帖子详情
    this.getCommentList(postID); // 获取评论列表
  },

  // 获取帖子详情
  async getPostDetail(postID) {
    try {
      const res = await wx.request({
        url: 'https://localhost:8080/auth/getPostDetail',
        method: 'POST',
        data: { postID },
      });
      if (res.data && !res.data.code) {
        this.setData({
          post: res.data,
          isPostLoaded: true,
        });
      } else {
        console.error('获取帖子详情失败:', res.data.msg);
      }
    } catch (error) {
      console.error('请求失败:', error);
    }
  },

  // 获取评论列表
  async getCommentList(postID) {
    try {
      const res = await wx.request({
        url: 'https://localhost:8080/auth/getCommentsByPostID',
        method: 'POST',
        data: { postID },
      });
      if (res.data && !res.data.code) {
        this.setData({
          comments: res.data.reverse(), // 反转评论列表，最新评论在前
        });
      } else {
        console.error('获取评论列表失败:', res.data.msg);
      }
    } catch (error) {
      console.error('请求失败:', error);
    }
  },

  // 设置排序类型
  setSortType(event) {
    const type = event.currentTarget.dataset.type;
    this.setData({ sortType: type });
  },

  // 切换子评论显示
  toggleComment(event) {
    const commentID = event.currentTarget.dataset.id;
    this.setData({
      postCommentID: this.data.postCommentID === commentID ? -1 : commentID,
    });
  },

  // 计算排序后的评论列表
  sortedComments() {
    const { comments, sortType } = this.data;
    if (sortType === 'time') {
      return comments; // 按时间排序（默认顺序）
    } else if (sortType === 'likes') {
      return comments.slice().sort((a, b) => b.LikeNum - a.LikeNum); // 按点赞数排序
    }
    return comments;
  },
});