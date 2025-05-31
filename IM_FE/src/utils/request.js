import axios from 'axios';

const request = axios.create({
  baseURL: '/api',
  timeout: 10000,
  header: {
    'Content-Type': 'application/json', // 默认情况下所有通过这个axios实例发出的请求
  },
  withCredentials: true, // 允许跨域请求时携带凭据（如 Cookies）
});

export default request;
