Component({
  properties: {
    src: {
      type: String,
      value: '',
    },
    userId: {
      type: Number,
      value: 0,
    },
    userIdentity: {
      type: String,
      value: '',
    },
  },

  data: {
    defaultAvatar: '/assets/default-avatar.svg', // 默认头像路径
  },

  methods: {
    // 跳转到用户主页
    navigate() {
      if (this.data.userId > 0) {
        wx.navigateTo({
          url: `/pages/user-profile/user-profile?id=${this.data.userId}`,
        });
      }
    },
  },
});