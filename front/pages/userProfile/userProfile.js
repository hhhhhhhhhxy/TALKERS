Page({
  data: {
    user: null,
    userInfo: {},
  },

  onLoad(options) {
    const { id } = options;
    this.setData({ userInfo: getApp().globalData.userInfo });

    this.getInfoById(Number(id))
      .then((res) => {
        if (!res.code) {
          this.setData({
            user: {
              ...res,
              title: this.getUserTitle(res.score),
            },
          });
        } else {
          console.log(`请求无效：${res.msg} (${res.code})`);
        }
      })
      .catch((error) => {
        if (error.response) {
          const res = error.response;
          console.log(`请求无效：${res.msg} (${res.code})`);
        } else {
          console.log('请求无效');
        }
      });
  },

  navigateChat() {
    const { user } = this.data;
    if (user.userID > 0) {
      wx.navigateTo({
        url: `/pages/chat/chat?user=${user.userID}`,
      });
    }
  },

  getUserTitle(userScore) {
    if (userScore < 100) return '菜鸟';
    if (userScore >= 100 && userScore < 300) return '大虾';
    if (userScore >= 300 && userScore < 600) return '码农';
    if (userScore >= 600 && userScore < 1000) return '程序猿';
    if (userScore >= 1000 && userScore < 2000) return '工程师';
    if (userScore >= 2000 && userScore < 3000) return '大牛';
    if (userScore >= 3000 && userScore < 4000) return '专家';
    if (userScore >= 4000 && userScore < 5000) return '大神';
    return '祖师爷';
  },

  /**
   * 获取用户信息
   * @param {number} userID - 用户ID
   * @returns {Promise<Object>} - 返回用户信息
   */
  async getInfoById(userID) {
    try {
      const res = await wx.request({
        url: 'https://localhost:8080/auth/getInfo',
        method: 'POST',
        header: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${wx.getStorageSync('token')}`,
        },
        data: { userID },
      });
      return res.data;
    } catch (e) {
      console.error(e);
      return null;
    }
  },
});