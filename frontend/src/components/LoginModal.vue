<template>
  <v-dialog v-model="dialog" max-width="600" persistent @click:outside="dialog = false">
    <v-card class="login-card">
      <v-card-title class="login-title">
        <v-icon size="large" class="mr-3">mdi-login</v-icon>
        登入
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-8">
        <v-form @submit.prevent="handleSubmit">
          <v-text-field
            v-model="request.username"
            label="使用者名稱"
            prepend-inner-icon="mdi-account-circle"
            variant="outlined"
            required
            class="mb-4"
            density="comfortable"
            :rules="[(v) => !!v || '請輸入使用者名稱']"
          />

          <v-text-field
            v-model="request.password"
            label="密碼"
            type="password"
            prepend-inner-icon="mdi-lock"
            variant="outlined"
            required
            class="mb-2"
            density="comfortable"
            :rules="[(v) => !!v || '請輸入密碼']"
          />
        </v-form>
      </v-card-text>

      <v-divider />

      <v-card-actions class="pa-6">
        <v-spacer />
        <v-btn variant="text" @click="dialog = false" class="mr-2">取消</v-btn>
        <v-btn color="primary" variant="elevated" :loading="loading" @click="handleSubmit" size="large">
          登入
          <v-icon end>mdi-send</v-icon>
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { useAuthStore } from "@/store/auth";
import type { LoginResponseData, ResponseType } from "@/types/response";
import type { UserType } from "@/types/user";

interface Props {
  modelValue: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits(["update:modelValue"]);

const router = useRouter();
const authStore = useAuthStore();
const { closeLoginModal, setSession, setUser } = authStore;

const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value),
});

const request = ref({
  username: "",
  password: "",
});

const loading = ref(false);

const LOGIN_API_BASE = `${import.meta.env.VITE_API_URL}/login`;

const loginDataToUser = (data: LoginResponseData): UserType => ({
  username: data.username,
  email: data.email,
  avatar: data.avatar,
  management: data.isAdmin,
  createdAt: new Date(data.createdAt),
  createName: "",
});

const handleSubmit = async () => {
  if (!request.value.username || !request.value.password) return;

  loading.value = true;
  try {
    const response = await axios.post<ResponseType<LoginResponseData>>(LOGIN_API_BASE, request.value);
    sessionStorage.setItem("msg", response.data.message);
    setUser(loginDataToUser(response.data.data));
    setSession(true);
    closeLoginModal();
    dialog.value = false;
  } catch (error) {
    closeLoginModal();
    dialog.value = false;
    if (axios.isAxiosError(error)) {
      sessionStorage.setItem("msg", `${error.response?.status}: ${error.response?.data?.message ?? error.message}`);
      router.push("/error");
    } else {
      sessionStorage.setItem("msg", String(error));
      router.push("/error");
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-card {
  background-color: rgb(var(--v-theme-background));
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}

.login-title {
  display: flex;
  align-items: center;
  font-size: 1.5rem;
  font-weight: 600;
  color: rgb(var(--v-theme-info));
  padding: 24px 32px;
}
</style>
