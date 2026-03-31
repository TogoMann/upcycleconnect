<script setup lang="ts">
defineProps<{
    title: string
    description: string
    icon: string
}>()
</script>

<template>
    <div class="service-card">
        <div class="icon-container">
            <span v-if="icon === 'recycle'">♻️</span>
            <span v-else-if="icon === 'tools'">🛠️</span>
            <span v-else-if="icon === 'shop'">🛍️</span>
            <span v-else>✨</span>
        </div>
        <h3>{{ title }}</h3>
        <p>{{ description }}</p>
    </div>
</template>

<style scoped>
.service-card {
    background: white;
    padding: 30px;
    border-radius: 12px;
    text-align: center;
    border: 1px solid #d7ece1;
    transition: transform 0.3s ease;
}
.service-card:hover {
    transform: translateY(-5px);
}
.icon-container {
    width: 50px;
    height: 50px;
    background: #f8f5ee;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 15px;
    font-size: 1.5rem;
}
h3 {
    color: #086a35;
    margin-bottom: 10px;
}
p {
    color: #353535;
    font-size: 0.9rem;
}
</style>
