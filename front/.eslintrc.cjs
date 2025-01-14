module.exports = {
  root: true,
  extends: [
    'eslint:recommended',
    'plugin:miniprogram/recommended'
  ],
  parserOptions: {
    ecmaVersion: 'latest',
    ecmaFeatures: {
      jsx: false // 微信小程序不支持 JSX
    }
  },
  plugins: ['miniprogram'],
  rules: {
    // 自定义规则
  }
};