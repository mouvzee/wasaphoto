<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
import { setAuth } from '@/services/axios.js'

export default {
  data() {
    return {
      // Stato reattivo per l'utente corrente
      reactiveUserId: null,
      storageCheckInterval: null,
      showUploadModal: false
    }
  },
  computed: {
    isLoginPage() {
      return this.$route.path === '/login'
    },
    currentUserId() {
      // Ora usa lo stato reattivo invece del localStorage direttamente
      return this.reactiveUserId;
    }
  },
  created() {
    // ✅ INIZIALIZZA il token all'avvio
    this.initializeAuth();
    this.updateCurrentUser();
    
    // Controlla periodicamente per cambiamenti nel localStorage
    this.storageCheckInterval = setInterval(() => {
      this.updateCurrentUser();
    }, 500); // Controlla ogni 500ms
  },
  beforeUnmount() {
    // Pulisci l'interval quando il componente viene distrutto
    if (this.storageCheckInterval) {
      clearInterval(this.storageCheckInterval);
    }
  },
  watch: {
    // Aggiorna anche quando cambia la route
    '$route'() {
      this.updateCurrentUser();
    }
  },
  methods: {
    // ✅ NUOVO: Inizializza l'autenticazione
    initializeAuth() {
      const token = localStorage.getItem('token');
      if (token && token !== '0' && token !== 'null' && token !== 'undefined') {
        setAuth(token);
      }
    },

    updateCurrentUser() {
      try {
        const userStr = localStorage.getItem('user');
        
        if (!userStr || userStr === 'null' || userStr === 'undefined') {
          this.reactiveUserId = null;
          return;
        }
        
        const user = JSON.parse(userStr);
        
        if (!user || !user.UserID || user.UserID === 0) {
          this.reactiveUserId = null;
          return;
        }
        
        // Aggiorna SOLO se è cambiato
        if (this.reactiveUserId !== user.UserID) {
          this.reactiveUserId = user.UserID;
        }
      } catch (error) {
        console.error('App.vue - Error parsing user:', error);
        this.reactiveUserId = null;
      }
    },
    
    logout() {
      localStorage.clear();
      setAuth(null); 
      // Forza l'aggiornamento dell'utente corrente
      this.reactiveUserId = null;
      
      window.location.href = '/#/login';
    },

    openUploadModal() {
      this.showUploadModal = true;
    },

    closeUploadModal() {
      this.showUploadModal = false;
    },

    onPostUploaded() {
      this.showUploadModal = false;
      this.$nextTick(() => {
        // Forza il reload dei post se siamo sulla pagina del profilo
        window.dispatchEvent(new CustomEvent('postUploaded'));
      });
    }
  }
}
</script>

<template>
  <div class="container-fluid">
    <!-- Sidebar e layout normale solo se NON siamo nella pagina di login -->
    <div v-if="!isLoginPage" class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <!-- Header con logo e logout -->
          <div class="sidebar-header px-3 mb-4">
            <h4 class="sidebar-brand mb-3">WASAPhoto</h4>
            <button class="btn btn-outline-danger btn-sm w-100" @click="logout">
              <i class="fas fa-sign-out-alt me-2"></i>
              Sign out
            </button>
          </div>

          <!-- Separatore -->
          <hr class="mx-3 text-muted">
          
          <!-- Pulsante Upload Post -->
          <div class="px-3 mb-3">
            <button 
              v-if="currentUserId"
              class="btn btn-primary w-100 upload-btn"
              @click="openUploadModal"
            >
              <i class="fas fa-plus-circle me-2"></i>
              <span>New Post</span>
            </button>
          </div>
          
          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>General</span>
          </h6>
          
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink 
                v-if="currentUserId" 
                :to="`/profiles/${currentUserId}/feed`" 
                class="nav-link"
                :key="`home-${currentUserId}`"
              >
                <i class="fas fa-home feather"></i>
                <span>Home</span>
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink 
                v-if="currentUserId" 
                :to="`/profiles/${currentUserId}`" 
                class="nav-link"
                :key="`profile-${currentUserId}`"
              >
                <i class="fas fa-user feather"></i>
                <span>My Profile</span>
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink to="/search" class="nav-link">
                <i class="fas fa-search feather"></i>
                <span>Search</span>
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>

      <!-- Contenuto principale ora occupa tutto lo spazio disponibile -->
      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4 main-content">
        <RouterView :key="$route.fullPath" />
      </main>
    </div>
    
    <!-- Layout full-page per la pagina di login -->
    <div v-else>
      <RouterView />
    </div>

    <!-- Upload Modal -->
    <UploadModal 
      v-if="showUploadModal"
      @close="closeUploadModal"
      @uploaded="onPostUploaded"
    />
  </div>
</template>

<style>
/* Maggiore spaziatura tra icona e testo */
.feather {
  width: 16px;
  height: 16px;
  margin-right: 12px; /* Aumentato da 8px a 12px */
}

/* Sidebar senza navbar */
.sidebar {
  position: fixed;
  top: 0; /* Ora parte da 0 invece che da 48px */
  bottom: 0;
  left: 0;
  z-index: 100;
  padding: 0; /* Rimosso padding top */
  box-shadow: inset -1px 0 0 rgba(0, 0, 0, .1);
  background-color: #f8f9fa !important;
}

.sidebar-sticky {
  position: relative;
  top: 0;
  height: 100vh; /* Ora occupa tutta l'altezza */
  padding-top: 1rem;
  overflow-x: hidden;
  overflow-y: auto;
}

/* Stili per l'header della sidebar */
.sidebar-header {
  border-bottom: 1px solid #e9ecef;
  padding-bottom: 1rem;
}

.sidebar-brand {
  color: #2470dc;
  font-weight: 600;
  text-align: center;
  margin-bottom: 1rem;
}

/* Pulsante Upload Post */
.upload-btn {
  font-weight: 600;
  border-radius: 25px;
  padding: 0.75rem 1rem;
  background: linear-gradient(45deg, #405de6, #5851db, #833ab4, #c13584, #e1306c, #fd1d1d);
  border: none;
  color: white;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.upload-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
  background: linear-gradient(45deg, #405de6, #5851db, #833ab4, #c13584, #e1306c, #fd1d1d);
}

/* Migliori stili per i nav links */
.sidebar .nav-link {
  font-weight: 500;
  color: #333;
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem; /* Aumentato padding per più spazio */
  border-radius: 0.375rem;
  margin: 0.25rem 0.5rem; /* Margini per distanziare i link */
  transition: all 0.2s ease;
}

.sidebar .nav-link:hover {
  color: #2470dc;
  background-color: rgba(36, 112, 220, 0.1);
}

.sidebar .nav-link.router-link-active {
  color: #2470dc;
  background-color: rgba(36, 112, 220, 0.15);
  font-weight: 600;
}

/* Contenuto principale ora parte da sinistra */
.main-content {
  padding-top: 2rem;
  min-height: 100vh;
}

@media (max-width: 767.98px) {
  .sidebar {
    position: relative;
    height: auto;
  }
  
  .main-content {
    padding-top: 1rem;
  }
  
  .sidebar-header {
    padding: 1rem;
  }
}
</style>
