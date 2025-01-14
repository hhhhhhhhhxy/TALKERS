Page({
  data: {
    heatPosts: [], // 热榜帖子列表
  },

  onLoad() {
    this.getHeatPosts();
  },

  // 获取热榜帖子
  async getHeatPosts() {
    try {
      const res = await wx.request({
        url: 'https://localhost:8080/api/auth/calculateHeat', // 替换为你的 API 地址
        method: 'GET',
      });
      if (res.data && res.data.code === 0) {
        this.setData({ heatPosts: res.data.data });
      } else {
        wx.showToast({
          title: res.data.msg || '获取热榜失败',
          icon: 'none',
        });
      }
    } catch (error) {
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none',
      });
      console.error('获取热榜失败:', error);
    }
  },
});