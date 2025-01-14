Page({
  data: {
    partitions: [
      { name: '吐槽', description: '言短情长，抒发心声。' },
      { name: '求助', description: '疑虑相询，忧难互助。' },
      { name: '生活', description: '生活分享，交流思辨。' },
      { name: '职业', description: '专业相关，信息发布。' },
      { name: '求职招募', description: '梦想职路，共创未来。' },
      { name: '其他', description: '畅所欲言，创意无限。' },
    ],
  },

  // 跳转到分区
  sendPartition(event) {
    const partitionName = event.currentTarget.dataset.name;
    wx.setStorageSync('selectedPartition', partitionName); // 存储选择的分区
    wx.switchTab({
      url: '/pages/index/index', // 跳转到首页
    });
  },
});