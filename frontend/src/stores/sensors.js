import axios from 'axios';
import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useSensorsStore = defineStore('sensors', () => {
    const data = ref(null);
    const loading = ref(true);

    const fetchData = async () => {
        loading.value = true;
        try {
            const response = await fetch('http://localhost:3000/sensors');
            data.value = await response.json();
        } catch (error) {
            console.error('Error fetching data:', error);
        } finally {
            loading.value = false;
        }
    }
    return { data, loading, fetchData };
});
