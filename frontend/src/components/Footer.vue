<template>
    <footer class="py-[2rem]">
        <div class="container text-center">
            <p class="text-[0.875rem]">&copy;{{ year }} <a href="https://www.rayluna.dev" target="_blank" class="underline">rayluna.dev</a></p>
        </div>
    </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>