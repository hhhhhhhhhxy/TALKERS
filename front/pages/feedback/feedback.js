Page({
  data: {
    feedbackContent: '', // 反馈内容
  },

  // 监听输入内容
  onInput(event) {
    this.setData({
      feedbackContent: event.detail.value.trim(), // 去除空格
    });
  },

  // 提交反馈
  async submitFeedback() {
    const { feedbackContent } = this.data;
    if (feedbackContent) {
      try {
        const res = await wx.request({
          url: 'https://localhost:8080/auth/feedback', // 替换为你的 API 地址
          method: 'POST',
          data: {
            content: feedbackContent,
            type: '', // 如果有其他字段，可以在这里添加
          },
        });
        if (res.data && res.data.code === 0) {
          wx.showToast({
            title: res.data.msg || '提交成功',
            icon: 'success',
          });
          this.setData({ feedbackContent: '' }); // 清空输入框
        } else {
          wx.showToast({
            title: res.data.msg || '提交失败',
            icon: 'none',
          });
        }
      } catch (error) {
        wx.showToast({
          title: '网络错误，请重试',
          icon: 'none',
        });
        console.error('提交反馈失败:', error);
      }
    } else {
      wx.showToast({
        title: '请输入反馈内容',
        icon: 'none',
      });
    }
  },
});