<template>
    <header class="container">
        <div class="flex justify-between items-center py-[1rem]">

            <div>
                <h1 class="h5">IoT Dashboard</h1>
            </div>
            <div>
                <button class="btn w-full" @click="sensors.toggleUnits">°C/°F</button>
            </div>
        </div>
    </header>
</template>

<script setup>
import { useSensorStore } from '@/stores/sensors';

const sensors = useSensorStore();
</script>