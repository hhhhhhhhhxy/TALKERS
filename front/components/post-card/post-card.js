Component({
  properties: {
    post: {
      type: Object,
      value: {},
    },
    // 不再需要 deleteHandler 属性
  },

  data: {
    postData: {}, // 帖子数据
    userInfo: {}, // 用户信息
  },

  lifetimes: {
    attached() {
      // 初始化数据
      this.setData({
        postData: this.properties.post,
        userInfo: getApp().globalData.userInfo,
      });
      // 截取内容
      this.setData({
        'postData.Content': this.data.postData.Content.slice(0, 15),
      });
    },
  },

  methods: {
    // 跳转到帖子详情页
    navigateToPostDetail() {
      wx.navigateTo({
        url: `/pages/postDetail/postDetail?postId=${this.data.postData.PostID}`,
      });
    },

    // 收藏帖子
    async handleSave() {
      try {
        await savePost(this.data.postData.PostID, this.data.userInfo.phone);
        this.setData({
          'postData.IsSaved': !this.data.postData.IsSaved,
        });
        wx.showToast({
          title: this.data.postData.IsSaved ? '收藏成功' : '取消成功',
          icon: 'none',
        });
      } catch (e) {
        wx.showToast({
          title: '失败了:-(',
          icon: 'none',
        });
      }
    },

    // 点赞帖子
    async like() {
      try {
        const res = await likePost(
          this.data.postData.PostID,
          this.data.userInfo.phone
        );
        return res ? true : false;
      } catch (e) {
        console.error(e);
        return false;
      }
    },

    // 删除帖子
    async deleteHandler() {
      try {
        // 调用删除接口
        await delPost(this.data.postData.PostID);
        // 提示用户删除成功
        wx.showToast({
          title: '删除成功',
          icon: 'none',
        });
        // 更新组件状态（例如隐藏或标记为已删除）
        this.setData({
          postData: { ...this.data.postData, isDeleted: true }, // 标记为已删除
        });
      } catch (e) {
        console.error(e);
        wx.showToast({
          title: '删除失败',
          icon: 'none',
        });
      }
    },
  },
});