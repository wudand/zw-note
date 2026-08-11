import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { OutlineItem } from '@/components/DocumentOutline/index.vue';
import {
  getDocumentList,
  getDocumentById,
  getDocumentOutline,
  getDocumentContent,
  updateDocumentContent,
  createDocument,
  updateDocument,
  deleteDocument,
  createOutline,
  updateOutline,
  deleteOutline,
  reorderOutlines,
  type OutlineApiNode,
} from '@/service/api/documentList';

export interface DocumentItem {
  id: string;
  title: string;
  description?: string;
  author?: string;
  created_at?: string;
  updated_at?: string;
}

/** 后端 OutlineNode → 前端 OutlineItem（parent_id → parentId） */
function mapOutlineNode(node: OutlineApiNode): OutlineItem {
  return {
    id: String(node.id),
    title: node.title,
    parentId: node.parent_id != null && node.parent_id !== '' ? String(node.parent_id) : undefined,
    children: (node.children ?? []).map(mapOutlineNode),
  };
}

function toParentIdNumber(parentId?: string | number | null): number | undefined {
  if (parentId === undefined || parentId === null || parentId === '') return undefined;
  const n = Number(parentId);
  return Number.isFinite(n) ? n : undefined;
}

export const useDocumentStore = defineStore('document', () => {
  const documents = ref<DocumentItem[]>([]);
  const documentsLoading = ref(false);

  const currentDocumentId = ref<string>('');
  const currentDocument = ref<DocumentItem | null>(null);
  const currentOutline = ref<OutlineItem[]>([]);
  const currentContent = ref<string>('');
  const currentOutlineId = ref<string>('');
  const outlineLoading = ref(false);
  const contentLoading = ref(false);
  const contentSaving = ref(false);
  const documentLoading = ref(false);

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
    const doc = res.data as DocumentItem;
    documents.value.push(doc);
    return doc;
  }

  async function editDocument(data: DocumentItem) {
    const res = await updateDocument(data);
    const updated = (res.data ?? data) as DocumentItem;
    const idx = documents.value.findIndex((d) => d.id === data.id);
    if (idx !== -1) {
      documents.value[idx] = { ...documents.value[idx], ...updated };
    }
  }

  async function removeDocument(id: string) {
    await deleteDocument(id);
    documents.value = documents.value.filter((d) => d.id !== id);
  }

  async function fetchDocumentById(id: string) {
    documentLoading.value = true;
    currentDocumentId.value = id;
    try {
      const res = await getDocumentById(id);
      currentDocument.value = res.data as DocumentItem;
      return currentDocument.value;
    } finally {
      documentLoading.value = false;
    }
  }

  // ──────────────── 目录树 ────────────────

  async function fetchOutline(documentId: string) {
    currentDocumentId.value = documentId;
    outlineLoading.value = true;
    try {
      const res = await getDocumentOutline(documentId);
      const tree = (res.data ?? []) as OutlineApiNode[];
      currentOutline.value = tree.map(mapOutlineNode);
    } finally {
      outlineLoading.value = false;
    }
  }

  async function createOutlineNode(payload: {
    title: string;
    parentId?: string | number | null;
  }) {
    const documentId = currentDocumentId.value;
    if (!documentId) throw new Error('未打开文档');

    const body: { title: string; parent_id?: number } = { title: payload.title };
    const parentNum = toParentIdNumber(payload.parentId);
    if (parentNum !== undefined) body.parent_id = parentNum;

    await createOutline(documentId, body);
    await fetchOutline(documentId);
  }

  async function updateOutlineNode(payload: {
    id: string;
    title: string;
    parentId?: string | number | null;
  }) {
    const documentId = currentDocumentId.value;
    if (!documentId) throw new Error('未打开文档');

    const body: {
      title: string;
      parent_id?: number;
      clear_parent?: boolean;
    } = { title: payload.title };

    const parentNum = toParentIdNumber(payload.parentId);
    if (parentNum === undefined) {
      // 清空上级 → 挪到根级
      body.clear_parent = true;
    } else if (String(parentNum) !== String(payload.id)) {
      body.parent_id = parentNum;
    }

    await updateOutline(payload.id, body);
    await fetchOutline(documentId);
  }

  async function removeOutlineNode(outlineId: string) {
    const documentId = currentDocumentId.value;
    if (!documentId) throw new Error('未打开文档');

    await deleteOutline(outlineId);
    if (currentOutlineId.value === outlineId) {
      currentOutlineId.value = '';
      currentContent.value = '';
    }
    await fetchOutline(documentId);
  }

  /** 批量重排目录（整棵树的 id / parent_id / sort_order） */
  async function reorderOutlineNodes(
    items: Array<{ id: number; parent_id?: number | null; sort_order: number }>,
  ) {
    const documentId = currentDocumentId.value;
    if (!documentId) throw new Error('未打开文档');
    await reorderOutlines(documentId, items);
    await fetchOutline(documentId);
  }

  // ──────────────── 节点内容 ────────────────

  async function fetchContent(outlineId: string) {
    contentLoading.value = true;
    currentOutlineId.value = outlineId;
    try {
      const res = await getDocumentContent(outlineId);
      currentContent.value = res.data?.content ?? '';
    } finally {
      contentLoading.value = false;
    }
  }

  async function saveContent(outlineId: string, content: string) {
    if (!outlineId) throw new Error('未选中目录节点');
    contentSaving.value = true;
    try {
      const res = await updateDocumentContent(outlineId, content);
      currentContent.value = res.data?.content ?? content;
      currentOutlineId.value = outlineId;
      return currentContent.value;
    } finally {
      contentSaving.value = false;
    }
  }

  function resetCurrentDocument() {
    currentDocumentId.value = '';
    currentDocument.value = null;
    currentOutline.value = [];
    currentContent.value = '';
    currentOutlineId.value = '';
  }

  return {
    documents,
    documentsLoading,
    currentDocumentId,
    currentDocument,
    currentOutline,
    currentContent,
    currentOutlineId,
    outlineLoading,
    contentLoading,
    contentSaving,
    documentLoading,
    fetchDocuments,
    fetchDocumentById,
    addDocument,
    editDocument,
    removeDocument,
    fetchOutline,
    createOutlineNode,
    updateOutlineNode,
    removeOutlineNode,
    reorderOutlineNodes,
    fetchContent,
    saveContent,
    resetCurrentDocument,
  };
});
