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
    // Se il server risponde 401, 404 o 403 (utente non esiste/bannato)
    if (error.response?.status === 401 || 
        error.response?.status === 404 || 
        error.response?.status === 403) {
      
      localStorage.clear();
      delete instance.defaults.headers.common['Authorization'];
      
      if (!window.location.href.includes('/login')) {
        window.location.href = '/#/login';
      }
    }
    return Promise.reject(error);
  }
);


const validateToken = async () => {
  
  const token = localStorage.getItem('token');
  const user = localStorage.getItem('user');
  
  
  if (!token || !user || token === '0' || token === 'null' || token === 'undefined') {
    localStorage.clear();
    return false;
  }
  
  try {
    const userData = JSON.parse(user);
    
    // Controlla che userData sia valido
    if (!userData || !userData.UserID || userData.UserID <= 0) {
      localStorage.clear();
      return false;
    }
    
    
    // Prova a fare una chiamata per verificare se l'utente esiste ancora
    const response = await instance.get(`/profiles/${userData.UserID}`);

    // Controlla se l'utente esiste davvero nei dati di risposta
    const profileData = response.data;
    
    // Se il backend restituisce un profilo vuoto o senza User, l'utente non esiste
    if (!profileData || 
        !profileData.User || 
        !profileData.User.UserID || 
        !profileData.User.Username ||
        profileData.User.UserID !== userData.UserID) {
      
      localStorage.clear();
      delete instance.defaults.headers.common['Authorization'];
      return false;
    }
    
    return true;
    
  } catch (error) {
    
    localStorage.clear();
    delete instance.defaults.headers.common['Authorization'];
    return false;
  }
};

const clearAuthData = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  delete instance.defaults.headers.common['Authorization'];
};

export { setAuth, validateToken, clearAuthData };
export default instance;
