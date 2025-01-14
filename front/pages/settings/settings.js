Page({
  data: {
    allInfo: {
      score: 0 // 默认值
    }, // 用户信息
    vCode: '', // 验证码
    password1: '', // 新密码
    password2: '', // 确认密码
  },

  onLoad() {
    this.getUserInfo();
    this.calculateProgress(); // 计算进度条宽度
  },
  
  // 计算进度条宽度
  calculateProgress() {
    const score = this.data.allInfo.score || 0; 
    const levelExp = this.levelExpHandler(score);
    const width = levelExp > 0 ? (score / levelExp) * 100 : 0; // 计算宽度，避免除零错误
    const progressWidth = Math.min(width, 100);
    this.setData({
      progressStyle: `width: ${progressWidth}%;`
    });
  },

  // 计算当前等级所需经验
  levelExpHandler(score) {
    if (score < 100) return 100;
    if (score < 500) return 500;
    return 1000;
  },

  // 获取用户信息
  async getUserInfo() {
    const res = await wx.cloud.callFunction({
      name: 'getAllInfo',
      data: { phone: this.data.userInfo.phone },
    });
    this.setData({ allInfo: res.result });
  },

  // 修改用户名
  updateName(e) {
    this.setData({ 'allInfo.name': e.detail.value });
  },

  // 修改简介
  updateIntro(e) {
    this.setData({ 'allInfo.intro': e.detail.value });
  },

  // 修改验证码
  updateVCode(e) {
    this.setData({ vCode: e.detail.value });
  },

  // 修改密码
  updatePassword1(e) {
    this.setData({ password1: e.detail.value });
  },

  // 确认密码
  updatePassword2(e) {
    this.setData({ password2: e.detail.value });
  },

  // 发送验证码
  async codeHandler() {
    try {
      await wx.cloud.callFunction({
        name: 'sendCode',
        data: { email: this.data.allInfo.email, type: 1 },
      });
      wx.showToast({ title: '验证码已发送', icon: 'none' });
    } catch (e) {
      wx.showToast({ title: '发送失败', icon: 'none' });
    }
  },

  // 重置密码
  async updatePasswordFunc() {
    const { password1, password2, vCode } = this.data;
    if (password1 === '' || password2 === '') {
      wx.showToast({ title: '密码不能为空', icon: 'none' });
      return;
    }
    if (password1 !== password2) {
      wx.showToast({ title: '两次密码不一致', icon: 'none' });
      return;
    }
    try {
      const res = await wx.cloud.callFunction({
        name: 'updatePassword',
        data: {
          email: this.data.allInfo.email,
          password1,
          password2,
          vCode,
        },
      });
      wx.showToast({ title: res.result.msg, icon: 'none' });
    } catch (e) {
      wx.showToast({ title: '重置失败', icon: 'none' });
    }
  },

  // 修改头像
  chooseAvatar() {
    wx.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        this.uploadAvatarFunc(res.tempFilePaths[0]);
      },
    });
  },

  // 上传头像
  async uploadAvatarFunc(filePath) {
    try {
      const res = await wx.cloud.uploadFile({
        cloudPath: `avatars/${Date.now()}.png`,
        filePath,
      });
      this.setData({ 'allInfo.avatarURL': res.fileID });
      wx.showToast({ title: '头像上传成功', icon: 'none' });
    } catch (e) {
      wx.showToast({ title: '上传失败', icon: 'none' });
    }
  },

  // 修改用户信息
  async updateUserInfoFunc() {
    const { avatarURL, intro, name } = this.data.allInfo;
    if (/^\s+|\s+$/.test(name)) {
      wx.showToast({ title: '用户名不能以空格开头或结尾', icon: 'none' });
      return;
    }
    try {
      const res = await wx.cloud.callFunction({
        name: 'updateUserInfo',
        data: { avatarURL, intro, name },
      });
      wx.showToast({ title: res.result.msg, icon: 'none' });
    } catch (e) {
      wx.showToast({ title: '修改失败', icon: 'none' });
    }
  },

  // 退出登录
  logout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          wx.clearStorageSync();
          wx.reLaunch({ url: '/pages/index/index' });
        }
      },
    });
  },

  // 跳转到收藏
  navigateToFavorites() {
    wx.navigateTo({ url: '/pages/favorites/favorites' });
  },

  // 跳转到历史
  navigateToHistory() {
    wx.navigateTo({ url: '/pages/history/history' });
  },

  // 跳转到反馈
  navigateToFeedback() {
    wx.navigateTo({ url: '/pages/feedback/feedback' });
  },
});