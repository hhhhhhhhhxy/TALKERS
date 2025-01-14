Component({
  properties: {
    cardData: {
      type: Object,
      value: {},
    },
    likeHandler: {
      type: null, // 将类型改为 null
      value: null,
    },
    isPost: {
      type: Boolean,
      value: false,
    },
  },

  data: {
    basicData: {
      Like: 0,
      IsLiked: false,
    },
  },

  lifetimes: {
    attached() {
      // 初始化数据
      this.setData({
        basicData: {
          Like: this.properties.cardData.Like,
          IsLiked: this.properties.cardData.IsLiked,
        },
      });
    },
  },

  methods: {
    // 点赞
    async like() {
      try {
        this.setData({
          'basicData.IsLiked': !this.data.basicData.IsLiked,
        });
        const res = await this.properties.likeHandler();
        if (res) {
          this.setData({
            'basicData.Like': this.data.basicData.IsLiked
              ? this.data.basicData.Like + 1
              : this.data.basicData.Like - 1,
          });
          wx.showToast({
            title: this.data.basicData.IsLiked ? '点赞成功' : '取消成功',
            icon: 'none',
          });
        } else {
          wx.showToast({
            title: '失败了:-(',
            icon: 'none',
          });
        }
      } catch (e) {
        wx.showToast({
          title: '失败了:-(',
          icon: 'none',
        });
      }
    },

    // 等级名称处理
    levelNameHandler(score) {
      // 实现等级名称逻辑
      return 'Lv1';
    },

    // 等级样式处理
    levelClassHandler(score) {
      // 实现等级样式逻辑
      return 'level-1';
    },

    // 字符串处理
    strHandler(type, value) {
      // 实现字符串处理逻辑
      return value;
    },
  },
});