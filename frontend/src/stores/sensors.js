import axios from 'axios';
import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

export const useSensorStore = defineStore('sensors', () => {
    // State
    const data = ref(null);
    const deviceInfo = ref(true);
    // // TODO: set to local storage default to metric
    const units = ref('imperial');

    const temperature = computed(() => {
        if (data.value?.temperature?.value && units.value === 'imperial') {
            return Math.round(data.value.temperature.value * 9 / 5 + 32) + '°F';
        }
        else if (data.value?.temperature?.value) {
            return Math.round(data.value.temperature.value) + '°C';
        }
        return null;
    })

    // Methods
    const toggleUnits = () => {
        // Todo set to local storage
        units.value = units.value === 'metric' ? 'imperial' : 'metric';
    }
    const fetchData = async () => {
        try {
            const res = await axios.get('/sensors');
            // Filter data for single device setup
            data.value = res.data?.devices[0]?.data || null;
            deviceInfo.value = res.data?.devices[0]?.info || null;

        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }
    return { data, deviceInfo, temperature, fetchData };
});
