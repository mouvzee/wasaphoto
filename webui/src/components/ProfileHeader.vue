<template>
  <div class="profile-header py-4">
    <div class="row">
      <!-- Avatar placeholder (sinistra) -->
      <div class="col-md-4 text-center">
        <div class="profile-avatar">
          <div class="avatar-circle">
            <i class="fas fa-user fa-4x text-muted"></i>
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
                <button class="btn btn-outline-secondary dropdown-toggle" type="button" data-bs-toggle="dropdown">
                  <i class="fas fa-ellipsis-h"></i>
                </button>
                <ul class="dropdown-menu">
                  <li><a class="dropdown-item" href="#" @click="banUser">Ban User</a></li>
                </ul>
              </div>
            </div>

            <!-- Pulsante settings per il proprio profilo -->
            <div v-else class="action-buttons">
              <button class="btn btn-outline-secondary" @click="editProfile">
                Edit Profile
              </button>
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

.avatar-circle {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  border: 1px solid #e1e8ed;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  background-color: #f8f9fa;
}

.profile-username {
  font-size: 28px;
  font-weight: 300;
  margin: 0;
}

.action-buttons .btn {
  font-weight: 600;
  padding: 8px 24px;
}

@media (max-width: 768px) {
  .profile-header .row {
    text-align: center;
  }
  
  .profile-username {
    font-size: 24px;
    margin-bottom: 1rem;
  }

  .action-buttons {
    margin-bottom: 1rem;
  }
}
</style>