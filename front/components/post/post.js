Page({
  data: {
    partition: '', // 当前分区
    searchinfo: '', // 搜索信息
    searchsort: '', // 搜索排序
    userInfo: {}, // 用户信息
    posts: [], // 帖子列表
    totalNum: 0, // 总帖子数
    curPage: 0, // 当前页码
    limit: 10, // 每页限制
    isLoading: true, // 是否正在加载
  },

  onLoad(options) {
    // 初始化数据
    this.setData({
      partition: getApp().globalData.partition,
      searchinfo: getApp().globalData.searchinfo,
      searchsort: getApp().globalData.searchsort,
      userInfo: getApp().globalData.userInfo,
    });

    this.updateNum();
    if (this.data.partition === '课程专区') {
      this.fetchTeachers();
    }
    this.startObserver();
  },

  onUnload() {
    this.endObserver();
  },

  onShow() {
    // 恢复滚动位置
    wx.pageScrollTo({
      scrollTop: this.data.scrollTop,
    });
  },

  async updateNum() {
    const res = await this.getPostsNum({
      partition: this.data.partition,
      searchsort: this.data.searchsort,
      searchinfo: this.data.searchinfo,
      userTelephone: this.data.userInfo.phone,
      tag: this.data.tag,
    });
    this.setData({ totalNum: res });
  },

  async addPosts() {
    const res = await this.getPosts({
      limit: this.data.limit,
      offset: this.data.curPage,
      partition: this.data.partition,
      searchsort: this.data.searchsort,
      searchinfo: this.data.searchinfo,
      userTelephone: this.data.userInfo.phone,
      tag: this.data.tag,
    });
    if (res) {
      this.setData({ posts: [...this.data.posts, ...res] });
    } else {
      this.setData({ totalNum: this.data.posts.length });
    }
  },

  async updatePosts() {
    this.setData({ posts: [], curPage: 0 });
    await this.addPosts();
  },

  deleteHandler(e) {
    const postId = e.detail;
    const posts = this.data.posts.filter(post => post.PostID !== postId);
    this.setData({ posts });
  },

  startObserver() {
    const observer = wx.createIntersectionObserver(this);
    observer.relativeToViewport({ bottom: 0 }).observe('.bottomDiv', (res) => {
      if (res.intersectionRatio > 0 && this.data.isLoading) {
        this.getMore();
      }
    });
  },

  endObserver() {
    if (this.observer) {
      this.observer.disconnect();
    }
  },

  async getMore() {
    if (this.data.isLoading) {
      await this.addPosts();
      this.setData({ curPage: this.data.posts.length });
    }
  },

  /**
   * 获取帖子列表
   * @param {Object} object - 请求参数
   * @returns {Promise<Array>} - 返回帖子列表
   */
  async getPosts(object) {
    try {
      const res = await wx.request({
        url: 'https://localhost:8080/auth/auth/browse', // 替换为你的 API 地址
        method: 'POST',
        header: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${wx.getStorageSync('token')}`, // 如果需要 token
        },
        data: object,
      });
      return res.data;
    } catch (e) {
      console.error(e);
      return null;
    }
  },

  /**
   * 获取帖子数量
   * @param {Object} object - 请求参数
   * @returns {Promise<number>} - 返回帖子数量
   */
  async getPostsNum(object) {
    try {
      const res = await wx.request({
        url: 'https://localhost:8080/auth/getPostNum', // 替换为你的 API 地址
        method: 'POST',
        header: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${wx.getStorageSync('token')}`, // 如果需要 token
        },
        data: object,
      });
      return res.data.Postcount;
    } catch (e) {
      console.error(e);
      return null;
    }
  },
});