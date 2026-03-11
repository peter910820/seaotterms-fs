<template>
  <div class="folder-page">
    <v-row>
      <!-- 左：當前層 API 回傳的資料夾全數顯示（根目錄亦然）；.. 回上層 -->
      <v-col cols="12" md="4" lg="3">
        <v-card class="fs-folder-sidebar" variant="outlined">
          <v-card-title class="text-subtitle-1 py-3 d-flex align-center">
            <v-icon start size="small">mdi-folder</v-icon>
            {{ currentFolderLabel }}
          </v-card-title>
          <v-list density="compact" class="py-0">
            <v-list-item
              v-if="parentPath !== null"
              key="parent"
              rounded="lg"
              class="mb-1 folder-item"
              @click="loadPath(parentPath as string)"
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
          <v-card-text v-if="folderLoadError" class="text-caption text-error text-center py-2">
            {{ folderLoadError }}
          </v-card-text>
          <v-card-text
            v-else-if="currentDirectories.length === 0"
            class="text-caption text-medium-emphasis text-center py-2"
          >
            尚無子資料夾
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 右：當前層 API 回傳的檔案 -->
      <v-col cols="12" md="8" lg="9">
        <v-card class="fs-file-list" variant="outlined">
          <v-card-title class="text-subtitle-1 py-3 d-flex align-center">
            <v-icon start size="small">mdi-file-document-multiple-outline</v-icon>
            {{ currentFolderLabel }} 檔案
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
            {{ fileLoadError || "此資料夾尚無檔案" }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import axios from "axios";
import type { FileResponseData, ResponseType } from "@/types/response";

const API_TIMEOUT_MS = 5000;
const FILE_API_BASE = `${import.meta.env.VITE_API_URL}/file`;

const CDN_BASE = (import.meta.env.VITE_CDN_URL ?? "").replace(/\/$/, "");

const currentPath = ref("");
const currentFiles = ref<string[]>([]);
const currentDirectories = ref<string[]>([]);
const folderLoadError = ref("");
const fileLoadError = ref("");

const currentFolderLabel = computed(() => currentPath.value || "根目錄");

const parentPath = computed(() => {
  const p = currentPath.value;
  if (!p) return null;
  const parts = p.split("/").filter(Boolean);
  if (parts.length <= 1) return "";
  return parts.slice(0, -1).join("/");
});

const goIntoFolder = (name: string) => {
  const next = currentPath.value ? `${currentPath.value}/${name}` : name;
  loadPath(next);
};

const isImage = (name: string): boolean => /\.(png|jpg|jpeg|gif|webp|svg)$/i.test(name);

const fileUrl = (fileName: string): string => {
  const path = currentPath.value ? `${currentPath.value}/${fileName}` : fileName;
  return path ? `${CDN_BASE}/${path}` : CDN_BASE;
};

const fetchFileList = async (path: string): Promise<ResponseType<FileResponseData> | null> => {
  const url = path ? `${FILE_API_BASE}/${path}` : FILE_API_BASE;
  try {
    const response = await axios.get<ResponseType<FileResponseData>>(url, { timeout: API_TIMEOUT_MS });
    if (response?.status === 200 && response.data?.data) return response.data;
    return null;
  } catch {
    return null;
  }
};

const loadPath = async (path: string) => {
  currentPath.value = path;
  fileLoadError.value = "";
  folderLoadError.value = "";
  const result = await fetchFileList(path);
  if (result?.data) {
    currentFiles.value = result.data.files ?? [];
    currentDirectories.value = result.data.directories ?? [];
  } else {
    currentFiles.value = [];
    currentDirectories.value = [];
    fileLoadError.value = "無法載入此資料夾";
    folderLoadError.value = "無法載入資料夾列表";
  }
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
