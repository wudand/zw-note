import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';
import qs from 'qs';
import { ElMessage } from 'element-plus';

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 5 * 60 * 1000,
  paramsSerializer: {
    serialize: (params: any) => qs.stringify(params, { arrayFormat: 'repeat' }),
  },
});

service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers!['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

service.interceptors.response.use(
  (response: AxiosResponse<any>) => {
    if (response.data.code === 1) {
      throw response.data;
    }
    return response.data;
  },
  (error) => {
    const status = Number(error.response?.status) || 0;

    if (status === 401) {
      ElMessage.warning('登录状态已过期，请重新登录');
    } else if (status === 500) {
      ElMessage.error(error.response?.data?.message || '服务器异常，请稍后重试');
      return Promise.reject();
    }

    return Promise.reject(error.response?.data);
  },
);

export default service;
