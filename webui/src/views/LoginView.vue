<template>
  <LoadingSpinner :loading="loading">
    <div class="container mt-5">
      <div class="row justify-content-center">
        <div class="col-md-6">
          <div class="card">
            <div class="card-body">
              <h2 class="card-title text-center">WASAPhoto</h2>
              <form @submit.prevent="doLogin">
                <div class="mb-3">
                  <label for="username" class="form-label">Username</label>
                  <input 
                    type="text" 
                    class="form-control" 
                    id="username" 
                    v-model="username" 
                    required
                    minlength="3"
                    maxlength="13"
                    pattern="[a-zA-Z0-9_\-]+"
                  >
                </div>
                <button type="submit" class="btn btn-primary w-100">
                  Login
                </button>
              </form>
              <ErrorMsg v-if="error" :msg="error" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </LoadingSpinner>
</template>

<script>
import { setAuth } from '@/services/axios.js'

export default {
  name: 'LoginView',
  data() {
    return {
      username: '',
      loading: false,
      error: null
    }
  },
  created() {
    // Pulisci localStorage e assicurati che non ci siano header
    localStorage.clear();
    
    // Rimuovi qualsiasi header Authorization esistente
    delete this.$axios.defaults.headers.common['Authorization'];
  },
  methods: {
    async doLogin() {
      this.loading = true
      this.error = null
      
      try {
        const response = await this.$axios.post('/login', {
          username: this.username
        })
        
        // Salva i dati SOLO se la risposta è valida
        if (response.data.Token && response.data.User && response.data.User.UserID > 0) {
          localStorage.setItem('token', response.data.Token.toString())
          localStorage.setItem('user', JSON.stringify(response.data.User))
          
          // ORA imposta l'header Authorization
          setAuth();
          
          // IMPORTANTE: Forza il reload per pulire tutto lo stato precedente
          const userId = response.data.User.UserID
          window.location.href = `/#/profiles/${userId}/feed`;
          
        } else {
          throw new Error('Invalid response from server');
        }
        
      } catch (error) {
        console.error('Login error:', error)
        this.error = 'Login failed. Please try again.'
        
        // Pulisci localStorage in caso di errore
        localStorage.clear();
        delete this.$axios.defaults.headers.common['Authorization'];
      } finally {
        this.loading = false
      }
    }
  }
}
</script>