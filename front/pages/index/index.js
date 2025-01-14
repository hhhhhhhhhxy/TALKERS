Page({
  data: {
    currentPage: 'home', // 当前页面
    posts: [
      { title: '帖子标题1', author: '用户A', time: '2023-10-01' },
      { title: '帖子标题2', author: '用户B', time: '2023-10-02' },
      { title: '帖子标题3', author: '用户C', time: '2023-10-03' },
    ],
    searchInfo: '',
  },

  // 跳转到分区页面
  navigateToPartitions() {
    wx.navigateTo({ url: '/pages/partition/partition' });
  },

  // 跳转到发帖页面
  navigateToPost() {
    wx.navigateTo({ url: '/pages/post-editor/post-editor' });
  },

  // 处理搜索输入
  handleSearchInput(e) {
    this.setData({ searchInfo: e.detail.value });
  },

  // 搜索功能
  search() {
    const { searchInfo } = this.data;
    console.log('搜索内容:', searchInfo);
    // 这里可以调用搜索接口
  }
});