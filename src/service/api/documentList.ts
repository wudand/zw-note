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
 * 更新文档
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
 * 删除文档
 * @returns 
 */
export const deleteDocument = (id: string) => {
    return request({
        url: `/document/${id}`,
        method: 'delete'
    });
}