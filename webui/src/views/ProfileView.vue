<template>
  <LoadingSpinner :loading="loading">
    <div class="profile-view">
      <!-- Mostra errore se l'utente è bannato o non trovato -->
      <ErrorMsg v-if="error" :msg="error" />
      
      <!-- Header del profilo -->
      <ProfileHeader 
        v-if="profile && !error"
        :profile="profile"
        :isOwnProfile="isOwnProfile"
        @profileUpdated="loadProfile"
        @userBanned="handleUserBanned"
        @editProfile="showEditProfileModal"
        @showFollowers="showFollowersList"
        @showFollowing="showFollowingList"
        @showPosts="scrollToPosts"
      />

      <!-- Sezione posts -->
      <div v-if="profile && !error" class="posts-section mt-4">
        <div class="d-flex justify-content-between align-items-center mb-3">
          <h3>Posts</h3>
          <button 
            v-if="isOwnProfile" 
            class="btn btn-primary"
            @click="createNewPost"
          >
            <i class="fas fa-plus"></i> New Post
          </button>
        </div>

        <!-- Griglia posts -->
        <div v-if="posts.length > 0" class="posts-grid">
          <div 
            v-for="post in posts" 
            :key="post.PhotoID"
            class="post-thumbnail"
            @click="openPost(post)"
          >
            <img 
              :src="getImageUrl(post.ImageData)" 
              :alt="post.Caption"
              class="post-image"
            >
            <div class="post-overlay">
              <div class="post-stats">
                <span class="stat-item">
                  <i class="fas fa-heart"></i> {{ post.Nlike }}
                </span>
                <span class="stat-item">
                  <i class="fas fa-comment"></i> {{ post.Ncomment }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Messaggio se non ci sono posts -->
        <div v-else class="no-posts text-center py-5">
          <i class="fas fa-camera fa-3x text-muted mb-3"></i>
          <h4 class="text-muted">No posts yet</h4>
          <p v-if="isOwnProfile" class="text-muted">
            Share your first photo to get started
          </p>
        </div>
      </div>

      <!-- Modal per followers -->
      <div v-if="showFollowersModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h5>Followers</h5>
            <button class="btn-close" @click="closeModal">&times;</button>
          </div>
          <div class="modal-body">
            <LoadingSpinner :loading="loadingFollowers">
              <div v-if="followers.length > 0" class="users-list">
                <div v-for="follower in followers" :key="follower.UserID" class="user-item">
                  <div class="user-info">
                    <div class="user-avatar">
                      <i class="fas fa-user"></i>
                    </div>
                    <RouterLink 
                      :to="`/profiles/${follower.UserID}`" 
                      class="username-link"
                      @click="closeModal"
                    >
                      {{ follower.Username }}
                    </RouterLink>
                  </div>
                </div>
              </div>
              <div v-else class="text-center text-muted">
                No followers yet
              </div>
            </LoadingSpinner>
          </div>
        </div>
      </div>

      <!-- Modal per following -->
      <div v-if="showFollowingModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h5>Following</h5>
            <button class="btn-close" @click="closeModal">&times;</button>
          </div>
          <div class="modal-body">
            <LoadingSpinner :loading="loadingFollowing">
              <div v-if="following.length > 0" class="users-list">
                <div v-for="user in following" :key="user.UserID" class="user-item">
                  <div class="user-info">
                    <div class="user-avatar">
                      <i class="fas fa-user"></i>
                    </div>
                    <RouterLink 
                      :to="`/profiles/${user.UserID}`" 
                      class="username-link"
                      @click="closeModal"
                    >
                      {{ user.Username }}
                    </RouterLink>
                  </div>
                </div>
              </div>
              <div v-else class="text-center text-muted">
                Not following anyone yet
              </div>
            </LoadingSpinner>
          </div>
        </div>
      </div>

      <!-- Modal per edit profile -->
      <div v-if="showEditModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h5>Edit Profile</h5>
            <button class="btn-close" @click="closeModal">&times;</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveUsername">
              <div class="mb-3">
                <label for="username" class="form-label">Username</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="username" 
                  v-model="newUsername"
                  required
                  minlength="3"
                  maxlength="13"
                  pattern="[a-zA-Z0-9_\-]+"
                >
              </div>
              <div class="d-flex justify-content-end">
                <button type="button" class="btn btn-secondary me-2" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary" :disabled="savingUsername">
                  {{ savingUsername ? 'Saving...' : 'Save' }}
                </button>
              </div>
            </form>
            <ErrorMsg v-if="editError" :msg="editError" />
          </div>
        </div>
      </div>
    </div>
  </LoadingSpinner>
</template>

<script>
export default {
  name: 'ProfileView',
  data() {
    return {
      profile: null,
      posts: [],
      loading: true,
      error: null,
      
      // Modal states
      showFollowersModal: false,
      showFollowingModal: false,
      showEditModal: false,
      
      // Followers/Following data
      followers: [],
      following: [],
      loadingFollowers: false,
      loadingFollowing: false,
      
      // Edit profile data
      newUsername: '',
      savingUsername: false,
      editError: null
    }
  },
  computed: {
    profileUserId() {
      const userIdParam = this.$route.params.userID;
      
      // Se il parametro è "null" (stringa), restituisci null
      if (userIdParam === 'null' || userIdParam === 'undefined') {
        return null;
      }
      
      const id = parseInt(userIdParam);
      return isNaN(id) ? null : id;
    },
    currentUserId() {
      try {
        const userStr = localStorage.getItem('user');
        if (!userStr || userStr === 'null' || userStr === 'undefined') return null;
        
        const user = JSON.parse(userStr);
        return user && user.UserID ? user.UserID : null;
      } catch (error) {
        console.error('Error parsing current user:', error);
        return null;
      }
    },
    isOwnProfile() {
      const result = this.profileUserId && this.currentUserId && this.profileUserId === this.currentUserId;
      return result;
    }
  },
  async created() {
    
    // Se userID è "null" (stringa), reindirizza al profilo dell'utente corrente
    if (this.$route.params.userID === 'null' || this.$route.params.userID === 'undefined') {
      if (this.currentUserId) {
        this.$router.replace(`/profiles/${this.currentUserId}`);
        return;
      } else {
        // Se non c'è neanche un utente corrente, vai al login
        this.$router.replace('/login');
        return;
      }
    }
    
    if (!this.profileUserId) {
      console.error('No valid profileUserId, cannot load profile');
      this.error = 'Invalid user ID in URL';
      this.loading = false;
      return;
    }
    
    await this.loadProfile();
    await this.loadPosts();
  },
  watch: {
    '$route.params.userID'(newUserID, oldUserID) {
      // Se il nuovo userID è "null", reindirizza
      if (newUserID === 'null' || newUserID === 'undefined') {
        if (this.currentUserId) {
          this.$router.replace(`/profiles/${this.currentUserId}`);
        }
        return;
      }
      
      this.loadProfile();
      this.loadPosts();
    }
  },
  methods: {
    async loadProfile() {
      if (!this.profileUserId) {
        console.error('Cannot load profile: profileUserId is', this.profileUserId);
        this.error = 'Invalid user ID';
        this.loading = false;
        return;
      }
      
      this.loading = true
      this.error = null
      
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}`)
        this.profile = response.data
      } catch (error) {
        console.error('Profile load error:', error)
        if (error.response?.status === 403) {
          this.error = 'This user has blocked you or their profile is private.'
        } else if (error.response?.status === 404) {
          this.error = 'User not found.'
        } else if (error.response?.status === 400) {
          this.error = 'Invalid user ID.'
        } else {
          this.error = 'Failed to load profile. Please try again.'
        }
      } finally {
        this.loading = false
      }
    },

    async loadPosts() {
      if (!this.profileUserId) return
      
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/posts`)
        this.posts = response.data || []
      } catch (error) {
        console.error('Posts load error:', error)
        this.posts = []
      }
    },

    async loadFollowers() {
      this.loadingFollowers = true
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/followers`)
        this.followers = response.data || []
      } catch (error) {
        console.error('Followers load error:', error)
        this.followers = []
      } finally {
        this.loadingFollowers = false
      }
    },

    async loadFollowing() {
      this.loadingFollowing = true
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/followings`)
        this.following = response.data || []
      } catch (error) {
        console.error('Following load error:', error)
        this.following = []
      } finally {
        this.loadingFollowing = false
      }
    },

    showFollowersList() {
      this.showFollowersModal = true
      this.loadFollowers()
    },

    showFollowingList() {
      this.showFollowingModal = true
      this.loadFollowing()
    },

    showEditProfileModal() {
      this.newUsername = this.profile.User.Username
      this.editError = null
      this.showEditModal = true
    },

    async saveUsername() {
      this.savingUsername = true
      this.editError = null
      
      try {
        await this.$axios.put(`/profiles/${this.currentUserId}/username`, {
          username: this.newUsername
        })
        
        // Aggiorna localStorage
        const user = JSON.parse(localStorage.getItem('user'))
        user.Username = this.newUsername
        localStorage.setItem('user', JSON.stringify(user))
        
        // Ricarica il profilo
        await this.loadProfile()
        this.closeModal()
        
      } catch (error) {
        console.error('Username save error:', error)
        if (error.response?.status === 400) {
          this.editError = 'Username already taken or invalid.'
        } else {
          this.editError = 'Failed to update username. Please try again.'
        }
      } finally {
        this.savingUsername = false
      }
    },

    handleUserBanned() {
      this.$router.push(`/profiles/${this.currentUserId}/feed`)
    },

    closeModal() {
      this.showFollowersModal = false
      this.showFollowingModal = false
      this.showEditModal = false
      this.editError = null
    },

    scrollToPosts() {
      const postsSection = document.querySelector('.posts-section')
      if (postsSection) {
        postsSection.scrollIntoView({ behavior: 'smooth' })
      }
    },

    createNewPost() {
      console.log('Create new post clicked')
    },

    openPost(post) {
      console.log('Open post:', post.PhotoID)
    },

    getImageUrl(imageData) {
      if (!imageData) return '/placeholder-image.jpg'
      
      if (typeof imageData === 'string' && imageData.startsWith('http')) {
        return imageData
      }
      
      try {
        const blob = new Blob([new Uint8Array(imageData)], { type: 'image/jpeg' })
        return URL.createObjectURL(blob)
      } catch (error) {
        console.error('Error creating image URL:', error)
        return '/placeholder-image.jpg'
      }
    }
  }
}
</script>

<style scoped>
.profile-view {
  max-width: 975px;
  margin: 0 auto;
  padding: 20px;
}

.posts-section {
  border-top: 1px solid #e1e8ed;
  padding-top: 2rem;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.post-thumbnail {
  position: relative;
  aspect-ratio: 1;
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.2s;
}

.post-thumbnail:hover {
  transform: scale(1.02);
}

.post-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.post-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.post-thumbnail:hover .post-overlay {
  opacity: 1;
}

.post-stats {
  color: white;
  text-align: center;
}

.stat-item {
  display: inline-block;
  margin: 0 10px;
  font-weight: 600;
}

.no-posts {
  color: #8e8e8e;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #e1e8ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-body {
  padding: 20px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
}

.users-list {
  max-height: 400px;
  overflow-y: auto;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.user-item:last-child {
  border-bottom: none;
}

.user-info {
  display: flex;
  align-items: center;
  width: 100%;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  color: #8e8e8e;
}

.username-link {
  color: #262626;
  text-decoration: none;
  font-weight: 500;
}

.username-link:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .profile-view {
    padding: 10px;
  }

  .posts-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 10px;
  }
}
</style>