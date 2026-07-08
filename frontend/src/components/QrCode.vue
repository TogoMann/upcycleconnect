<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import QRCode from 'qrcode'

const props = defineProps<{
  value: string
  size?: number
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)

async function render() {
  if (canvasRef.value && props.value) {
    await QRCode.toCanvas(canvasRef.value, props.value, {
      width: props.size || 160,
      margin: 1,
      color: {
        dark: '#086a35',
        light: '#ffffff',
      },
    })
  }
}

onMounted(render)
watch(() => props.value, render)
</script>

<template>
  <div class="qrcode-container">
    <canvas ref="canvasRef"></canvas>
  </div>
</template>

<style scoped>
.qrcode-container {
  display: flex;
  justify-content: center;
  padding: 10px;
  background: white;
  border-radius: 8px;
}
</style>
