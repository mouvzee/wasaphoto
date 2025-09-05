<!-- filepath: /home/void/Uni/wasaphoto/webui/src/components/FeedPost.vue -->
<template>
  <div class="feed-post">
    <!-- Header del post -->
    <div class="post-header">
      <div class="user-info">
        <div class="user-avatar" :style="avatarStyle">
          <span class="avatar-initials">{{ userInitials }}</span>
        </div>
        <div class="user-details">
          <RouterLink 
            :to="`/profiles/${postUserId}`" 
            class="username"
          >
            {{ displayUsername }}
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Immagine del post -->
    <div class="post-image-container">
      <img 
        :src="`data:image/jpeg;base64,${post.ImageData}`"
        :alt="post.Caption"
        class="post-image"
        @click="openPostModal"
      />
    </div>

    <!-- Azioni del post -->
    <div class="post-actions">
      <div class="action-buttons">
        <button 
          class="action-btn like-btn"
          @click="toggleLike"
          :class="{ 'liked': post.Liked }"
          :disabled="likingInProgress"
        >
          <i class="fas fa-heart"></i>
        </button>
        <button class="action-btn" @click="openPostModal">
          <i class="fas fa-comment"></i>
        </button>
      </div>
    </div>

    <!-- Conteggi like -->
    <div class="post-stats">
      <div class="likes-count" v-if="post.Nlike > 0">
        <strong>{{ post.Nlike }} {{ post.Nlike === 1 ? 'like' : 'likes' }}</strong>
      </div>
    </div>

    <!-- Caption -->
    <div class="post-caption" v-if="post.Caption">
      <RouterLink :to="`/profiles/${postUserId}`" class="username">
        {{ displayUsername }}
      </RouterLink>
      <span class="caption-text">{{ post.Caption }}</span>
    </div>

    <!-- Link ai commenti -->
    <div class="comments-preview" v-if="post.Ncomment > 0">
      <button class="view-comments-btn" @click="openPostModal">
        View {{ post.Ncomment === 1 ? 'comment' : `all ${post.Ncomment} comments` }}
      </button>
    </div>

    <!-- Timestamp -->
    <div class="post-timestamp">
      {{ formatTime(post.CreatedAt) }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'FeedPost',
  props: {
    post: {
      type: Object,
      required: true
    }
  },
  emits: ['openModal', 'postUpdated'],
  data() {
    return {
      likingInProgress: false,
      showOptions: false
    }
  },
  computed: {
    userInitials() {
      const username = this.post?.Username || 
                      this.post?.User?.Username || 
                      this.post?.username || 
                      this.post?.user?.username;
      if (!username) return '?';
      return username.charAt(0).toUpperCase();
    },
    
    avatarStyle() {
      const username = this.post?.Username || 
                      this.post?.User?.Username || 
                      this.post?.username || 
                      this.post?.user?.username || 
                      'default';
      
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
    },

    displayUsername() {
      return this.post?.Username || 
             this.post?.User?.Username || 
             this.post?.username || 
             this.post?.user?.username || 
             'Unknown User';
    },

    postUserId() {
      return this.post?.UserID || 
             this.post?.User?.UserID || 
             this.post?.userId || 
             this.post?.user?.userID ||
             this.post?.user?.UserID;
    },

    currentUserId() {
      try {
        const user = JSON.parse(localStorage.getItem('user'));
        return user ? user.UserID : null;
      } catch (error) {
        return null;
      }
    }
  },
  methods: {
    async toggleLike() {
      if (!this.post || !this.currentUserId || this.likingInProgress) return;

      this.likingInProgress = true;

      try {
        if (this.post.Liked) {
          await this.$axios.delete(`/profiles/${this.postUserId}/posts/${this.post.PhotoID}/likes/0`);
          this.post.Liked = false;
          this.post.Nlike = Math.max(0, this.post.Nlike - 1);
        } else {
          await this.$axios.put(`/profiles/${this.postUserId}/posts/${this.post.PhotoID}/likes/0`);
          this.post.Liked = true;
          this.post.Nlike++;
        }
        
        this.$emit('postUpdated', {
          PhotoID: this.post.PhotoID,
          Nlike: this.post.Nlike,
          Liked: this.post.Liked,
          Ncomment: this.post.Ncomment
        });
        
      } catch (error) {
        console.error('Error toggling like:', error);
      } finally {
        this.likingInProgress = false;
      }
    },

    openPostModal() {
      const postData = {
        ...this.post,
        UserID: this.postUserId,
        Username: this.displayUsername
      };
      
      this.$emit('openModal', postData);
    },

    formatTime(timestamp) {
      if (!timestamp) return '';
      
      const date = new Date(timestamp);
      const now = new Date();
      const diffInSeconds = Math.floor((now - date) / 1000);
      
      if (diffInSeconds < 60) {
        return 'now';
      } else if (diffInSeconds < 3600) {
        const minutes = Math.floor(diffInSeconds / 60);
        return `${minutes}m`;
      } else if (diffInSeconds < 86400) {
        const hours = Math.floor(diffInSeconds / 3600);
        return `${hours}h`;
      } else if (diffInSeconds < 604800) {
        const days = Math.floor(diffInSeconds / 86400);
        return `${days}d`;
      } else {
        return date.toLocaleDateString();
      }
    }
  }
}
</script>

<style scoped>
.feed-post {
  background: white;
  border: 1px solid #dbdbdb;
  border-radius: 8px;
  margin-bottom: 24px;
  max-width: 600px; /* ✅ AUMENTATO: da 470px a 600px */
  width: 100%;
}

.post-header {
  display: flex;
  align-items: center;
  /* ✅ RIMOSSO: justify-content: space-between; */
  padding: 16px;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 40px; /* ✅ AUMENTATO: da 32px a 40px */
  height: 40px; /* ✅ AUMENTATO: da 32px a 40px */
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-weight: 600;
  font-size: 16px; /* ✅ AUMENTATO: da 14px a 16px */
}

.avatar-initials {
  line-height: 1;
}

.username {
  font-weight: 600;
  color: #262626;
  text-decoration: none;
  font-size: 16px; /* ✅ AUMENTATO: da 14px a 16px */
}

.username:hover {
  color: #262626;
  text-decoration: none;
}

.post-image-container {
  position: relative;
  width: 100%;
  background: #f8f9fa;
}

.post-image {
  width: 100%;
  height: auto;
  max-height: 700px; /* ✅ AUMENTATO: da 600px a 700px */
  object-fit: cover;
  cursor: pointer;
  display: block;
}

.post-actions {
  padding: 16px; /* ✅ AUMENTATO: da 12px 16px 4px a 16px */
}

.action-buttons {
  display: flex;
  gap: 20px; /* ✅ AUMENTATO: da 16px a 20px */
}

.action-btn {
  background: none;
  border: none;
  color: #262626;
  cursor: pointer;
  padding: 8px;
  font-size: 28px; /* ✅ AUMENTATO: da 24px a 28px */
  transition: color 0.2s ease;
}

.action-btn:hover {
  color: #8e8e8e;
}

.like-btn.liked {
  color: #ed4956;
  animation: likeAnimation 0.3s ease;
}

@keyframes likeAnimation {
  0% { transform: scale(1); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.post-stats {
  padding: 0 16px;
  margin-bottom: 12px; /* ✅ AUMENTATO: da 8px a 12px */
}

.likes-count {
  font-size: 16px; /* ✅ AUMENTATO: da 14px a 16px */
  color: #262626;
}

.post-caption {
  padding: 0 16px;
  margin-bottom: 8px; /* ✅ AUMENTATO: da 4px a 8px */
  font-size: 16px; /* ✅ AUMENTATO: da 14px a 16px */
  line-height: 20px; /* ✅ AUMENTATO: da 18px a 20px */
}

.caption-text {
  color: #262626;
  margin-left: 4px;
}

.comments-preview {
  padding: 0 16px;
  margin-bottom: 12px; /* ✅ AUMENTATO: da 8px a 12px */
}

.view-comments-btn {
  background: none;
  border: none;
  color: #8e8e8e;
  cursor: pointer;
  font-size: 16px; /* ✅ AUMENTATO: da 14px a 16px */
  padding: 0;
}

.view-comments-btn:hover {
  color: #262626;
}

.post-timestamp {
  padding: 0 16px 16px;
  font-size: 12px; /* ✅ MANTENUTO */
  color: #8e8e8e;
  text-transform: uppercase;
  letter-spacing: 0.2px;
}

@media (max-width: 768px) {
  .feed-post {
    border-radius: 0;
    border-left: none;
    border-right: none;
    margin-bottom: 0;
    max-width: none;
  }
  
  .post-image {
    max-height: 400px;
  }
}
</style>