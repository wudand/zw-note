import request from '../request';

/**
 * 获取文档列表
 */
export const getDocumentList = () => {
  return request({
    url: '/notes/v1/documents',
    method: 'get',
  });
};

/**
 * 创建文档
 */
export const createDocument = (data: {
  title: string;
  author: string;
  description?: string;
}) => {
  return request({
    url: '/notes/v1/documents',
    method: 'post',
    data: {
      title: data.title,
      author: data.author,
      description: data.description ?? '',
    },
  });
};

/**
 * 获取单篇文档详情
 */
export const getDocumentById = (id: string) => {
  return request({
    url: `/notes/v1/documents/${id}`,
    method: 'get',
  });
};

/**
 * 更新文档信息（列表页「设置」弹窗）
 */
export const updateDocument = (data: {
  id: string;
  title?: string;
  author?: string;
  description?: string;
}) => {
  const { id, ...body } = data;
  return request({
    url: `/notes/v1/documents/${id}`,
    method: 'put',
    data: body,
  });
};

/**
 * 删除文档
 */
export const deleteDocument = (id: string) => {
  return request({
    url: `/notes/v1/documents/${id}`,
    method: 'delete',
  });
};

// ──────────────── 目录树 ────────────────

export type OutlineApiNode = {
  id: string;
  title: string;
  parent_id?: string | null;
  children?: OutlineApiNode[];
};

/**
 * 获取文档目录树
 */
export const getDocumentOutline = (documentId: string) => {
  return request({
    url: `/notes/v1/documents/${documentId}/outlines`,
    method: 'get',
  });
};

/**
 * 新建目录节点
 */
export const createOutline = (
  documentId: string,
  data: { title: string; parent_id?: number; sort_order?: number },
) => {
  return request({
    url: `/notes/v1/documents/${documentId}/outlines`,
    method: 'post',
    data,
  });
};

/**
 * 更新目录节点（标题 / 父节点 / 排序）
 */
export const updateOutline = (
  outlineId: string,
  data: {
    title?: string;
    parent_id?: number | null;
    clear_parent?: boolean;
    sort_order?: number;
  },
) => {
  return request({
    url: `/notes/v1/outlines/${outlineId}`,
    method: 'put',
    data,
  });
};

/**
 * 删除目录节点（子节点由后端级联删除）
 */
export const deleteOutline = (outlineId: string) => {
  return request({
    url: `/notes/v1/outlines/${outlineId}`,
    method: 'delete',
  });
};

/**
 * 批量重排目录节点
 */
export const reorderOutlines = (
  documentId: string,
  items: Array<{ id: number; parent_id?: number | null; sort_order: number }>,
) => {
  return request({
    url: `/notes/v1/documents/${documentId}/outlines/reorder`,
    method: 'put',
    data: { items },
  });
};

// ──────────────── 节点内容 ────────────────

/**
 * 获取目录节点 Markdown 内容
 */
export const getDocumentContent = (outlineId: string) => {
  return request({
    url: `/notes/v1/outlines/${outlineId}/content`,
    method: 'get',
  });
};

/**
 * 保存目录节点 Markdown 内容
 */
export const updateDocumentContent = (outlineId: string, content: string) => {
  return request({
    url: `/notes/v1/outlines/${outlineId}/content`,
    method: 'put',
    data: { content },
  });
};

// ──────────────── 文件上传 ────────────────

export type UploadImageResult = {
  /** 相对 upload.dir 的路径，如 images/20260811/xxx.png */
  path: string;
  /** 可直接用于 <img src> 的站点相对路径，如 /uploads/images/20260811/xxx.png */
  url: string;
};

/**
 * 上传图片（插入图片工具复用后台本地文件上传能力）
 */
export const uploadNoteImage = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  return request({
    url: '/notes/v1/upload/image',
    method: 'post',
    data: formData,
  });
};
