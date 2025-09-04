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

        <!-- Griglia posts -->
        <div v-if="posts.length > 0" class="posts-grid">
          <PostCard
            v-for="post in posts" 
            :key="post.PhotoID"
            :post="post"
            @click="openPostModal"
          />
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

      <!-- Modal per visualizzazione post -->
      <PostModal
        v-if="selectedPost"
        :postData="selectedPost"
        @close="closePostModal"
        @deleted="onPostDeleted"
        @postUpdated="onPostUpdated"
      />
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
      editError: null,

      // Post modal
      selectedPost: null
    }
  },
  computed: {
    profileUserId() {
      const userIdParam = this.$route.params.userID;
      
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
      return this.profileUserId && this.currentUserId && this.profileUserId === this.currentUserId;
    }
  },
  async created() {
    if (this.$route.params.userID === 'null' || this.$route.params.userID === 'undefined') {
      if (this.currentUserId) {
        this.$router.replace(`/profiles/${this.currentUserId}`);
        return;
      } else {
        this.$router.replace('/login');
        return;
      }
    }
    
    if (!this.profileUserId) {
      this.error = 'Invalid user ID in URL';
      this.loading = false;
      return;
    }
    
    await this.loadProfile();
    await this.loadPosts();

    // ✅ Ascolta per nuovi post uploadati
    window.addEventListener('postUploaded', this.handlePostUploaded);
  },

  beforeUnmount() {
    // ✅ Rimuovi listener
    window.removeEventListener('postUploaded', this.handlePostUploaded);
  },

  watch: {
    '$route.params.userID'(newUserID, oldUserID) {
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
    // ✅ Nuovo metodo per gestire post uploadati
    handlePostUploaded() {
      // Ricarica i post solo se siamo sul profilo corrente
      if (this.isOwnProfile) {
        this.loadPosts();
      }
    },

    async loadProfile() {
      if (!this.profileUserId) {
        this.error = 'Invalid user ID';
        this.loading = false;
        return;
      }
      
      this.loading = true;
      this.error = null;
      
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}`);
        this.profile = response.data;
      } catch (error) {
        console.error('Profile load error:', error);
        if (error.response?.status === 403) {
          this.error = 'This user has blocked you or their profile is private.';
        } else if (error.response?.status === 404) {
          this.error = 'User not found.';
        } else if (error.response?.status === 400) {
          this.error = 'Invalid user ID.';
        } else if (error.response?.status === 401) {
          this.$router.push('/login');
          return;
        } else {
          this.error = 'Failed to load profile. Please try again.';
        }
      } finally {
        this.loading = false;
      }
    },

    async loadPosts() {
      if (!this.profileUserId) return;
      
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/posts`);
        this.posts = response.data || [];
      } catch (error) {
        console.error('Posts load error:', error);
        if (error.response?.status === 401) {
          this.$router.push('/login');
          return;
        }
        this.posts = [];
      }
    },

    async loadFollowers() {
      this.loadingFollowers = true;
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/followers`);
        this.followers = response.data || [];
      } catch (error) {
        console.error('Followers load error:', error);
        this.followers = [];
      } finally {
        this.loadingFollowers = false;
      }
    },

    async loadFollowing() {
      this.loadingFollowing = true;
      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/followings`);
        this.following = response.data || [];
      } catch (error) {
        console.error('Following load error:', error);
        this.following = [];
      } finally {
        this.loadingFollowing = false;
      }
    },

    showFollowersList() {
      this.showFollowersModal = true;
      this.loadFollowers();
    },

    showFollowingList() {
      this.showFollowingModal = true;
      this.loadFollowing();
    },

    showEditProfileModal() {
      this.newUsername = this.profile.User.Username;
      this.editError = null;
      this.showEditModal = true;
    },

    async saveUsername() {
      this.savingUsername = true;
      this.editError = null;
      
      try {
        await this.$axios.put(`/profiles/${this.currentUserId}/username`, {
          username: this.newUsername
        });
        
        const user = JSON.parse(localStorage.getItem('user'));
        user.Username = this.newUsername;
        localStorage.setItem('user', JSON.stringify(user));
        
        await this.loadProfile();
        this.closeModal();
        
      } catch (error) {
        console.error('Username save error:', error);
        if (error.response?.status === 400) {
          this.editError = 'Username already taken or invalid.';
        } else {
          this.editError = 'Failed to update username. Please try again.';
        }
      } finally {
        this.savingUsername = false;
      }
    },

    handleUserBanned() {
      this.$router.push(`/profiles/${this.currentUserId}/feed`);
    },

    closeModal() {
      this.showFollowersModal = false;
      this.showFollowingModal = false;
      this.showEditModal = false;
      this.editError = null;
    },

    scrollToPosts() {
      const postsSection = document.querySelector('.posts-section');
      if (postsSection) {
        postsSection.scrollIntoView({ behavior: 'smooth' });
      }
    },

    openPostModal(post) {
      if (!post.UserID) {
        post.UserID = this.profileUserId;
      }
      
      if (!post.Username && this.profile?.User?.Username) {
        post.Username = this.profile.User.Username;
      }
      
      // ✅ CLONA il post per evitare modifiche dirette
      this.selectedPost = { ...post };
    },

    closePostModal() {
      this.selectedPost = null;
    },

    onPostDeleted() {
      this.loadPosts();
      this.closePostModal();
    },

    // ✅ NUOVO: Aggiorna il post nella lista quando cambia nel modal
    onPostUpdated(updatedData) {
      const postIndex = this.posts.findIndex(p => p.PhotoID === updatedData.PhotoID);
      if (postIndex !== -1) {
        // Aggiorna il post nella lista
        this.posts[postIndex] = {
          ...this.posts[postIndex],
          Nlike: updatedData.Nlike,
          IsLiked: updatedData.IsLiked,
          Ncomment: updatedData.Ncomment
        };
        
        // Forza il re-render del componente
        this.$forceUpdate();
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
  grid-template-columns: repeat(3, 1fr);
  gap: 15px;
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

/* ✅ RESPONSIVE: su mobile 2 colonne, su tablet 3 colonne */
@media (max-width: 768px) {
  .profile-view {
    padding: 10px;
  }

  .posts-grid {
    grid-template-columns: repeat(2, 1fr); /* 2 colonne su mobile */
    gap: 10px;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .posts-grid {
    grid-template-columns: repeat(3, 1fr); /* 3 colonne su tablet */
    gap: 12px;
  }
}

@media (min-width: 1025px) {
  .posts-grid {
    grid-template-columns: repeat(3, 1fr); /* 3 colonne su desktop */
    gap: 15px;
  }
}
</style>