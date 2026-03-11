<template>
  <div class="error-page">
    <v-card class="error-card mx-auto" variant="outlined">
      <v-card-title class="error-title d-flex align-center text-error">
        <v-icon start size="large">mdi-alert-circle</v-icon>
        錯誤
      </v-card-title>
      <v-divider />
      <v-card-text class="error-content">
        <p class="error-message">{{ errorMessage }}</p>
        <v-btn to="/" color="primary" variant="tonal" size="large" class="mt-4">
          <v-icon start>mdi-home</v-icon>
          返回首頁
        </v-btn>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

const errorMessage = ref("發生錯誤！");

// 暫時保留舊邏輯，一樣從sessionStorage抓error msg
onMounted(() => {
  const msg = sessionStorage.getItem("msg");
  if (msg) {
    errorMessage.value = msg;
    sessionStorage.removeItem("msg");
  }
});
</script>

<style scoped>
.error-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 70vh;
  padding: 24px;
}
.error-card {
  border-radius: 16px;
  max-width: 720px;
  width: 100%;
}
.error-title {
  font-size: 1.75rem;
  font-weight: 600;
  padding: 28px 32px;
}
.error-content {
  padding: 28px 32px 32px;
}
.error-message {
  font-size: 1.25rem;
  line-height: 1.6;
  color: rgb(var(--v-theme-on-surface));
  word-break: break-word;
  margin: 0;
}
</style>
