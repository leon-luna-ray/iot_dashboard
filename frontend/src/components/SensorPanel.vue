<template>
    <section class="container">
        <ul class="grid md:grid-cols-2 lg:grid-cols-3 gap-[2rem] lg:gap-[4rem]">
            <li v-if="sensors.data?.battery?.value" class="card">
                <h2 class="h5 capitalize">Battery</h2>
                <p class="h4">{{ sensors.data.battery.value }}%</p>
            </li>
            <li v-if="sensors.timestemp" class="card">
                <h2 class="h5 capitalize">Updated at</h2>
                <p class="h4">{{ sensors.timestemp }}</p>
            </li>
            <li v-if="sensors.temperature" class="card">
                <h2 class="h5 capitalize">Temperature</h2>
                <p class="h4">{{ sensors.temperature }}</p>
            </li>
            <li v-if="sensors.data?.humidity?.value" class="card">
                <h2 class="h5 capitalize">Humidity</h2>
                <p class="h4">{{ Math.round(sensors.data.humidity.value) }}%</p>
            </li>
            <li v-if="sensors.data?.co2?.value" class="card">
                <h2 class="h5 capitalize">Co2</h2>
                <p class="h4">{{ sensors.data.co2.value }} ppm</p>
            </li>
            <li v-if="sensors.data?.pm25?.value" class="card">
                <h2 class="h5 capitalize">pm25</h2>
                <p class="h4">{{ sensors.data.pm25.value }}</p>
            </li>
            <li v-if="sensors.data?.pm10?.value" class="card">
                <h2 class="h5 capitalize">pm10</h2>
                <p class="h4">{{ sensors.data.pm10.value }}</p>
            </li>
        </ul>
    </section>
</template>


<script setup>
import { useSensorStore } from '@/stores/sensors';

const sensors = useSensorStore();
</script>