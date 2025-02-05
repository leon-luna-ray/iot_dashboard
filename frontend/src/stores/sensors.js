import axios from 'axios';

import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useSensorStore = defineStore('sensors', () => {
    // State
    const data = ref(null);
    const deviceInfo = ref(true);

    // Methods
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
    return { data, deviceInfo, fetchData };
});
