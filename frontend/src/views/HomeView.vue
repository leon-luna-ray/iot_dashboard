<template>
    <div class="container flex-col-2">
        <h1 class="text-center h2">Home Stats</h1>
        <ul v-if="gridItems.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-[2rem] lg:gap-[4rem]">
            <li v-for="item in gridItems" class="text-center flex-col-1 border py-[1rem]" :key="item.key">
                <h2 class="h4 capitalize">{{ item.key }}</h2>
                <p class="h3">{{ item.value.value }}</p>
            </li>
        </ul>
    </div>
</template>

<script setup>
import { computed } from 'vue';
import { useSensorStore } from '@/stores/sensors';

const sensors = useSensorStore();

// Just for testing
const gridItems = computed(() => {
    if(!sensors.data) return [];
    return Object.entries(sensors.data).map(([key, value]) => ({ key, value }));
});
</script>