import axios from "axios";

const instance = axios.create({
    baseURL: __API_URL__,
    timeout: 1000 * 5
});

const setAuth = () => {
  const token = localStorage.getItem('token');
  if (token && token !== '0' && token !== 'null' && token !== 'undefined') {
    instance.defaults.headers.common['Authorization'] = token;
  } else {
    delete instance.defaults.headers.common['Authorization'];
  }
}

instance.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      delete instance.defaults.headers.common['Authorization'];
      window.location.href = '/#/login';
    }
    return Promise.reject(error);
  }
);

export { setAuth };
export default instance;
