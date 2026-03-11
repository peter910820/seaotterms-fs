<template>
  <div class="folder-page">
    <v-row>
      <v-col cols="12" md="4" lg="3">
        <v-card class="fs-folder-sidebar" variant="outlined">
          <v-card-title class="text-subtitle-1 py-3 d-flex align-center">
            <v-icon start size="small">mdi-folder</v-icon>
            {{ pathLabel }}
          </v-card-title>
          <v-list density="compact" class="py-0">
            <v-list-item
              v-if="parentPath !== null"
              key="parent"
              rounded="lg"
              class="mb-1 folder-item"
              @click="parentPath !== null && loadPath(parentPath)"
            >
              <template #prepend>
                <v-icon size="small">mdi-folder-arrow-up</v-icon>
              </template>
              <v-list-item-title class="text-body2">..</v-list-item-title>
            </v-list-item>
            <v-list-item
              v-for="(name, index) in currentDirectories"
              :key="'dir-' + index"
              rounded="lg"
              class="mb-1 folder-item"
              @click="goIntoFolder(name)"
            >
              <template #prepend>
                <v-icon size="small">mdi-folder-outline</v-icon>
              </template>
              <v-list-item-title class="text-body2">{{ name }}</v-list-item-title>
            </v-list-item>
          </v-list>
          <v-card-text v-if="error" class="text-caption text-error text-center py-2">
            {{ error }}
          </v-card-text>
          <v-card-text
            v-else-if="currentDirectories.length === 0"
            class="text-caption text-medium-emphasis text-center py-2"
          >
            尚無子資料夾
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="8" lg="9">
        <v-card class="fs-file-list" variant="outlined">
          <v-card-title class="text-subtitle-1 py-3 d-flex align-center">
            <v-icon start size="small">mdi-file-document-multiple-outline</v-icon>
            {{ pathLabel }} 檔案
          </v-card-title>
          <v-divider />
          <v-list v-if="currentFiles.length > 0" class="py-0">
            <v-list-item
              v-for="(name, index) in currentFiles"
              :key="'file-' + index"
              :href="fileUrl(name)"
              target="_blank"
              rel="noopener"
              class="file-item"
              rounded="lg"
            >
              <template #prepend>
                <v-icon size="small" color="primary">
                  {{ isImage(name) ? "mdi-file-image-outline" : "mdi-file-outline" }}
                </v-icon>
              </template>
              <v-list-item-title class="text-truncate">{{ name }}</v-list-item-title>
              <template #append>
                <v-icon size="small">mdi-open-in-new</v-icon>
              </template>
            </v-list-item>
          </v-list>
          <v-card-text v-else class="text-medium-emphasis text-center py-8">
            {{ error || "此資料夾尚無檔案" }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";

import { useFileList } from "@/composables/useFileList";

const CDN_BASE = (import.meta.env.VITE_CDN_URL ?? "").replace(/\/$/, "");

const { currentPath, currentDirectories, currentFiles, error, pathLabel, parentPath, loadPath, goIntoFolder } =
  useFileList();

const isImage = (name: string): boolean => /\.(png|jpg|jpeg|gif|webp|svg)$/i.test(name);

const fileUrl = (fileName: string): string => {
  const path = currentPath.value ? `${currentPath.value}/${fileName}` : fileName;
  return path ? `${CDN_BASE}/${path}` : CDN_BASE;
};

onMounted(() => loadPath(""));
</script>

<style scoped>
.folder-page {
  min-height: 60vh;
}
.fs-folder-sidebar,
.fs-file-list {
  border-radius: 12px;
}
.folder-item,
.file-item {
  cursor: pointer;
}
</style>
