<template>
  <div class="toy-wrapper">
    <a href="/tools" class="back-btn">← 返回博客</a>
    <iframe
      src="https://playstarfling.com"
      frameborder="0"
      allow="fullscreen"
      class="toy-iframe"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })
useSeoMeta({
  title: 'Starfling · 太空弹弓游戏',
  description: '一个用单个 HTML 文件实现的太空轨道弹弓游戏，一键发射，无限循环。',
})
</script>

<style scoped>
.toy-wrapper {
  position: fixed;
  inset: 0;
  background: #000;
}
.toy-iframe {
  width: 100%;
  height: 100%;
  border: none;
  display: block;
}
.back-btn {
  position: fixed;
  top: 16px;
  left: 16px;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  text-decoration: none;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  transition: background 0.2s;
}
.back-btn:hover {
  background: rgba(0, 0, 0, 0.85);
}
</style>
