Page({
  data: {
    partition: '', // 当前分区
    posts: [], // 帖子列表
    totalNum: 0, // 帖子总数
    curPage: 0, // 当前页码
    limit: 10, // 每页加载数量
    isLoading: false, // 是否正在加载
    scrollTop: 0, // 滚动位置
  },

  onLoad() {
    this.updateNum(); // 初始化帖子数量
    this.startObserver(); // 启动滚动加载监听
  },

  onUnload() {
    this.endObserver(); // 停止滚动加载监听
  },

  onPageScroll(e) {
    this.setData({ scrollTop: e.scrollTop }); // 记录滚动位置
  },

  // 获取帖子数量
  updateNum() {
    getPostsNum({
      partition: this.data.partition,
      searchsort: this.data.searchsort,
      searchinfo: this.data.searchinfo,
      userTelephone: this.data.userInfo.phone,
    }).then((id) => {
      this.setData({ totalNum: id });
    });
  },

  // 加载更多帖子
  addPosts() {
    getPosts({
      limit: this.data.limit,
      offset: this.data.curPage,
      partition: this.data.partition,
      searchsort: this.data.searchsort,
      searchinfo: this.data.searchinfo,
      userTelephone: this.data.userInfo.phone,
    }).then((res) => {
      if (res) {
        this.setData({ posts: [...this.data.posts, ...res] });
      } else {
        this.setData({ totalNum: this.data.posts.length });
      }
    });
  },

  // 删除帖子
  deleteHandler(callback) {
    callback().then((res) => {
      if (res) {
        this.setData({ curPage: this.data.curPage - 1 });
        this.updateNum();
        this.updatePosts(res);
      }
    });
  },

  // 更新帖子列表
  updatePosts(id) {
    const index = this.data.posts.findIndex((post) => post.PostID === id);
    if (index !== -1) {
      const posts = this.data.posts;
      posts.splice(index, 1);
      this.setData({ posts });
    }
  },

  // 滚动加载逻辑
  getMore() {
    if (this.data.isLoading) {
      this.addPosts();
      this.setData({ curPage: this.data.posts.length });
    }
  },

  // 启动滚动加载监听
  startObserver() {
    const observer = wx.createIntersectionObserver(this);
    observer.relativeToViewport({ bottom: 0 }).observe('.bottomDiv', (res) => {
      if (res.intersectionRatio > 0) {
        this.getMore();
      }
    });
    this.observer = observer;
  },

  // 停止滚动加载监听
  endObserver() {
    if (this.observer) {
      this.observer.disconnect();
    }
  },
});