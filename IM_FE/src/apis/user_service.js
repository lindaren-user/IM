import request from '@/utils/request';

export const user_service = {
  login: (form) => request.post('/user', form),

  logout: () => request.post('/user/logout'),

  search: (type, keyword, signal) =>
    request.get(`/user/search?type=${type}&keyword=${keyword}`, {
      signal,
    }),
};
