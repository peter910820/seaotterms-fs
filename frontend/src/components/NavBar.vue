<template>
  <v-app-bar color="surface" density="compact" elevation="1">
    <v-app-bar-nav-icon v-if="mobile" @click="drawer = !drawer" aria-label="選單" />
    <v-toolbar-title class="text-primary font-weight-bold d-flex align-center">
      <v-icon start color="primary">mdi-folder-network</v-icon>
      海獺的資源伺服器
    </v-toolbar-title>
    <v-spacer />
    <v-btn
      v-for="item in navItems"
      :key="item.to"
      :to="item.to"
      :active="isActive(item.to)"
      variant="text"
      color="primary"
      class="d-none d-md-inline-flex"
    >
      <v-icon start size="small">{{ item.icon }}</v-icon>
      {{ item.title }}
    </v-btn>
    <div v-if="hasSession && user" class="d-none d-md-flex align-center gap-3 mr-3 nav-bar-user">
      <v-avatar color="primary" size="36">
        <v-img v-if="avatarUrl" :src="avatarUrl" cover />
        <v-icon v-else size="20">mdi-account</v-icon>
      </v-avatar>
      <span class="text-body2 text-medium-emphasis nav-bar-user-name">{{ user.username }}</span>
    </div>
    <v-btn v-if="hasSession" variant="text" color="primary" class="d-none d-md-inline-flex" @click="handleLogout">
      <v-icon start size="small">mdi-logout</v-icon>
      登出
    </v-btn>
    <v-btn v-else variant="text" color="primary" class="d-none d-md-inline-flex" @click="openLoginModal">
      <v-icon start size="small">mdi-login</v-icon>
      登入
    </v-btn>
  </v-app-bar>
  <v-navigation-drawer v-model="drawer" :mobile-breakpoint="960" temporary location="start" class="fs-drawer">
    <div v-if="hasSession && user" class="pa-3 d-flex align-center gap-3 nav-drawer-user">
      <v-avatar color="primary" size="40">
        <v-img v-if="avatarUrl" :src="avatarUrl" cover />
        <v-icon v-else size="24">mdi-account</v-icon>
      </v-avatar>
      <span class="text-body2 text-medium-emphasis nav-drawer-user-name">{{ user.username }}</span>
    </div>
    <v-divider v-if="hasSession && user" class="my-0" />
    <v-list density="compact" nav>
      <v-list-item
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        :active="isActive(item.to)"
        :prepend-icon="item.icon"
        :title="item.title"
        active-color="primary"
        rounded="lg"
      />
      <v-divider class="my-2" />
      <v-list-item v-if="hasSession" :prepend-icon="'mdi-logout'" title="登出" rounded="lg" @click="handleLogout" />
      <v-list-item v-else :prepend-icon="'mdi-login'" title="登入" rounded="lg" @click="openLoginModal" />
    </v-list>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import { useDisplay } from "vuetify";
import { storeToRefs } from "pinia";
import Cookies from "js-cookie";

import { useAuthStore } from "@/store/auth";

const route = useRoute();
const { mobile } = useDisplay();
const drawer = ref(false);

const authStore = useAuthStore();
const { hasSession, user } = storeToRefs(authStore);

const avatarUrl = computed(() => user.value?.avatar?.trim() || undefined);

const openLoginModal = () => {
  drawer.value = false;
  authStore.openLoginModal();
};

const handleLogout = () => {
  drawer.value = false;
  Cookies.remove("session_id");
  authStore.setSession(false);
  authStore.clearUser();
};

const navItems = [
  { to: "/", title: "首頁", icon: "mdi-home" },
  { to: "/folder", title: "檔案夾", icon: "mdi-folder-open" },
  { to: "/markdown-writer", title: "Markdown預覽器", icon: "mdi-text-box-edit-outline" },
];

// 點亮navbar用(判斷是否在當前頁面)
const isActive = (path: string) => {
  if (path === "/") return route.path === "/" || route.name === "home";
  return route.path === path;
};
</script>

<style scoped>
.nav-bar-user,
.nav-drawer-user {
  min-width: 0;
}
.nav-bar-user-name,
.nav-drawer-user-name {
  min-width: 0;
  margin-left: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.nav-bar-user-name {
  max-width: 120px;
}
</style>
