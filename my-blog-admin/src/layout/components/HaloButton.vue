<template>
    <div class="container" :class="[sizeClass, typeClass, { 'loading': loading, 'disabled': disabled }]"
         @click="handleClick">
        <button class="halo-button" :disabled="disabled || loading">
            <el-icon v-if="loading" class="is-loading">
                <Loading />
            </el-icon>
            <el-icon v-else-if="icon">
                <component :is="icon" />
            </el-icon>
            <span v-if="content || $slots.default" class="button-content">
                <slot>{{ content }}</slot>
            </span>
        </button>
    </div>
</template>

<script setup>
import { computed } from 'vue';
import { Loading } from '@element-plus/icons-vue';

const props = defineProps({
    content: {
        type: String,
        default: ''
    },
    size: {
        type: String,
        default: 'large',
        validator: (value) => ['small', 'medium', 'large'].includes(value)
    },
    type: {
        type: String,
        default: 'primary',
        validator: (value) => ['default', 'primary', 'success', 'warning', 'danger'].includes(value)
    },
    loading: {
        type: Boolean,
        default: false
    },
    disabled: {
        type: Boolean,
        default: false
    },
    icon: {
        type: [String, Object],
        default: null
    }
});

const emit = defineEmits(['click']);

const sizeClass = computed(() => props.size);
const typeClass = computed(() => `type-${props.type}`);

const handleClick = (event) => {
    if (!props.disabled && !props.loading) {
        emit('click', event);
    }
};
</script>

<style scoped>
.container {
    position: relative;
    display: inline-block;
}

/* 尺寸 */
.container.small {
    width: 70px;
    height: 28px;
}

.container.medium {
    width: 90px;
    height: 32px;
}

.container.large {
    width: 110px;
    height: 36px;
}

/* 类型 - 默认 (灰色系) */
.container.type-default .halo-button {
    background: linear-gradient(270deg, #909399, #b0b3b8, #909399);
    background-size: 200%;
}

.container.type-default .halo-button::before {
    background: linear-gradient(270deg, #909399, #b0b3b8, #909399);
    background-size: 200%;
}

/* 类型 - 主要 (蓝色系) */
.container.type-primary .halo-button {
    background: linear-gradient(270deg, #03a9f4, #4fc3f7, #03a9f4);
    background-size: 200%;
}

.container.type-primary .halo-button::before {
    background: linear-gradient(270deg, #03a9f4, #4fc3f7, #03a9f4);
    background-size: 200%;
}

/* 类型 - 成功 (绿色系) */
.container.type-success .halo-button {
    background: linear-gradient(270deg, #4caf50, #8bc34a, #4caf50);
    background-size: 200%;
}

.container.type-success .halo-button::before {
    background: linear-gradient(270deg, #4caf50, #8bc34a, #4caf50);
    background-size: 200%;
}

/* 类型 - 警告 (橙黄色系) */
.container.type-warning .halo-button {
    background: linear-gradient(270deg, #ff9800, #ffcc02, #ff9800);
    background-size: 200%;
}

.container.type-warning .halo-button::before {
    background: linear-gradient(270deg, #ff9800, #ffcc02, #ff9800);
    background-size: 200%;
}

/* 类型 - 危险 (红色系) */
.container.type-danger .halo-button {
    background: linear-gradient(270deg, #f44336, #e57373, #f44336);
    background-size: 200%;
}

.container.type-danger .halo-button::before {
    background: linear-gradient(270deg, #f44336, #e57373, #f44336);
    background-size: 200%;
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
    border: none;
    border-radius: 15px;
    cursor: pointer;
    z-index: 1;
    padding: 0;
    margin: 0;
    outline: none;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
}

.container.small .halo-button {
    font-size: 10px;
    border-radius: 8px;
    padding: 0 8px;
    width: auto;
    min-width: 60px;
}

.container.medium .halo-button {
    font-size: 12px;
    border-radius: 10px;
    padding: 0 10px;
    width: auto;
    min-width: 80px;
}

.container.large .halo-button {
    font-size: 14px;
    border-radius: 12px;
    padding: 0 12px;
    width: auto;
    min-width: 100px;
}

.halo-button .el-icon {
    font-size: 1.1em;
}

.container.small .halo-button .el-icon {
    font-size: 1em;
}

.container.medium .halo-button .halo-button .el-icon {
    font-size: 1.05em;
}

.container.large .halo-button .el-icon {
    font-size: 1.15em;
}

.halo-button:hover:not(.disabled) {
    animation: animate 8s linear infinite;
}

.halo-button::before {
    content: '';
    position: absolute;
    inset: -3px;
    z-index: -1;
    border-radius: 20px;
    opacity: 0;
    transition: opacity 0.3s ease;
}

.container.small .halo-button::before {
    border-radius: 12px;
}

.container.medium .halo-button::before {
    border-radius: 16px;
}

.container.large .halo-button::before {
    border-radius: 20px;
}

.halo-button:hover:not(.disabled)::before {
    filter: blur(10px);
    opacity: 1;
    animation: animate 8s linear infinite;
}

.halo-button.disabled {
    opacity: 0.6;
    cursor: not-allowed;
    filter: grayscale(0.5);
}

.halo-button.disabled:hover {
    animation: none;
}

.halo-button.disabled:hover::before {
    opacity: 0;
    animation: none;
}

.is-loading {
    animation: rotating 2s linear infinite;
}

@keyframes animate {
    0% {
        background-position: 200%;
    }

    100% {
        background-position: 0%;
    }
}

@keyframes rotating {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}

.button-content {
    white-space: nowrap;
}
</style>
