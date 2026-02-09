<script setup>
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-section">
          <h3>关于博客</h3>
          <p>分享技术、记录生活、传播知识</p>
        </div>

        <div class="footer-section">
          <h3>快速链接</h3>
          <router-link to="/">首页</router-link>
          <router-link to="/articles">文章列表</router-link>
          <router-link to="/category/1">分类</router-link>
          <router-link to="/about">关于</router-link>
        </div>

        <div class="footer-section">
          <h3>联系方式</h3>
          <p>Email: 2641070015@qq.com</p>
          <p>微信: yj2641070015</p>
        </div>
      </div>

      <div class="footer-bottom">
        <p>&copy; {{ currentYear }} 我的博客. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  background: linear-gradient(135deg, #38bdf8 0%, #0ea5e9 50%, #0284c7 100%);
  background-size: 200% 200%;
  animation: gradientShift 10s ease infinite;
  color: #fff;
  padding: 50px 0 30px;
  margin-top: 80px;
  position: relative;
  overflow: hidden;
}

@keyframes gradientShift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: #66a9bd;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 50px;
  margin-bottom: 50px;
  position: relative;
  z-index: 1;
}

.footer-section h3 {
  font-size: 20px;
  margin-bottom: 24px;
  background: linear-gradient(135deg, #667eea 0%, #f093fb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 700;
  letter-spacing: -0.3px;
}

.footer-section p {
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 12px;
  line-height: 1.7;
  font-size: 15px;
}

.footer-section a {
  display: block;
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 14px;
  transition: all 0.3s;
  font-size: 15px;
}

.footer-section a:hover {
  color: #fff;
  transform: translateX(8px);
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.15);
  padding-top: 25px;
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  position: relative;
  z-index: 1;
}

.footer-bottom p {
  margin: 0;
  font-size: 14px;
}
</style>
