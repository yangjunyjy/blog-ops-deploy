<template>
  <div class="about-page" :class="{ 'dark-mode': themeStore.isDark }">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="bg-circle circle-1"></div>
      <div class="bg-circle circle-2"></div>
      <div class="bg-circle circle-3"></div>
      <div class="star star-1">⭐</div>
      <div class="star star-2">⭐</div>
      <div class="star star-3">⭐</div>
      <div class="star star-4">⭐</div>
      <div class="star star-5">⭐</div>
    </div>

    <!-- 主容 -->
    <div class="main-wrapper">
      <!-- 个人信息卡片 -->
      <div class="profile-card">
        <!-- <div class="card-decoration"></div> -->
        <div class="avatar-section">
          <div class="avatar-glow"></div>
          <div class="avatar-border">
            <div class="avatar-inner">
              <span class="avatar-emoji">🎮</span>
            </div>
          </div>
          <div class="status-badge">在线</div>
        </div>
        <h1 class="profile-name">运维工程师</h1>
        <p class="profile-tagline">🌟 热爱代码 · 探索世界 · 记录成长</p>
        <div class="profile-stats">
          <div class="stat-item">
            <div class="stat-icon">📚</div>
            <div class="stat-value">100+</div>
            <div class="stat-label">文章</div>
          </div>
          <div class="stat-divider">|</div>
          <div class="stat-item">
            <div class="stat-icon">💻</div>
            <div class="stat-value">3+</div>
            <div class="stat-label">项目</div>
          </div>
          <div class="stat-divider">|</div>
          <div class="stat-item">
            <div class="stat-icon">⏰</div>
            <div class="stat-value">1+</div>
            <div class="stat-label">年经验</div>
          </div>
        </div>
        <div class="social-links">
          <a href="mailto:example@email.com" class="social-btn">
            <span class="social-icon">📧</span>
            <span class="social-text">邮箱</span>
          </a>
          <a href="#" class="social-btn">
            <span class="social-icon">🐙</span>
            <span class="social-text">GitHub</span>
          </a>
          <a href="#" class="social-btn">
            <span class="social-icon">🐦</span>
            <span class="social-text">Twitter</span>
          </a>
        </div>
      </div>

      <!-- 技能展示-->
      <div class="section-card skills-card">
        <div class="section-header">
          <h2 class="section-title">
            <span class="title-icon">⚔️</span>
            技能库
          </h2>
        </div>
        <div class="skills-container">
          <div class="skill-category">
            <h3 class="category-title">前端</h3>
            <div class="skill-list">
              <span class="skill-badge" style="--skill-color: #42b883">Vue.js</span>
              <span class="skill-badge" style="--skill-color: #f7df1e">JavaScript</span>
              <span class="skill-badge" style="--skill-color: #264de4">CSS3</span>
            </div>
          </div>
          <div class="skill-category">
            <h3 class="category-title">后端</h3>
            <div class="skill-list">
              <span class="skill-badge" style="--skill-color: #00add8">Go</span>
              <span class="skill-badge" style="--skill-color: #68a063">Java</span>
              <span class="skill-badge" style="--skill-color: #3776ab">Python</span>
            </div>
          </div>
          <div class="skill-category">
            <h3 class="category-title">工具</h3>
            <div class="skill-list">
              <span class="skill-badge" style="--skill-color: #f05032">nginx</span>
              <span class="skill-badge" style="--skill-color: #2496ed">Docker</span>
              <span class="skill-badge" style="--skill-color: #4479a1">k8s</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 时间�?-->
      <div class="section-card timeline-card">
        <div class="section-header">
          <h2 class="section-title">
            <span class="title-icon">🗺</span>
            成长轨迹
          </h2>
        </div>
        <div class="timeline">
          <div class="timeline-item" v-for="(item, index) in timelineData" :key="index">
            <div class="timeline-dot"></div>
            <div class="timeline-content">
              <div class="timeline-year">{{ item.year }}</div>
              <div class="timeline-title">{{ item.title }}</div>
              <div class="timeline-desc">{{ item.desc }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 兴趣爱好 -->
      <div class="section-card hobbies-card">
        <div class="section-header">
          <h2 class="section-title">
            <span class="title-icon">🎯</span>
            兴趣爱好
          </h2>
        </div>
        <div class="hobbies-grid">
          <div class="hobby-item">
            <div class="hobby-icon">🎮</div>
            <div class="hobby-name">游戏</div>
            <div class="hobby-desc">RPG & 独立游戏爱好 </div>
          </div>
          <div class="hobby-item">
            <div class="hobby-icon">📺</div>
            <div class="hobby-name">动漫</div>
            <div class="hobby-desc">追番日常</div>
          </div>
          <div class="hobby-item">
            <div class="hobby-icon">🎵</div>
            <div class="hobby-name">音乐</div>
            <div class="hobby-desc">电子音乐 & J-Pop</div>
          </div>
          <div class="hobby-item">
            <div class="hobby-icon">📖</div>
            <div class="hobby-name">阅读</div>
            <div class="hobby-desc">技术书籍 & 轻小说</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import HaloButton from '@/components/HaloButton.vue'

const themeStore = useThemeStore()

onMounted(() => {
  document.title = '关于 - 我的博客'
})

const timelineData = ref([
  { year: '2022', title: '启程', desc: '2020年开始学习编程' },
  { year: '2023', title: '突破', desc: '完成第一个全栈项目' },
  { year: '2024', title: '探索', desc: '深入研究更多技术领域' },
  { year: '2025', title: '成长', desc: '持续学习，保持好奇心' }
])
</script>

<style scoped>
.about-page {
  min-height: calc(100vh - 80px);
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
  padding: 40px 20px;
  position: relative;
  overflow: hidden;
  transition: background 0.3s ease;
}

.about-page.dark-mode {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

/* 背景装饰 */
.bg-decoration {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 0;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: float 8s ease-in-out infinite;
}

.circle-1 {
  width: 400px;
  height: 400px;
  background: #3b82f6;
  top: -100px;
  right: -100px;
  animation-delay: 0s;
}

.circle-2 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -50px;
  left: -50px;
  animation-delay: 2s;
}

.circle-3 {
  width: 200px;
  height: 200px;
  background: linear-gradient(135deg, #4ecdc4 0%, #45b7d1 100%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-20px) scale(1.05);
  }
}

.star {
  position: absolute;
  font-size: 24px;
  opacity: 0.3;
  animation: twinkle 3s ease-in-out infinite;
}

.star-1 { top: 10%; left: 15%; animation-delay: 0s; }
.star-2 { top: 20%; right: 20%; animation-delay: 0.5s; }
.star-3 { top: 60%; left: 10%; animation-delay: 1s; }
.star-4 { top: 70%; right: 15%; animation-delay: 1.5s; }
.star-5 { bottom: 15%; left: 50%; animation-delay: 2s; }

@keyframes twinkle {
  0%, 100% { opacity: 0.3; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.2); }
}

.main-wrapper {
  max-width: 1000px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 40px;
}

/* 卡片通用样式 */
.section-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.dark-mode .section-card {
  background: rgba(45, 55, 72, 0.85);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.section-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.2);
}

/* 个人信息卡片 */
.profile-card {
  text-align: center;
  position: relative;
  overflow: hidden;
}

.card-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 6px;
  background: #3b82f6;
}

.avatar-section {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.avatar-glow {
  position: absolute;
  width: 180px;
  height: 180px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(102, 126, 234, 0.3) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: pulse-glow 2s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 1; transform: translate(-50%, -50%) scale(1.1); }
}

.avatar-border {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background: #3b82f6;
  padding: 4px;
  position: relative;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
}

.avatar-inner {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.3s ease;
}

.dark-mode .avatar-inner {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

.avatar-emoji {
  font-size: 56px;
  animation: bounce 2s ease-in-out infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.status-badge {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(72, 187, 120, 0.4);
}

.profile-name {
  font-size: 40px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 12px;
  letter-spacing: 2px;
}

.profile-tagline {
  font-size: 16px;
  color: #718096;
  margin-bottom: 32px;
  transition: color 0.3s ease;
}

.dark-mode .profile-tagline {
  color: #a0aec0;
}

.profile-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  margin-bottom: 32px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 24px;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  background: #3b82f6;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  font-size: 14px;
  color: #718096;
  transition: color 0.3s ease;
}

.dark-mode .stat-label {
  color: #a0aec0;
}

.stat-divider {
  font-size: 24px;
  color: #cbd5e0;
  opacity: 0.5;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

.social-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border: 2px solid rgba(102, 126, 234, 0.3);
  border-radius: 50px;
  text-decoration: none;
  transition: all 0.3s ease;
}

.dark-mode .social-btn {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  border-color: rgba(102, 126, 234, 0.4);
}

.social-btn:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
  background: #3b82f6;
  border-color: #667eea;
}

.social-btn:hover .social-icon,
.social-btn:hover .social-text {
  color: white !important;
}

.social-icon {
  font-size: 20px;
  color: #667eea;
  transition: color 0.3s ease;
}

.social-text {
  font-size: 14px;
  font-weight: 600;
  color: #4a5568;
  transition: color 0.3s ease;
}

.dark-mode .social-text {
  color: #a0aec0;
}

/* 章节标题 */
.section-header {
  margin-bottom: 24px;
  text-align: center;
}

.section-title {
  font-size: 24px;
  font-weight: 800;
  color: #2d3748;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin: 0;
  transition: color 0.3s ease;
}

.dark-mode .section-title {
  color: #f7fafc;
}

.title-icon {
  font-size: 28px;
}

/* 技能卡�?*/
.skills-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.skill-category {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-title {
  font-size: 16px;
  font-weight: 700;
  color: #4a5568;
  margin: 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #e2e8f0;
  position: relative;
  transition: color 0.3s ease;
}

.category-title::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 60px;
  height: 2px;
  background: linear-gradient(90deg, #667eea 0%, transparent 100%);
}

.dark-mode .category-title {
  color: #cbd5e0;
  border-bottom-color: #4a5568;
}

.skill-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.skill-badge {
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--skill-color) 0%, rgba(0, 0, 0, 0.05) 100%);
  border: 2px solid var(--skill-color);
  border-radius: 50px;
  font-size: 14px;
  font-weight: 600;
  color: white;
  transition: all 0.3s ease;
  cursor: default;
}

.dark-mode .skill-badge {
  background: linear-gradient(135deg, var(--skill-color) 0%, rgba(0, 0, 0, 0.2) 100%);
}

.skill-badge:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 8px 20px var(--skill-color);
}

/* 时间�?*/
.timeline {
  position: relative;
  max-width: 700px;
  margin: 0 auto;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 3px;
  height: 100%;
  background: linear-gradient(to bottom, rgba(59, 130, 246, 0), #3b82f6, rgba(59, 130, 246, 0));
}

.timeline-item {
  position: relative;
  padding: 20px;
  margin-bottom: 32px;
}

.timeline-item:nth-child(odd) {
  text-align: right;
  padding-right: calc(50% + 48px);
}

.timeline-item:nth-child(even) {
  text-align: left;
  padding-left: calc(50% + 48px);
}

.timeline-dot {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #3b82f6;
  box-shadow: 0 0 0 6px rgba(102, 126, 234, 0.2);
  animation: pulse-dot 2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { box-shadow: 0 0 0 6px rgba(102, 126, 234, 0.2); }
  50% { box-shadow: 0 0 0 10px rgba(102, 126, 234, 0.1); }
}

.timeline-content {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.dark-mode .timeline-content {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-color: rgba(102, 126, 234, 0.3);
}

.timeline-content:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.timeline-year {
  font-size: 20px;
  font-weight: 800;
  background: #3b82f6;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.timeline-title {
  font-size: 18px;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 8px;
  transition: color 0.3s ease;
}

.dark-mode .timeline-title {
  color: #f7fafc;
}

.timeline-desc {
  font-size: 14px;
  color: #718096;
  line-height: 1.6;
  margin: 0;
  transition: color 0.3s ease;
}

.dark-mode .timeline-desc {
  color: #a0aec0;
}

/* 兴趣爱好 */
.hobbies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.hobby-item {
  text-align: center;
  padding: 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 16px;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.dark-mode .hobby-item {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.hobby-item:hover {
  transform: translateY(-6px) scale(1.05);
  border-color: rgba(102, 126, 234, 0.5);
  box-shadow: 0 12px 28px rgba(102, 126, 234, 0.2);
}

.hobby-icon {
  font-size: 48px;
  margin-bottom: 12px;
  animation: bounce 2s ease-in-out infinite;
}

.hobby-name {
  font-size: 18px;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 8px;
  transition: color 0.3s ease;
}

.dark-mode .hobby-name {
  color: #f7fafc;
}

.hobby-desc {
  font-size: 14px;
  color: #718096;
  line-height: 1.5;
  margin: 0;
  transition: color 0.3s ease;
}

.dark-mode .hobby-desc {
  color: #a0aec0;
}

/* 响应�?*/
@media (max-width: 768px) {
  .about-page {
    padding: 20px 15px;
  }

  .main-wrapper {
    gap: 24px;
  }

  .section-card {
    padding: 24px 20px;
  }

  .profile-name {
    font-size: 32px;
  }

  .profile-tagline {
    font-size: 14px;
  }

  .profile-stats {
    gap: 16px;
  }

  .stat-value {
    font-size: 24px;
  }

  .social-links {
    gap: 12px;
  }

  .social-btn {
    padding: 10px 16px;
  }

  .section-title {
    font-size: 20px;
  }

  .timeline::before {
    left: 16px;
  }

  .timeline-item:nth-child(odd),
  .timeline-item:nth-child(even) {
    text-align: left;
    padding-left: 48px;
    padding-right: 12px;
  }

  .timeline-dot {
    left: 16px;
  }

  .hobbies-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .hobby-item {
    padding: 20px 16px;
  }
}

@media (max-width: 480px) {
  .hobbies-grid {
    grid-template-columns: 1fr;
  }

  .social-links {
    flex-direction: column;
  }

  .social-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
