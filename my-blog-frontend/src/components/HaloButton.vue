<template>
  <div class="container" :class="sizeClass">
    <button class="halo-button">{{ content }}</button>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  content: {
    type: String,
    required: true
  },
  size: {
    type: String,
    required: true,
    validator: (value) => ['small', 'large'].includes(value)
  }
});

const sizeClass = computed(() => props.size);
</script>

<style scoped>
.container {
  position: relative;
  display: inline-block;
}

.container.small {
  width: 50px;
  height: 20px;
}

.container.large {
  width: 80px;
  height: 32px;
}

.halo-button {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  text-align: center;
  color: #fff;
  text-decoration: none;
  font-family: sans-serif;
  box-sizing: border-box;
  background: linear-gradient(90deg, #03a9f4, #e56bb0, #ffeb3b, #03a9f4);
  border: none;
  border-radius: 15px;
  cursor: pointer;
  z-index: 1;
  background-size: 400%;
  padding: 0;
  margin: 0;
  outline: none;
  transition: all 0.3s ease;
}

.container.small .halo-button {
  font-size: 8px;
  line-height: 10px;
  border-radius: 5px;
}

.container.large .halo-button {
  font-size: 12px;
  line-height: 30px;
  border-radius: 12px;
}

.halo-button:hover {
  animation: animate 8s linear infinite;
}

.halo-button::before {
  content: '';
  position: absolute;
  inset: -5px;
  z-index: -1;
  background: linear-gradient(90deg, #03a9f4, #e56bb0, #ffeb3b, #03a9f4);
  background-size: 400%;
  border-radius: 20px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.container.small .halo-button::before {
  border-radius: 10px;
}

.container.large .halo-button::before {
  border-radius: 20px;
}

.halo-button:hover::before {
  filter: blur(10px);
  opacity: 1;
  animation: animate 8s linear infinite;
}

@keyframes animate {
  0% {
    background-position: 0%;
  }
  100% {
    background-position: 400%;
  }
}
</style>