Page({
  data: {
    title: '', // 标题
    partition: '主页', // 当前分区
    partitions: ['主页', '吐槽', '求助', '生活', '职业', '求职招募', '其他'], // 分区列表
    postContent: '', // 正文内容
  },

  // 更新标题
  updateTitle(e) {
    this.setData({ title: e.detail.value });
  },

  // 更新分区
  updatePartition(e) {
    const index = e.detail.value;
    this.setData({ partition: this.data.partitions[index] });
  },

  // 提交帖子
  async submitPost() {
    const { title, partition, postContent } = this.data;
    const userInfo = getApp().globalData.userInfo;

    if (!title || !postContent) {
      wx.showToast({ title: '请填写完整信息', icon: 'none' });
      return;
    }

    try {
      const res = await wx.request({
        url: 'https://localhost:8080/auth/Post', // 替换为实际的发送帖子接口
        method: 'POST',
        data: {
          content: postContent,
          partition,
          title,
          phone: userInfo.phone,
        },
      });
      wx.showToast({ title: res.data.msg, icon: 'none' });
      wx.navigateBack(); // 返回上一页
    } catch (err) {
      console.error('提交失败', err);
      wx.showToast({ title: '提交失败，请重试', icon: 'none' });
    }
  },
});