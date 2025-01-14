Page({
  data: {
    modelValue: '', // 输入框内容
    isPreview: false, // 是否预览
    textareaHeight: 450, // 输入框高度
    routePath: '', // 当前路由路径
  },

  onLoad() {
    // 获取当前路由路径
    const pages = getCurrentPages();
    this.setData({ routePath: pages[pages.length - 1].route });

    // 读取草稿
    if (this.data.routePath === '/post') {
      const draft = wx.getStorageSync('draft');
      if (draft) {
        this.setData({ modelValue: draft });
        wx.showToast({ title: '读取草稿成功，已删除', icon: 'none' });
        wx.removeStorageSync('draft');
      }
    }
  },

  // 输入框内容变化
  handleInput(e) {
    const value = e.detail.value;
    this.setData({ modelValue: value });
    this.autoResize();
  },

  // 自动调整输入框高度
  autoResize() {
    const query = wx.createSelectorQuery();
    query.select('textarea').boundingClientRect((rect) => {
      this.setData({ textareaHeight: rect.height });
    }).exec();
  },

  // 上传图片
  uploadFile(e) {
    const file = e.detail.files[0];
    this.upload(file);
  },

  // 处理粘贴事件
  handlePaste(e) {
    const items = e.clipboardData.items;
    for (let i = 0; i < items.length; i++) {
      const item = items[i];
      if (item.kind === 'file') {
        const file = item.getAsFile();
        e.preventDefault();
        this.upload(file);
      }
    }
  },

  // 上传文件
  upload(file) {
    if (file) {
      wx.uploadFile({
        url: 'https://example.com/upload', // 替换为实际的上传接口
        filePath: file.path,
        name: 'file',
        success: (res) => {
          const data = JSON.parse(res.data);
          wx.showToast({ title: data.message, icon: 'none' });
          this.setData({
            modelValue: this.data.modelValue + `<img src="${data.fileURL}" alt="${file.name}" />`,
          });
          this.autoResize();
        },
        fail: (err) => {
          console.error('上传失败', err);
          wx.showToast({ title: '上传失败，请重试', icon: 'none' });
        },
      });
    } else {
      wx.showToast({ title: '请选择一个文件', icon: 'none' });
    }
  },

  // 编辑内容
  editContent(e) {
    const type = e.currentTarget.dataset.type;
    let insertion = '';
    let cursorOffset = 0;

    switch (type) {
      case '标题':
        insertion = '### 标题\n';
        cursorOffset = -1;
        break;
      case '粗体':
        insertion = '**粗体**';
        cursorOffset = -2;
        break;
      case '斜体':
        insertion = '*斜体*';
        cursorOffset = -1;
        break;
      case '删除线':
        insertion = '~~删除线~~';
        cursorOffset = -2;
        break;
      case '引用':
        insertion = '> 引用\n';
        cursorOffset = -1;
        break;
      case '无序列表':
        insertion = '- 无序列表';
        cursorOffset = 0;
        break;
      case '有序列表':
        insertion = '1. 有序列表';
        cursorOffset = 0;
        break;
      case '表格':
        insertion = '\n| 列1 | 列2 | 列3 |\n| --- | --- | --- |\n| 内容1 | 内容2 | 内容3 |\n';
        cursorOffset = 0;
        break;
      case '分割线':
        insertion = '\n---\n';
        cursorOffset = 0;
        break;
      case '代码块':
        insertion = '```\n代码块\n```';
        cursorOffset = -4;
        break;
      default:
        return;
    }

    const start = this.data.modelValue.length;
    const newValue = this.data.modelValue + insertion;
    this.setData({ modelValue: newValue });

    // 设置光标位置
    wx.nextTick(() => {
      const textarea = wx.createSelectorQuery().select('textarea');
      textarea.focus();
      textarea.setSelectionRange(start + insertion.length + cursorOffset, start + insertion.length + cursorOffset);
    });
  },

  // 切换预览
  togglePreview() {
    this.setData({ isPreview: !this.data.isPreview });
  },

  // 保存草稿
  savePost() {
    wx.setStorageSync('draft', this.data.modelValue);
    wx.showToast({ title: '已暂存为草稿', icon: 'none' });
  },

  // 发送
  send() {
    this.triggerEvent('send', { content: this.data.modelValue });
  },
});