// pages/user-list/user-list.js
Page({
  data: {
    users: [
      {
        userID: 1,
        name: '用户A',
        intro: '这是用户A的简介',
        unreadCount: 2, // 未读消息数量
      },
      {
        userID: 2,
        name: '用户B',
        intro: '这是用户B的简介',
        unreadCount: 0, // 未读消息数量
      },
    ], // 用户列表
  },

  // 跳转到聊天页面
  navigateToChat(event) {
    const user = event.currentTarget.dataset.user;
    wx.navigateTo({
      url: `/pages/chat/chat?userID=${user.userID}&name=${user.name}&intro=${user.intro}`,
    });

    // 清空该用户的未读消息数量
    this.clearUnreadCount(user.userID);
  },

  // 清空未读消息数量
  clearUnreadCount(userID) {
    const users = this.data.users.map((user) => {
      if (user.userID === userID) {
        return { ...user, unreadCount: 0 };
      }
      return user;
    });
    this.setData({ users });
  },
});