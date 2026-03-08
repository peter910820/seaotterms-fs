<template>
  <div class="md-writer-page">
    <v-card class="md-writer-card" variant="outlined">
      <v-card-title class="d-flex align-center text-h5 py-4">
        <v-icon start color="primary">mdi-text-box-edit-outline</v-icon>
        Markdown 預覽器
      </v-card-title>
      <v-divider />
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-textarea
              id="markdown-input"
              v-model="markdownText"
              label="輸入 Markdown"
              variant="outlined"
              auto-grow
              rows="16"
              class="font-monospace"
              hide-details
            />
          </v-col>
          <v-col cols="12" md="6">
            <div id="markdown-display" class="markdown-preview pa-4 rounded" v-html="renderedHtml" />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import markdownit from "markdown-it";
import hljs from "highlight.js";

const md = markdownit({
  highlight(str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return (
          '<pre><code class="hljs">' +
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
          "</code></pre>"
        );
      } catch {
        // ignore
      }
    }
    return '<pre><code class="hljs">' + markdownit().utils.escapeHtml(str) + "</code></pre>";
  },
});

const markdownText = ref("");
const renderedHtml = computed(() => md.render(markdownText.value));

onMounted(() => {
  hljs.highlightAll();
});

watch(renderedHtml, () => {
  setTimeout(() => hljs.highlightAll(), 0);
});
</script>

<style scoped>
.md-writer-page {
  max-width: 1200px;
  margin: 0 auto;
}
.md-writer-card {
  border-radius: 12px;
}
.markdown-preview {
  min-height: 320px;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.markdown-preview :deep(pre) {
  padding: 0.75rem;
  border-radius: 8px;
  overflow-x: auto;
}
.markdown-preview :deep(code) {
  font-family: ui-monospace, monospace;
}
</style>
