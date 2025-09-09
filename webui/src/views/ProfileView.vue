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
        @showBannedUsers="showBannedUsersList"
      />

      <!-- Sezione posts -->
      <div v-if="profile && !error" class="posts-section mt-4" ref="postsSection">

        <!-- Griglia posts -->
        <div class="posts-grid">
          <PostCard
            v-for="post in posts" 
            :key="post.PhotoID"
            :post="post"
            @click="openPostModal"
          />
        </div>
        
        <!-- Empty state per i post -->
        <div v-if="posts.length === 0" class="no-posts">
          <p class="text-muted">No posts yet</p>
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
                    <div class="user-avatar" :style="getUserAvatarStyle(follower.Username)">
                      <span class="avatar-initials">{{ getUserInitials(follower.Username) }}</span>
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
                    <div class="user-avatar" :style="getUserAvatarStyle(user.Username)">
                      <span class="avatar-initials">{{ getUserInitials(user.Username) }}</span>
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

      <div v-if="showBannedModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h5>Banned Users</h5>
            <button class="btn-close" @click="closeModal">&times;</button>
          </div>
          <div class="modal-body">
            <LoadingSpinner :loading="loadingBanned">
              <div v-if="bannedUsers.length > 0" class="users-list">
                <div v-for="user in bannedUsers" :key="user.UserID" class="user-item">
                  <div class="user-info">
                    <div class="user-avatar" :style="getUserAvatarStyle(user.Username)">
                      <span class="avatar-initials">{{ getUserInitials(user.Username) }}</span>
                    </div>
                    <span class="username">{{ user.Username }}</span>
                  </div>
                  <button 
                    class="btn btn-sm btn-outline-primary"
                    @click="unbanUser(user)"
                    :disabled="unbanning === user.UserID"
                  >
                    {{ unbanning === user.UserID ? 'Unbanning...' : 'Unban' }}
                  </button>
                </div>
              </div>
              <div v-else class="text-center text-muted">
                No banned users
              </div>
            </LoadingSpinner>
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
      showBannedModal: false,
      
      // Followers/Following data
      followers: [],
      following: [],
      loadingFollowers: false,
      loadingFollowing: false,
      
      
      bannedUsers: [],
      loadingBanned: false,
      unbanning: null,
      
      // Edit profile data
      newUsername: '',
      savingUsername: false,
      editError: null,

      // Post modal
      selectedPost: null,
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
    await this.loadProfile();
    await this.loadPosts();

    window.addEventListener('postUploaded', this.onPostUploaded);
  },

  beforeUnmount() {
    window.removeEventListener('postUploaded', this.onPostUploaded);
  },

  methods: {
    async onPostUploaded() {
      // Ricarica sia il profilo (per il counter) che i post
      await Promise.all([
        this.loadProfile(),
        this.loadPosts()
      ]);
    },

    async loadProfile() {
      this.loading = true;
      this.error = null;

      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}`);
        
        if (response.data) {
          this.profile = response.data;
        } else {
          throw new Error('No profile data received');
        }
      } catch (error) {
        console.error('Profile load error:', error);
        this.error = 'Failed to load profile. Please try again.';
        
        if (error.response?.status === 404) {
          this.error = 'User not found.';
        } else if (error.response?.status === 401) {
          this.$router.push('/login');
          return;
        }
      } finally {
        this.loading = false;
      }
    },

    async loadPosts() {
      this.loadingPosts = true;
      this.postsError = null;

      try {
        const response = await this.$axios.get(`/profiles/${this.profileUserId}/posts`);
        this.posts = response.data || [];
      
        if (this.profile && this.profile.PostsCount !== this.posts.length) {
          this.profile.PostsCount = this.posts.length;
        }
        
      } catch (error) {
        console.error('Posts load error:', error);
        this.postsError = 'Failed to load posts.';
        
        if (error.response?.status === 404) {
          this.posts = [];
        } else if (error.response?.status === 401) {
          this.$router.push('/login');
          return;
        }
      } finally {
        this.loadingPosts = false;
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

    async loadBannedUsers() {
      this.loadingBanned = true;
      try {
        const response = await this.$axios.get(`/profiles/${this.currentUserId}/bans`);
        this.bannedUsers = response.data || [];
      } catch (error) {
        console.error('Banned users load error:', error);
        this.bannedUsers = [];
      } finally {
        this.loadingBanned = false;
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

    showBannedUsersList() {
      this.showBannedModal = true;
      this.loadBannedUsers();
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

    async unbanUser(user) {
      if (!confirm(`Are you sure you want to unban ${user.Username}?`)) {
        return;
      }

      this.unbanning = user.UserID;
      try {
        await this.$axios.delete(`/profiles/${this.currentUserId}/bans/${user.UserID}`);
        
        // Rimuovi dalla lista locale
        this.bannedUsers = this.bannedUsers.filter(u => u.UserID !== user.UserID);
        
      } catch (error) {
        console.error('Unban error:', error);
        alert('Failed to unban user. Please try again.');
      } finally {
        this.unbanning = null;
      }
    },

    handleUserBanned() {
      this.$router.push(`/profiles/${this.currentUserId}/feed`);
    },

    closeModal() {
      this.showFollowersModal = false;
      this.showFollowingModal = false;
      this.showEditModal = false;
      this.showBannedModal = false;
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
      
      this.selectedPost = { ...post };
    },

    closePostModal() {
      this.selectedPost = null;
    },

    async onPostDeleted() {
      await Promise.all([
        this.loadProfile(),
        this.loadPosts()
      ]);
      
      this.closePostModal();
    },

    onPostUpdated(updatedData) {
      const postIndex = this.posts.findIndex(p => p.PhotoID === updatedData.PhotoID);
      if (postIndex !== -1) {
        this.posts[postIndex] = {
          ...this.posts[postIndex],
          Nlike: updatedData.Nlike,
          Liked: updatedData.Liked,
          Ncomment: updatedData.Ncomment
        };
        
        this.$forceUpdate();
      }
    },
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
  justify-content: space-between; /* Per allineare il pulsante a destra */
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
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-weight: 600;
  font-size: 18px;
}

.avatar-initials {
  line-height: 1;
}

.username {
  color: #262626;
  font-weight: 500;
  margin-left: 12px;
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