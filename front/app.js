// app.js
const API_BASE_URL = 'https://localhost:8080/api'; // 全局 API 地址

App({
  onLaunch() {
    // 小程序初始化时执行
    console.log('小程序初始化');

    wx.getUserInfo({
      success: res => {
        this.globalData.userInfo = res.userInfo;
      }
    });
    this.checkLoginStatus();
  },

  globalData: {
    userInfo: null, // 全局用户信息
    apiBaseUrl: API_BASE_URL, // 全局 API 地址
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token');
    if (token) {
      console.log('用户已登录');
    } else {
      console.log('用户未登录');
      // 可以跳转到登录页面
      wx.navigateTo({
        url: '/pages/login/login',
      });
    }
  },
});