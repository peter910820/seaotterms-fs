import { ref, computed } from "vue";
import axios from "axios";
import type { FileResponseData, ResponseType } from "@/types/response";

const API_TIMEOUT_MS = 5000;
const FILE_API_BASE = `${import.meta.env.VITE_API_URL}/file`;

export function useFileList() {
  const currentPath = ref("");
  const currentDirectories = ref<string[]>([]);
  const currentFiles = ref<string[]>([]);
  const error = ref("");

  const pathLabel = computed(() => currentPath.value || "根目錄");

  const parentPath = computed(() => {
    const p = currentPath.value;
    if (!p) return null;
    const parts = p.split("/").filter(Boolean);
    if (parts.length <= 1) return "";
    return parts.slice(0, -1).join("/");
  });

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
    error.value = "";
    const result = await fetchFileList(path);
    if (result?.data) {
      currentDirectories.value = result.data.directories ?? [];
      currentFiles.value = result.data.files ?? [];
    } else {
      currentDirectories.value = [];
      currentFiles.value = [];
      error.value = "無法載入資料夾";
    }
  };

  const goIntoFolder = (name: string) => {
    const next = currentPath.value ? `${currentPath.value}/${name}` : name;
    loadPath(next);
  };

  return {
    currentPath,
    currentDirectories,
    currentFiles,
    error,
    pathLabel,
    parentPath,
    loadPath,
    goIntoFolder,
  };
}
