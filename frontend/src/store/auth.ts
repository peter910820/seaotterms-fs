import { defineStore } from "pinia";
import { ref } from "vue";

import type { UserType } from "@/types/user";

export const useAuthStore = defineStore("auth", () => {
  const showLoginModal = ref(false);
  // 依 cookie session_id 判斷登入，由 router 同步
  const hasSession = ref(false);
  const user = ref<UserType>();

  const setSession = (value: boolean) => {
    hasSession.value = value;
  };

  const setUser = (data: UserType) => {
    user.value = data;
  };

  const clearUser = () => {
    user.value = undefined;
  };

  const openLoginModal = () => {
    showLoginModal.value = true;
  };

  const closeLoginModal = () => {
    showLoginModal.value = false;
  };

  return {
    showLoginModal,
    hasSession,
    user,
    setSession,
    setUser,
    clearUser,
    openLoginModal,
    closeLoginModal,
  };
});
