import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig} from 'axios';
import qs from 'qs';
import { ElMessage } from 'element-plus';

/**
 * 创建并配置一个 Axios 实例对象
 */
const service: AxiosInstance = axios.create({
	baseURL: import.meta.env.VITE_API_URL,
	timeout: 5*60*1000, // 全局超时时间
    paramsSerializer: {
        serialize: (params: any) => {
            return qs.stringify(params, {arrayFormat: 'repeat'});
        }
    }
});

/**
 * Axios请求拦截器，对请求进行处理
 * 1. 序列化get请求参数
 * 2. 统一增加Authorization和TENANT-ID请求头
 * 3. 自动适配单体、微服务架构不同的URL
 * @param config AxiosRequestConfig对象，包含请求配置信息
 */
service.interceptors.request.use(

	(config: InternalAxiosRequestConfig) => {
		// 统一增加Authorization请求头, skipToken 跳过增加token
		const token = localStorage.getItem('token');
		if (token && !config.headers?.skipToken) {
			config.headers!['Authorization'] = `Bearer ${token}`;
		}

		// 处理完毕，返回config对象
		return config;
	},
	(error) => {
		// 对请求错误进行处理
		return Promise.reject(error);
	}
);

/**
 * 响应拦截器处理函数
 * @param response 响应结果
 * @returns 如果响应成功，则返回响应的data属性；否则，抛出错误或者执行其他操作
 */
const handleResponse = (response: AxiosResponse<any>) => {
	if (response.data.code === 1) {
		throw response.data;
	}

	// // 针对密文返回解密
	// if (response.data.encryption) {
	// 	const originData = JSON.parse(other.decryption(response.data.encryption, import.meta.env.VITE_PWD_ENC_KEY));
	// 	response.data = originData;
	// 	return response.data;
	// }

	return response.data;
};

/**
 * 添加 Axios 的响应拦截器，用于全局响应结果处理
 */
service.interceptors.response.use(handleResponse, (error) => {
	const status = Number(error.response.status) || 200;

	if (status === 424) {
		ElMessage.warning('令牌状态已过期，请点击重新登录');
	}

	if (status === 500) {
		ElMessage.error(error.response.data.message || '系统异常，请联系管理员');
		return Promise.reject();
	}
	return Promise.reject(error.response.data);
});

// // 常用header
// export enum CommonHeaderEnum {
// 	'TENANT_ID' = 'TENANT-ID',
// 	'ENC_FLAG' = 'Enc-Flag',
// 	'AUTHORIZATION' = 'Authorization',
// }

// 导出 axios 实例
export default service;