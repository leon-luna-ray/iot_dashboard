import axios from 'axios';
import { ref, computed, onMounted } from 'vue';
import { defineStore } from 'pinia';
import { setLocalStorage } from '@/utils/localStorage.js';

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
    });
    const timestemp = computed(() => {
        if (data.value?.timestamp?.value) {
            return new Date(data.value.timestamp.value).toLocaleTimeString();
        }
        return null;
    });

    // Methods
    const toggleUnits = async () => {
        units.value = units.value === 'imperial' ? 'metric' : 'imperial';
        await setLocalStorage('units', units.value);
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

    const setUnitsFromLocalStorage = () => {
        const value = localStorage.getItem('units');
        if (value) {
            units.value = JSON.parse(value);
        }
    }

    onMounted(() => {
        setUnitsFromLocalStorage();
    });
    return { data, deviceInfo, temperature, timestemp, fetchData, toggleUnits };
});
