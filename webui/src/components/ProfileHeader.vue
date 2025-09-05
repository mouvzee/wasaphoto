<template>
  <div class="profile-header py-4">
    <div class="row">
      <!-- Avatar placeholder (sinistra) -->
      <div class="col-md-4 text-center">
        <div class="profile-avatar">
          <div class="avatar-circle" :style="avatarStyle">
            <!-- ✅ Mostra iniziali invece dell'icona -->
            <span class="avatar-initials">{{ userInitials }}</span>
          </div>
        </div>
      </div>

      <!-- Info profilo (destra) -->
      <div class="col-md-8">
        <div class="profile-info">
          <!-- Username e pulsanti -->
          <div class="d-flex align-items-center mb-3">
            <h1 class="profile-username me-4">{{ profile.User.Username }}</h1>
            
            <!-- Pulsanti azione (solo se non è il proprio profilo) -->
            <div v-if="!isOwnProfile" class="action-buttons">
              <button 
                v-if="!profile.IsFollowed" 
                @click="followUser"
                class="btn btn-primary me-2"
                :disabled="loading"
              >
                Follow
              </button>
              <button 
                v-else 
                @click="unfollowUser"
                class="btn btn-outline-secondary me-2"
                :disabled="loading"
              >
                Following
              </button>
              
              <div class="dropdown d-inline">
                <button class="btn-dots" type="button" data-bs-toggle="dropdown">
                  <i class="fas fa-ellipsis-h"></i>
                </button>
                <ul class="dropdown-menu">
                  <li><a class="dropdown-item" href="#" @click.prevent.stop="banUser">Ban User</a></li>
                </ul>
              </div>
            </div>

            <!-- Pulsanti per il proprio profilo -->
            <div v-else class="action-buttons">
              <button class="btn btn-outline-secondary me-2" @click="editProfile">
                Edit Profile
              </button>
              <div class="dropdown d-inline">
                <button class="btn-dots" type="button" data-bs-toggle="dropdown">
                  <i class="fas fa-ellipsis-h"></i>
                </button>
                <ul class="dropdown-menu">
                  <li><a class="dropdown-item" href="#" @click.prevent.stop="showBannedUsers">Banned Users</a></li>
                </ul>
              </div>
            </div>
          </div>

          <!-- Stats bar integrata nell'header -->
          <ProfileStats 
            :profile="profile" 
            :isOwnProfile="isOwnProfile"
            @showFollowers="$emit('showFollowers')"
            @showFollowing="$emit('showFollowing')"
            @showPosts="$emit('showPosts')"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProfileHeader',
  props: {
    profile: {
      type: Object,
      required: true
    },
    isOwnProfile: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      loading: false
    }
  },
  computed: {
    userInitials() {
      if (!this.profile?.User?.Username) return '?';
      const username = this.profile.User.Username;
      return username.charAt(0).toUpperCase();
    },
    
    avatarStyle() {
      // Genera un colore basato sul username
      const username = this.profile?.User?.Username || 'default';
      const colors = [
        '#e74c3c', '#3498db', '#2ecc71', '#f39c12', 
        '#9b59b6', '#1abc9c', '#34495e', '#e67e22'
      ];
      
      let hash = 0;
      for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
      }
      
      const color = colors[Math.abs(hash) % colors.length];
      
      return {
        backgroundColor: color,
        color: 'white'
      };
    }
  },
  methods: {
    async followUser() {
      this.loading = true
      try {
        await this.$axios.put(`/profiles/${this.getCurrentUserId()}/followings/${this.profile.User.UserID}`)
        this.$emit('profileUpdated')
      } catch (error) {
        console.error('Follow error:', error)
      } finally {
        this.loading = false
      }
    },

    async unfollowUser() {
      this.loading = true
      try {
        await this.$axios.delete(`/profiles/${this.getCurrentUserId()}/followings/${this.profile.User.UserID}`)
        this.$emit('profileUpdated')
      } catch (error) {
        console.error('Unfollow error:', error)
      } finally {
        this.loading = false
      }
    },

    async banUser() {
      if (confirm(`Are you sure you want to ban ${this.profile.User.Username}?`)) {
        try {
          await this.$axios.put(`/profiles/${this.getCurrentUserId()}/bans/${this.profile.User.UserID}`)
          this.$emit('userBanned')
        } catch (error) {
          console.error('Ban error:', error)
        }
      }
    },

    editProfile() {
      this.$emit('editProfile')
    },

    showBannedUsers() {
      this.$emit('showBannedUsers')
    },

    getCurrentUserId() {
      const user = JSON.parse(localStorage.getItem('user'))
      return user.UserID
    }
  }
}
</script>

<style scoped>
.profile-header {
  border-bottom: 1px solid #e1e8ed;
}

.profile-avatar .avatar-circle {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background-color: #f8f9fa;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.profile-username {
  font-size: 28px;
  font-weight: 300;
  margin: 0;
  color: #262626;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-dots {
  background: none;
  border: none;
  color: #262626;
  font-size: 16px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: background-color 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
}

.btn-dots:hover {
  background-color: #f5f5f5;
  color: #262626;
}

.btn-dots:focus {
  outline: none;
  box-shadow: none;
}

.btn-dots:active {
  background-color: #e5e5e5;
}

/* Dropdown menu styling */
.dropdown-menu {
  border: 1px solid #e1e8ed;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 8px 0;
  margin-top: 4px;
}

.dropdown-item {
  padding: 8px 16px;
  color: #262626;
  font-size: 14px;
}

.dropdown-item:hover {
  background-color: #f5f5f5;
  color: #262626;
}

.avatar-initials {
  font-size: 3rem;
  font-weight: 600;
  line-height: 1;
}

/* Rimuovi l'icona di default */
.profile-avatar .avatar-circle {
  background-color: transparent; /* Sarà sovrascritto dallo style computed */
}

@media (max-width: 768px) {
  .profile-username {
    font-size: 20px;
  }
  
  .profile-avatar .avatar-circle {
    width: 100px;
    height: 100px;
  }
  
  .profile-avatar .avatar-circle i {
    font-size: 2.5rem !important;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
  
  .action-buttons .btn {
    width: 100%;
  }
}
</style>