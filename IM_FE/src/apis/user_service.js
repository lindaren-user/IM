import request from '@/utils/request';

export const user_service = {
  login: (form) => request.post('/', form),

  logout: () => request.post('/logout'),
};
