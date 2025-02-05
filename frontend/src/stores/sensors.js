import axios from 'axios';
import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useSensorStore = defineStore('sensors', () => {
    const data = ref(null);
    const loading = ref(true);

    const fetchData = async () => {
        loading.value = true;
        try {
            const response = await axios.get('/sensors');
            data.value = response.data;
            console.log('Fetched data:', data.value);
        } catch (error) {
            console.error('Error fetching data:', error);
        } finally {
            loading.value = false;
        }
    }
    return { data, loading, fetchData };
});
