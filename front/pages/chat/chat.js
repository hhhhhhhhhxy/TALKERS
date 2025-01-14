// pages/chat/chat.js
Page({
  data: {
    draft: '', // 输入框内容
    messages: [], // 聊天记录
    user: {
      userID: 0,
      name: '我',
    }, // 当前用户信息
    currentUser: {}, // 当前聊天对象
  },

  onLoad(options) {
    // 从路由参数中获取聊天对象信息
    const currentUser = {
      userID: options.userID,
      name: options.name,
      intro: options.intro,
    };
    this.setData({ currentUser });

    // 设置导航栏标题为对方的用户名
    wx.setNavigationBarTitle({
      title: currentUser.name,
    });

    // 模拟加载聊天记录
    this.loadChatHistory();
  },

  // 加载聊天记录
  loadChatHistory() {
    const messages = [
      {
        chatMsgID: 1,
        senderUserID: 1,
        content: '你好！',
        createdAt: Date.now() - 10000,
        read: true, // 对方已读
      },
      {
        chatMsgID: 2,
        senderUserID: 0,
        content: '你好，用户A！',
        createdAt: Date.now() - 5000,
        read: false, // 对方未读
      },
    ];
    this.setData({ messages });
  },

  // 处理输入框内容
  handleInput(event) {
    this.setData({ draft: event.detail.value });
  },

  // 发送消息
  sendMessage() {
    const { draft, user, currentUser } = this.data;
    if (draft.trim() === '') return;

    const message = {
      chatMsgID: Date.now(),
      senderUserID: user.userID,
      content: draft,
      createdAt: Date.now(),
      read: false, // 默认未读
    };

    this.setData({
      draft: '',
      messages: [...this.data.messages, message],
    });

    // 模拟对方查看消息
    setTimeout(() => {
      this.markMessageAsRead(message.chatMsgID);
    }, 2000); // 2秒后标记为已读
  },

  // 标记消息为已读
  markMessageAsRead(chatMsgID) {
    const messages = this.data.messages.map((item) => {
      if (item.chatMsgID === chatMsgID) {
        return { ...item, read: true };
      }
      return item;
    });
    this.setData({ messages });
  },
});