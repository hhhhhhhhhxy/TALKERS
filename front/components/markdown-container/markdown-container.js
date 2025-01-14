Component({
  properties: {
    markdownContent: {
      type: String,
      value: '',
    },
  },

  data: {
    safeHTML: '', // 渲染后的 HTML 内容
  },

  observers: {
    // 监听 markdownContent 变化
    markdownContent(content) {
      if (content) {
        this.renderMarkdown(content);
      }
    },
  },

  methods: {
    // 渲染 Markdown 内容
    renderMarkdown(content) {
      // 使用微信小程序的 rich-text 渲染 Markdown
      const html = this.parseMarkdown(content);
      this.setData({ safeHTML: html });
    },

    // 解析 Markdown 内容
    parseMarkdown(content) {
      // 简单的 Markdown 解析逻辑
      const lines = content.split('\n');
      let html = '';

      lines.forEach((line) => {
        if (line.startsWith('# ')) {
          html += `<h1>${line.substring(2)}</h1>`;
        } else if (line.startsWith('## ')) {
          html += `<h2>${line.substring(3)}</h2>`;
        } else if (line.startsWith('### ')) {
          html += `<h3>${line.substring(4)}</h3>`;
        } else if (line.startsWith('```')) {
          html += `<pre><code>${line.substring(3)}</code></pre>`;
        } else {
          html += `<p>${line}</p>`;
        }
      });

      return html;
    },
  },
});