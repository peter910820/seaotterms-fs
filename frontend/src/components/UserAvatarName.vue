<template>
  <div class="user-avatar-name d-flex align-center gap-3" :class="rootClass">
    <v-avatar color="primary" :size="size">
      <v-img v-if="avatarUrl" :src="avatarUrl" cover />
      <v-icon v-else :size="iconSize">{{ fallbackIcon }}</v-icon>
    </v-avatar>
    <span class="user-avatar-name__text text-body2 text-medium-emphasis" :class="nameClass">{{ username }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    avatarUrl?: string;
    username: string;
    size?: number;
    rootClass?: string;
    nameClass?: string;
  }>(),
  {
    size: 36,
    rootClass: "",
    nameClass: "",
  },
);

const fallbackIcon = "mdi-account";
const iconSize = computed(() => Math.round((props.size ?? 36) * 0.56));
</script>

<style scoped>
.user-avatar-name {
  min-width: 0;
}
.user-avatar-name__text {
  min-width: 0;
  margin-left: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
