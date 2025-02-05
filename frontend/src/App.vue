<template>
  <Header />
  <main>
    <HomeView />
  </main>
  <Footer />
</template>

<script setup>
import { onMounted } from 'vue';
import { useSensorStore } from '@/stores/sensors'

import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import HomeView from '@/views/HomeView.vue'

const sensors = useSensorStore()

onMounted(() => {
  sensors.fetchData()
})
</script>