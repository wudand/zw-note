import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { OutlineItem } from '@/components/DocumentOutline/index.vue';
import {
  getDocumentList,
  getDocumentOutline,
  getDocumentContent,
  createDocument,
  updateDocument,
  deleteDocument,
} from '@/service/api/documentList';

export interface DocumentItem {
  id: string;
  title: string;
  description?: string;
  author?: string;
}

export const useDocumentStore = defineStore('document', () => {
  // 文档列表
  const documents = ref<DocumentItem[]>([]);
  const documentsLoading = ref(false);

  // 当前打开的文档
  const currentDocumentId = ref<string>('');
  const currentOutline = ref<OutlineItem[]>([]);
  const currentContent = ref<string>('');
  const outlineLoading = ref(false);
  const contentLoading = ref(false);

  // ──────────────── 文档列表 ────────────────

  async function fetchDocuments() {
    documentsLoading.value = true;
    try {
      const res = await getDocumentList();
      documents.value = res.data.list ?? [];
    } finally {
      documentsLoading.value = false;
    }
  }

  async function addDocument(data: Omit<DocumentItem, 'id'>) {
    const res = await createDocument(data);
    documents.value.push(res.data);
    return res.data as DocumentItem;
  }

  async function editDocument(data: DocumentItem) {
    await updateDocument(data);
    const idx = documents.value.findIndex((d) => d.id === data.id);
    if (idx !== -1) documents.value[idx] = { ...documents.value[idx], ...data };
  }

  async function removeDocument(id: string) {
    await deleteDocument(id);
    documents.value = documents.value.filter((d) => d.id !== id);
  }

  // ──────────────── 当前文档目录 & 内容 ────────────────

  async function fetchOutline(documentId: string) {
    currentDocumentId.value = documentId;
    outlineLoading.value = true;
    try {
      const res = await getDocumentOutline(documentId);
      currentOutline.value = res.data ?? [];
    } finally {
      outlineLoading.value = false;
    }
  }

  async function fetchContent(outlineId: string) {
    contentLoading.value = true;
    try {
      const res = await getDocumentContent(outlineId);
      currentContent.value = res.data ?? '';
    } finally {
      contentLoading.value = false;
    }
  }

  function resetCurrentDocument() {
    currentDocumentId.value = '';
    currentOutline.value = [];
    currentContent.value = '';
  }

  return {
    // state
    documents,
    documentsLoading,
    currentDocumentId,
    currentOutline,
    currentContent,
    outlineLoading,
    contentLoading,
    // actions
    fetchDocuments,
    addDocument,
    editDocument,
    removeDocument,
    fetchOutline,
    fetchContent,
    resetCurrentDocument,
  };
});
