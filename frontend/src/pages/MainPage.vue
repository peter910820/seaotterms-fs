<template>
  <div v-if="!isLoggedIn" class="main-page main-page--empty"></div>
  <div v-else class="main-page">
    <v-row>
      <v-col cols="12" md="4" lg="3">
        <v-card class="fs-card folder-picker-card" variant="outlined">
          <v-card-title class="d-flex align-center text-subtitle-1 py-3">
            <v-icon start color="primary">mdi-folder-open</v-icon>
            選擇上傳位置
          </v-card-title>
          <v-divider />
          <v-card-text class="py-2">
            <div class="text-caption text-medium-emphasis mb-2">目前位置</div>
            <div class="current-path mb-3">{{ currentPathLabel }}</div>
            <v-list density="compact" class="py-0">
              <v-list-item v-if="parentPath !== null" rounded="lg" class="mb-1" @click="goToPath(parentPath)">
                <template #prepend>
                  <v-icon size="small">mdi-folder-arrow-up</v-icon>
                </template>
                <v-list-item-title class="text-body2">..</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-for="(name, index) in currentDirectories"
                :key="'dir-' + index"
                rounded="lg"
                class="mb-1"
                @click="goIntoFolder(name)"
              >
                <template #prepend>
                  <v-icon size="small">mdi-folder-outline</v-icon>
                </template>
                <v-list-item-title class="text-body2">{{ name }}</v-list-item-title>
              </v-list-item>
            </v-list>
            <v-card-text v-if="folderError" class="text-caption text-error text-center py-2">
              {{ folderError }}
            </v-card-text>
            <v-card-text
              v-else-if="currentDirectories.length === 0 && parentPath === null"
              class="text-caption text-medium-emphasis text-center py-2"
            >
              尚無子資料夾
            </v-card-text>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 右：上傳表單 -->
      <v-col cols="12" md="8" lg="9">
        <v-card class="fs-card upload-card" variant="outlined">
          <v-card-title class="d-flex align-center text-h5 py-4">
            <v-icon start color="primary">mdi-upload</v-icon>
            資源伺服器目錄
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="8">
                <v-file-input
                  v-model="file"
                  label="選擇檔案上傳"
                  placeholder="只允許上傳單一檔案"
                  prepend-icon=""
                  prepend-inner-icon="mdi-upload-file"
                  variant="outlined"
                  density="comfortable"
                  hide-details="auto"
                  clearable
                  show-size
                  @update:model-value="onFileSelected"
                />
              </v-col>
              <v-col cols="12" md="8">
                <v-text-field
                  v-model="uploadFilename"
                  label="上傳後的檔名（選填，留空則使用原檔名）"
                  placeholder="例如：報告.pdf"
                  variant="outlined"
                  density="comfortable"
                  hide-details="auto"
                  clearable
                />
              </v-col>
              <v-col cols="12">
                <v-btn color="primary" size="large" :disabled="!hasFile" :loading="uploading" @click="upload">
                  <v-icon start>mdi-send</v-icon>
                  上傳至 {{ currentPathLabel }}
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import axios from "axios";

import { useAuthStore } from "@/store/auth";
import type { FileResponseData, ResponseType } from "@/types/response";

const router = useRouter();
const { hasSession } = storeToRefs(useAuthStore());
const file = ref<File[] | null>(null);
const uploadFilename = ref("");
const uploading = ref(false);

function onFileSelected(files: File | File[] | null) {
  if (files) {
    const single = Array.isArray(files) ? files[0] : files;
    uploadFilename.value = single?.name ?? "";
  } else {
    uploadFilename.value = "";
  }
}

const hasFile = computed(() => {
  const v = file.value;
  if (Array.isArray(v)) return (v?.length ?? 0) > 0;
  return !!v;
});

// 依 store 登入狀態（與 router 同步），使登入成功後頁面即時更新並觸發載入
const isLoggedIn = hasSession;

const API_TIMEOUT_MS = 5000;

const currentPath = ref("");
const currentDirectories = ref<string[]>([]);
const folderError = ref("");

const currentPathLabel = computed(() => (currentPath.value ? currentPath.value : "根目錄"));

const parentPath = computed(() => {
  const p = currentPath.value;
  if (!p) return null;
  const parts = p.split("/").filter(Boolean);
  if (parts.length <= 1) return "";
  return parts.slice(0, -1).join("/");
});

const fetchFolder = async (path: string): Promise<ResponseType<FileResponseData> | null> => {
  const url = path ? `${import.meta.env.VITE_API_URL}/file/${path}` : `${import.meta.env.VITE_API_URL}/file`;
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
  folderError.value = "";
  const result = await fetchFolder(path);
  if (result?.data) {
    currentDirectories.value = result.data.directories ?? [];
  } else {
    folderError.value = "無法載入資料夾";
    currentDirectories.value = [];
  }
};

const goToPath = (path: string) => {
  loadPath(path);
};

const goIntoFolder = (name: string) => {
  const next = currentPath.value ? `${currentPath.value}/${name}` : name;
  loadPath(next);
};

// 上傳時送給後端的路徑：根目錄送 "."，子資料夾送相對路徑
const uploadDirectory = computed(() => (currentPath.value ? currentPath.value : "."));

onMounted(async () => {
  if (!isLoggedIn.value) return;
  await loadPath("");
});

watch(isLoggedIn, (loggedIn) => {
  if (loggedIn) loadPath("");
});

const upload = async () => {
  const f = Array.isArray(file.value) ? file.value?.[0] : file.value;
  if (!f) {
    alert("請選擇檔案");
    return;
  }
  const formData = new FormData();
  formData.append("file", f);
  formData.append("directory", uploadDirectory.value);
  const nameOverride = uploadFilename.value?.trim();
  if (nameOverride) formData.append("filename", nameOverride);
  uploading.value = true;
  try {
    const response = await axios.post(`${import.meta.env.VITE_API_URL}/upload`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    const msg = response?.data?.message ?? "檔案上傳成功！";
    alert(msg);
    file.value = null;
    uploadFilename.value = "";
  } catch (err: unknown) {
    if (axios.isAxiosError(err) && err.response?.data?.message) {
      alert(err.response.data.message);
    } else {
      router.push("/error");
    }
  } finally {
    uploading.value = false;
  }
};
</script>

<style scoped>
.main-page {
  margin: 0 auto;
}
.main-page--empty {
  min-height: 60vh;
}
.upload-card,
.folder-picker-card {
  border-radius: 12px;
}
.current-path {
  font-size: 0.9rem;
  word-break: break-all;
  color: rgb(var(--v-theme-on-surface));
}
</style>
