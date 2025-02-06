<template>
    <section class="container">
        <ul class="grid md:grid-cols-2 lg:grid-cols-3 gap-[2rem] lg:gap-[4rem]">
            <li v-if="sensors.temperature" class="card">
                <h2 class="h4 capitalize">Temperature</h2>
                <p class="h3">{{ sensors.temperature }}</p>
            </li>
            <!-- <li v-for="item in gridItems" class="card" :key="item.key">
                <h2 class="h4 capitalize">{{ item.key }}</h2>
                <p class="h3">{{ item.value.value }}</p>
            </li> -->
        </ul>
    </section>
</template>


<script setup>
import { useSensorStore } from '@/stores/sensors';

const sensors = useSensorStore();
</script>