Page({
  data: {
    title: '登录', // 页面标题
    page: 0, // 当前页面状态：0-登录，1-注册，2-忘记密码
    email: '', // 邮箱
    username: '', // 用户名
    password1: '', // 密码
    password2: '', // 重复密码
    code: '', // 验证码
    CDkey: '', // 邀请码
    remembered: false, // 是否记住我
  },
  onLoad() {
    const token = wx.getStorageSync('token'); // 从本地存储获取 Token
    if (!token) {
        wx.showToast({ title: '请先登录', icon: 'none' });
        return;
    }

    wx.request({
        url: 'http://localhost:8080/api/auth/info',
        method: 'GET',
        header: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        success: (res) => {
            if (res.statusCode === 200) {
                console.log('用户信息:', res.data);
            } else {
                console.error('获取用户信息失败:', res.data);
            }
        },
        fail: (err) => {
            console.error('请求失败:', err);
        },
    });
},
  // 处理输入框内容
  handleEmailInput(event) {
    this.setData({ email: event.detail.value });
  },
  handleUsernameInput(event) {
    this.setData({ username: event.detail.value });
  },
  handlePassword1Input(event) {
    this.setData({ password1: event.detail.value });
  },
  handlePassword2Input(event) {
    this.setData({ password2: event.detail.value });
  },
  handleCodeInput(event) {
    this.setData({ code: event.detail.value });
  },
  handleCDkeyInput(event) {
    this.setData({ CDkey: event.detail.value });
  },
  handleRememberedChange(event) {
    this.setData({ remembered: event.detail.value });
  },

  // 切换页面状态
  jumpToReg() {
    const page = this.data.page === 0 ? 1 : 0;
    const title = page === 0 ? '登录' : '注册';
    this.setData({ page, title });
  },

  // 切换到忘记密码页面
  forget() {
    const page = this.data.page === 0 ? 2 : 0;
    const title = page === 0 ? '登录' : '忘记密码';
    this.setData({ page, title });
  },

  // 处理输入框内容
  handleEmailInput(event) {
    this.setData({ email: event.detail.value });
  },
  handleUsernameInput(event) {
    this.setData({ username: event.detail.value });
  },
  handlePassword1Input(event) {
    this.setData({ password1: event.detail.value });
  },
  handlePassword2Input(event) {
    this.setData({ password2: event.detail.value });
  },
  handleCodeInput(event) {
    this.setData({ code: event.detail.value });
  },
  handleCDkeyInput(event) {
    this.setData({ CDkey: event.detail.value });
  },
  handleRememberedChange(event) {
    this.setData({ remembered: event.detail.value });
  },

  // 登录
  login() {
    const { email, password1 } = this.data;
    if (email && password1) {
      wx.request({
        url: 'http://localhost:8080/api/auth/login',
        method: 'POST',
        data: {
          email,
          password: password1,
        },
        success: (res) => {
          if (res.statusCode === 200) {
            console.log('登录成功', res.data);
            wx.setStorageSync('token', res.data.token); // 保存 Token
            wx.switchTab({
              url: '/pages/index/index', // 登录成功后跳转到主页
            });
          } else {
            wx.showToast({ title: '登录失败', icon: 'none' });
          }
        },
        fail: (err) => {
          console.error('登录请求失败', err);
          wx.showToast({ title: '登录失败', icon: 'none' });
        },
      });
    } else {
      wx.showToast({ title: '请输入邮箱和密码', icon: 'none' });
    }
  },

  // 发送验证码
  getVCode(event) {
    const mode = event.currentTarget.dataset.mode; // 0-注册，1-重置密码
    const { email } = this.data;
    if (email) {
      wx.request({
        url: 'http://localhost:8080/api/auth/validateEmail',
        method: 'POST',
        data: {
          email,
          mode, // 0-注册，1-重置密码
        },
        success: (res) => {
          if (res.statusCode === 200) {
            console.log('验证码已发送', res.data);
            wx.showToast({ title: '验证码已发送', icon: 'none' });
          } else {
            wx.showToast({ title: '验证码发送失败', icon: 'none' });
          }
        },
        fail: (err) => {
          console.error('验证码请求失败', err);
          wx.showToast({ title: '验证码发送失败', icon: 'none' });
        },
      });
    } else {
      wx.showToast({ title: '请填写邮箱', icon: 'none' });
    }
  },

  // 注册
  reg() {
    const { CDkey, email, username, password1, password2, code } = this.data;
    if (CDkey && email && username && password1 && password2 && code) {
      if (password1 !== password2) {
        wx.showToast({ title: '两次密码不一致', icon: 'none' });
        return;
      }
      wx.request({
        url: 'http://localhost:8080/api/auth/register',
        method: 'POST',
        data: {
          CDkey,
          email,
          username,
          password: password1,
          code,
        },
        success: (res) => {
          if (res.statusCode === 200) {
            console.log('注册成功', res.data);
            this.setData({ page: 0, title: '登录' });
            wx.showToast({ title: '注册成功', icon: 'none' });
          } else {
            wx.showToast({ title: '注册失败', icon: 'none' });
          }
        },
        fail: (err) => {
          console.error('注册请求失败', err);
          wx.showToast({ title: '注册失败', icon: 'none' });
        },
      });
    } else {
      wx.showToast({ title: '请填写完整信息', icon: 'none' });
    }
  },

  // 重置密码
  resetPwd() {
    const { email, password1, password2, code } = this.data;
    if (email && password1 && password2 && code) {
      if (password1 !== password2) {
        wx.showToast({ title: '两次密码不一致', icon: 'none' });
        return;
      }
      wx.request({
        url: 'http://localhost:8080/api/auth/modifyPassword',
        method: 'POST',
        data: {
          email,
          password: password1,
          code,
        },
        success: (res) => {
          if (res.statusCode === 200) {
            console.log('密码重置成功', res.data);
            this.setData({ page: 0, title: '登录' });
            wx.showToast({ title: '密码重置成功', icon: 'none' });
          } else {
            wx.showToast({ title: '密码重置失败', icon: 'none' });
          }
        },
        fail: (err) => {
          console.error('密码重置请求失败', err);
          wx.showToast({ title: '密码重置失败', icon: 'none' });
        },
      });
    } else {
      wx.showToast({ title: '请填写完整信息', icon: 'none' });
    }
  },
});