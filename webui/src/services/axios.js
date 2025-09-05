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

// Controlla token prima di ogni chiamata
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token && token !== '0' && token !== 'null' && token !== 'undefined') {
      config.headers.Authorization = token;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.clear();
      delete instance.defaults.headers.common['Authorization'];
      window.location.href = '/#/login';
    }
    return Promise.reject(error);
  }
);

// Funzione per validare token
const validateToken = async () => {
  const token = localStorage.getItem('token');
  const user = localStorage.getItem('user');
  
  if (!token || !user || token === '0' || token === 'null' || token === 'undefined') {
    localStorage.clear();
    return false;
  }
  
  try {
    // Prova a fare una chiamata per verificare se il token è valido
    const userData = JSON.parse(user);
    await instance.get(`/profiles/${userData.UserID}`);
    return true;
  } catch (error) {
    // Token non valido, pulisci tutto
    localStorage.clear();
    delete instance.defaults.headers.common['Authorization'];
    return false;
  }
};

export { setAuth, validateToken };
export default instance;
