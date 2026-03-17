<template>
  <div class="flex flex-col">
    <div class="aspect-video w-full overflow-hidden mb-6 bg-gray-200">
      <img :src="image" :alt="title" class="w-full h-full object-cover grayscale-[20%]" />
    </div>
    <h3 class="font-bold text-lg mb-4 tracking-wider">{{ title }}</h3>
    <p class="text-gray-600 text-sm leading-relaxed">
      {{ description }}
    </p>
  </div>
</template>

<script setup>
defineProps(['title', 'description', 'image'])
</script>
