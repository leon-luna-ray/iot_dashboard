<template>
  <div class="flex flex-col justify-between h-screen">
    <div>
      <Header />
      <main>
        <HomeView />
      </main>
    </div>
    <Footer />
  </div>
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