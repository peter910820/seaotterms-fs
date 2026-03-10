<template>
  <div class="main-layout">
    <NavBar />
    <LoginModal v-model="showLoginModal" />
    <v-main class="fs-main">
      <v-container fluid class="pa-4">
        <router-view />
      </v-container>
    </v-main>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";

import NavBar from "@/components/NavBar.vue";
import LoginModal from "@/components/LoginModal.vue";
import { useLoginModal } from "@/store/loginModal";

const loginModalStore = useLoginModal();
const { showLoginModal } = storeToRefs(loginModalStore);
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
}
.fs-main {
  background: rgb(var(--v-theme-background));
}
</style>
