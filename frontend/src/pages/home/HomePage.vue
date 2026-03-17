<script setup lang="ts">
import { ref } from 'vue'
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F4] font-sans text-gray-800">
    <nav class="flex items-center justify-between px-8 py-6 max-w-7xl mx-auto">
      <div class="text-2xl font-semibold text-[#3D7A4D]">UpCycleConnect</div>

      <div class="hidden md:flex space-x-8 text-sm font-medium text-[#7A9C78]">
        <a href="#" class="hover:underline">Accueil</a>
        <a href="#" class="hover:underline">Prestations</a>
        <a href="#" class="hover:underline">Évènements</a>
        <a href="#" class="hover:underline">Forum</a>
        <a href="#" class="hover:underline">À propos</a>
      </div>

      <button
        class="bg-[#065F36] text-white px-6 py-2 rounded-md text-sm font-medium hover:bg-opacity-90 transition"
      >
        S'inscrire / Se connecter
      </button>
    </nav>

    <header class="text-center py-16 px-4">
      <h1 class="text-5xl md:text-6xl font-light text-[#333] mb-12">
        L'art de transformer vos objets.
      </h1>

      <div class="max-w-7xl mx-auto overflow-hidden rounded-lg shadow-xl">
        <img
          src="https://images.unsplash.com/photo-1500382017468-9049fed747ef?auto=format&fit=crop&q=80&w=2000"
          alt="Paysage côtier vert"
          class="w-full h-[400px] object-cover"
        />
      </div>

      <p class="mt-16 text-3xl font-light text-[#333]">
        Valorisons nos ressources pour un avenir plus
        <span class="text-[#3D7A4D] font-normal">durable</span>.
      </p>
    </header>

    <!-- <section
      class="max-w-7xl mx-auto px-8 py-12 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8"
    >
      <ServiceCard
        v-for="service in services"
        :key="service.title"
        :title="service.title"
        :description="service.description"
        :image="service.image"
      />
    </section> -->
  </div>
</template>

<!-- <script setup>
import ServiceCard from '@components/ServiceCard.vue'

const services = [
  {
    title: 'RECYCLER',
    description:
      'Donnez une nouvelle utilité à vos matériaux en fin de vie grâce à notre réseau de collecte et de tri optimisé.',
    image:
      'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?auto=format&fit=crop&q=80&w=500',
  },
  {
    title: 'RÉPARER',
    description:
      "Prolongez la durée de vie de vos objets en accédant à l'expertise de nos artisans spécialisés dans la restauration.",
    image:
      'https://images.unsplash.com/photo-1581092160562-40aa08e78837?auto=format&fit=crop&q=80&w=500',
  },
  {
    title: 'TRANSFORMER',
    description: 'Métamorphosez vos anciens produits en pièces uniques et design.',
    image:
      'https://images.unsplash.com/photo-1504148455328-c376907d081c?auto=format&fit=crop&q=80&w=500',
  },
  {
    title: 'SE FORMER',
    description:
      'Apprenez les gestes de demain avec nos ateliers et développez vos compétences pour devenir un acteur du changement.',
    image:
      'https://images.unsplash.com/photo-1542601906990-b4d3fb778b09?auto=format&fit=crop&q=80&w=500',
  },
]
</script> -->
