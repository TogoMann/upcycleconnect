<script setup lang="ts">
import { onMounted, ref } from 'vue'
import JsBarcode from 'jsbarcode'

const props = defineProps<{
  value: string
  format?: string
  width?: number
  height?: number
  displayValue?: boolean
}>()

const barcodeRef = ref<HTMLCanvasElement | null>(null)

onMounted(() => {
  if (barcodeRef.value) {
    JsBarcode(barcodeRef.value, props.value, {
      format: props.format || "CODE128",
      width: props.width || 2,
      height: props.height || 100,
      displayValue: props.displayValue !== undefined ? props.displayValue : true,
      lineColor: "#086a35", // UpcycleConnect Green
    })
  }
})
</script>

<template>
  <div class="barcode-container">
    <canvas ref="barcodeRef"></canvas>
  </div>
</template>

<style scoped>
.barcode-container {
  display: flex;
  justify-content: center;
  padding: 10px;
  background: white;
  border-radius: 8px;
}
</style>
