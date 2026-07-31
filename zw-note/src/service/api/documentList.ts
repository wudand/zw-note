import request from '../request';

/**
 * 获取文档列表
 * @returns 
 */
export const getDocumentList = () => {
    return request({
        url: '/documentList',
        method: 'get'
    });
}

/**
 * 创建文档
 * @returns 
 */
export const createDocument = (data: any) => {
    return request({
        url: '/document',
        method: 'post',
        data
    });
}

/**
 * 更新文档信息
 * @returns 
 */
export const updateDocument = (data: any) => {
    return request({
        url: '/document',
        method: 'put',
        data
    });
}

/**
 * 更新文档内容
 * @returns 
 */
export const updateDocumentContent = (data: any) => {
    return request({
        url: '/document/content',
        method: 'put',
        data
    });
}

/**
 * 获取文档目录
 * @returns 
 */
export const getDocumentOutline = (documentId: string) => {
    return request({
        url: `/document/outline/${documentId}`,
        method: 'get'
    });
}

/**
 * 获取文档目录对应的内容
 * @returns 
 */
export const getDocumentContent = (outlineId: string) => {
    return request({
        url: `/document/content/${outlineId}`,
        method: 'get'
    });
}

/**
 * 删除文档
 * @returns 
 */
export const deleteDocument = (id: string) => {
    return request({
        url: `/document/${id}`,
        method: 'delete'
    });
}