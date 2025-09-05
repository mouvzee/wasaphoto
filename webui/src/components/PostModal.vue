<template>
  <div class="modal-overlay" @click="close">
    <div class="modal-content post-modal" @click.stop>
      <!-- Header con X e menu opzioni -->
      <div class="modal-header">
        <div class="header-actions">
          <button class="btn-close-modal" @click="close">
            <i class="fas fa-times"></i>
          </button>
          
          <!-- Menu opzioni (solo per il proprietario) -->
          <div v-if="isOwner" class="post-options">
            <button 
              class="btn btn-link p-0"
              @click="toggleOptionsMenu"
            >
              <i class="fas fa-ellipsis-h"></i>
            </button>
            
            <div v-if="showOptionsMenu" class="options-menu">
              <button 
                class="dropdown-item text-danger"
                @click="confirmDelete"
              >
                <i class="fas fa-trash me-2"></i>
                Delete Post
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <LoadingSpinner :loading="loading">
        <div v-if="post" class="post-container">
          <!-- Immagine -->
          <div class="post-image-container">
            <img 
              :src="getImageUrl(post.ImageData)" 
              :alt="post.Caption"
              class="post-image"
            >
          </div>
          
          <!-- Sidebar con info e commenti -->
          <div class="post-sidebar">
            <!-- Header utente -->
            <div class="post-user-header">
              <div class="user-info">
                <div class="user-avatar" :style="avatarStyle">
                  <span class="avatar-initials">{{ userInitials }}</span>
                </div>
                <RouterLink 
                  :to="`/profiles/${post.UserID}`" 
                  class="username"
                  @click="close"
                >
                  {{ post.Username }}
                </RouterLink>
              </div>
            </div>
            
            <!-- Caption -->
            <div v-if="post.Caption" class="post-caption">
              <div class="caption-content">
                <RouterLink 
                  :to="`/profiles/${post.UserID}`" 
                  class="username me-2"
                  @click="close"
                >
                  {{ post.Username }}
                </RouterLink>
                {{ post.Caption }}
              </div>
              <div class="post-date">
                {{ formatDate(post.CreatedAt) }}
              </div>
            </div>
            
            <!-- Commenti -->
            <div class="comments-section">
              <div class="comments-list">
                <div 
                  v-for="comment in comments" 
                  :key="comment.CommentID"
                  class="comment-item"
                >
                  <div class="comment-content">
                    <RouterLink 
                      :to="`/profiles/${comment.UserID}`" 
                      class="comment-username"
                      @click="close"
                    >
                      {{ comment.Username }}
                    </RouterLink>
                    <span class="comment-text">{{ comment.Lyric }}</span>
                  </div>
                  <div class="comment-meta">
                    <span class="comment-date">{{ formatDate(comment.Created_At) }}</span>
                    
                    <!-- Menu opzioni commento -->
                    <div 
                      v-if="canDeleteComment(comment)" 
                      class="comment-options"
                    >
                      <button 
                        class="btn btn-link p-0"
                        @click="toggleCommentMenu(comment.CommentID)"
                      >
                        <i class="fas fa-ellipsis-h"></i>
                      </button>
                      
                      <div 
                        v-if="showCommentMenu === comment.CommentID" 
                        class="comment-options-menu"
                      >
                        <button 
                          class="dropdown-item text-danger"
                          @click="confirmDeleteComment(comment)"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
                
                <!-- Messaggio se non ci sono commenti -->
                <div v-if="comments.length === 0 && !loadingComments" class="no-comments">
                  <p class="text-muted text-center">No comments yet</p>
                </div>
              </div>
            </div>
            
            <!-- Azioni e form commento -->
            <div class="post-actions">
              <!-- Pulsanti like e commento -->
              <div class="action-buttons">
                <button 
                  class="btn btn-link p-0 me-3"
                  @click="toggleLike"
                  :class="{ 'text-danger': post.Liked }"
                  :disabled="likingInProgress"
                >
                  <i class="fas fa-heart fa-lg"></i>
                </button>
                <button 
                  class="btn btn-link p-0"
                  @click="focusCommentInput"
                >
                  <i class="fas fa-comment fa-lg"></i>
                </button>
              </div>
              
              <div class="likes-count">
                <strong>{{ post.Nlike }} {{ post.Nlike === 1 ? 'like' : 'likes' }}</strong>
              </div>
              
              <!-- Form nuovo commento -->
              <div class="comment-form">
                <form @submit.prevent="addComment">
                  <div class="input-group">
                    <input 
                      ref="commentInput"
                      v-model="newComment"
                      type="text"
                      class="form-control comment-input"
                      placeholder="Add a comment..."
                      maxlength="255"
                      :disabled="submittingComment"
                    >
                    <button 
                      type="submit"
                      class="btn btn-link"
                      :disabled="!newComment.trim() || submittingComment"
                    >
                      <span v-if="submittingComment">
                        <i class="fas fa-spinner fa-spin"></i>
                      </span>
                      <span v-else>Post</span>
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        </div>
      </LoadingSpinner>
      
      <ErrorMsg v-if="error" :msg="error" />
    </div>
  </div>
</template>

<script>
export default {
  name: 'PostModal',
  props: {
    postData: {
      type: Object,
      required: true
    }
  },
  emits: ['close', 'deleted', 'postUpdated'],
  data() {
    return {
      post: null,
      loading: false,
      error: null,
      showOptionsMenu: false,
      
      // Commenti
      comments: [],
      newComment: '',
      loadingComments: false,
      submittingComment: false,
      showCommentMenu: null,
      
      // Like
      likingInProgress: false
    }
  },
  computed: {
    currentUserId() {
      try {
        const user = JSON.parse(localStorage.getItem('user'));
        return user ? user.UserID : null;
      } catch (error) {
        return null;
      }
    },
    currentUsername() {
      try {
        const user = JSON.parse(localStorage.getItem('user'));
        return user ? user.Username : null;
      } catch (error) {
        return null;
      }
    },
    isOwner() {
      return this.post && this.currentUserId && this.post.UserID === this.currentUserId;
    },
    // ✅ NUOVO: Avatar initials come in FeedPost
    userInitials() {
      if (!this.post?.Username) return '?';
      return this.post.Username.charAt(0).toUpperCase();
    },
    
    // ✅ NUOVO: Avatar style come in FeedPost
    avatarStyle() {
      const username = this.post?.Username || 'default';
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
  async created() {
    this.post = {
      ...this.postData,
      Nlike: this.postData.Nlike || 0,
      Ncomment: this.postData.Ncomment || 0,
      Liked: this.postData.Liked || false,
      Username: this.postData.Username || 'Unknown User'
    };
    
    await this.loadPostData();
    await this.loadComments();
    
    document.addEventListener('click', this.closeAllMenus);
  },
  beforeUnmount() {
    document.removeEventListener('click', this.closeAllMenus);
  },
  methods: {
    close() {
      this.$emit('close');
    },

    async loadPostData() {
      try {
        const response = await this.$axios.get(`/profiles/${this.post.UserID}/posts`);
        const posts = response.data || [];
        const currentPost = posts.find(p => p.PhotoID === this.post.PhotoID);
        
        if (currentPost) {
          this.post.Nlike = currentPost.Nlike || 0;
          this.post.Liked = currentPost.Liked || false;
          this.post.Caption = this.post.Caption || currentPost.Caption;
          this.post.CreatedAt = this.post.CreatedAt || currentPost.CreatedAt;
        }
      } catch (error) {
        console.error('Error loading updated post data:', error);
      }
    },

    async loadComments() {
      this.loadingComments = true;
      try {
        const response = await this.$axios.get(`/profiles/${this.post.UserID}/posts/${this.post.PhotoID}/comments`);
        const rawComments = response.data || [];
        
        const processedComments = rawComments.map(comment => ({
          CommentID: comment.CommentID,
          Lyric: comment.Lyric,
          Created_At: comment.Created_At,
          UserID: comment.User?.UserID || null,
          Username: comment.User?.Username || 'Unknown User'
        }));
        
        this.comments = processedComments;
        this.post.Ncomment = this.comments.length;
      } catch (error) {
        console.error('Error loading comments:', error);
        if (error.response?.status === 404) {
          this.comments = [];
          this.post.Ncomment = 0;
        }
      } finally {
        this.loadingComments = false;
      }
    },

    async toggleLike() {
      if (!this.post || !this.currentUserId || this.likingInProgress) return;

      this.likingInProgress = true;

      try {
        if (this.post.Liked) {
          await this.$axios.delete(`/profiles/${this.post.UserID}/posts/${this.post.PhotoID}/likes/0`);
          this.post.Liked = false;
          this.post.Nlike = Math.max(0, this.post.Nlike - 1);
        } else {
          await this.$axios.put(`/profiles/${this.post.UserID}/posts/${this.post.PhotoID}/likes/0`);
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
        this.error = 'Failed to toggle like. Please try again.';
      } finally {
        this.likingInProgress = false;
      }
    },

    toggleOptionsMenu() {
      this.showOptionsMenu = !this.showOptionsMenu;
      this.showCommentMenu = null;
    },

    toggleCommentMenu(commentId) {
      this.showCommentMenu = this.showCommentMenu === commentId ? null : commentId;
      this.showOptionsMenu = false;
    },

    closeAllMenus() {
      this.showOptionsMenu = false;
      this.showCommentMenu = null;
    },

    confirmDelete() {
      if (confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
        this.deletePost();
      }
      this.showOptionsMenu = false;
    },

    async deletePost() {
      try {
        await this.$axios.delete(`/profiles/${this.post.UserID}/posts/${this.post.PhotoID}`);
        this.$emit('deleted');
        this.close();
      } catch (error) {
        console.error('Error deleting post:', error);
        this.error = 'Failed to delete post. Please try again.';
      }
    },

    async addComment() {
      if (!this.newComment.trim() || !this.currentUserId) {
        return;
      }

      this.submittingComment = true;

      try {
        const response = await this.$axios.post(
          `/profiles/${this.post.UserID}/posts/${this.post.PhotoID}/comments`,
          {
            Lyric: this.newComment.trim()
          }
        );

        const newCommentData = {
          CommentID: response.data.CommentID || Date.now(),
          Lyric: this.newComment.trim(),
          UserID: this.currentUserId,
          Username: this.currentUsername || 'You',
          Created_At: new Date().toISOString()
        };

        this.comments.push(newCommentData);
        this.newComment = '';
        this.post.Ncomment = this.comments.length;

        this.$emit('postUpdated', {
          PhotoID: this.post.PhotoID,
          Nlike: this.post.Nlike,
          Liked: this.post.Liked,
          Ncomment: this.post.Ncomment
        });

        this.$nextTick(() => {
          const container = this.$el.querySelector('.comments-list');
          if (container) {
            container.scrollTop = container.scrollHeight;
          }
        });

      } catch (error) {
        console.error('Error adding comment:', error);
        this.error = 'Failed to add comment. Please try again.';
      } finally {
        this.submittingComment = false;
      }
    },

    canDeleteComment(comment) {
      if (!this.currentUserId) return false;
      return comment.UserID === this.currentUserId || this.post.UserID === this.currentUserId;
    },

    confirmDeleteComment(comment) {
      if (confirm('Are you sure you want to delete this comment?')) {
        this.deleteComment(comment);
      }
      this.showCommentMenu = null;
    },

    async deleteComment(comment) {
      try {
        await this.$axios.delete(
          `/profiles/${this.post.UserID}/posts/${this.post.PhotoID}/comments/${comment.CommentID}`
        );

        this.comments = this.comments.filter(c => c.CommentID !== comment.CommentID);
        this.post.Ncomment = this.comments.length;

        this.$emit('postUpdated', {
          PhotoID: this.post.PhotoID,
          Nlike: this.post.Nlike,
          Liked: this.post.Liked,
          Ncomment: this.post.Ncomment
        });

      } catch (error) {
        console.error('Error deleting comment:', error);
        if (error.response?.status === 404) {
          this.comments = this.comments.filter(c => c.CommentID !== comment.CommentID);
          this.post.Ncomment = this.comments.length;
          
          this.$emit('postUpdated', {
            PhotoID: this.post.PhotoID,
            Nlike: this.post.Nlike,
            Liked: this.post.Liked,
            Ncomment: this.post.Ncomment
          });
        } else {
          this.error = 'Failed to delete comment. Please try again.';
        }
      }
    },

    focusCommentInput() {
      this.$nextTick(() => {
        if (this.$refs.commentInput) {
          this.$refs.commentInput.focus();
        }
      });
    },

    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      const now = new Date();
      const diffTime = Math.abs(now - date);
      const diffMinutes = Math.floor(diffTime / (1000 * 60));
      const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
      const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
      
      if (diffMinutes < 1) return 'now';
      if (diffMinutes < 60) return `${diffMinutes}m`;
      if (diffHours < 24) return `${diffHours}h`;
      if (diffDays < 7) return `${diffDays}d`;
      return `${Math.floor(diffDays / 7)}w`;
    },

    getImageUrl(imageData) {
      if (!imageData) return '/placeholder-image.jpg';
      
      if (typeof imageData === 'string' && imageData.startsWith('http')) {
        return imageData;
      }
      
      if (typeof imageData === 'string' && imageData.startsWith('data:image')) {
        return imageData;
      }
      
      if (typeof imageData === 'string') {
        return `data:image/jpeg;base64,${imageData}`;
      }
      
      try {
        const blob = new Blob([new Uint8Array(imageData)], { type: 'image/jpeg' });
        return URL.createObjectURL(blob);
      } catch (error) {
        console.error('Error creating image URL:', error);
        return '/placeholder-image.jpg';
      }
    }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.post-modal {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 1000px;
  max-height: 90vh;
  position: relative;
  overflow: hidden;
}

.modal-header {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 20;
  padding: 15px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-close-modal {
  background: rgba(0, 0, 0, 0.5);
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.post-container {
  display: flex;
  height: 80vh;
}

.post-image-container {
  flex: 1;
  background: black;
  display: flex;
  align-items: center;
  justify-content: center;
}

.post-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.post-sidebar {
  width: 350px;
  display: flex;
  flex-direction: column;
  border-left: 1px solid #e9ecef;
}

.post-user-header {
  padding: 16px;
  border-bottom: 1px solid #e9ecef;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-weight: 600;
  font-size: 14px;
}

.avatar-initials {
  line-height: 1;
}

.username {
  font-weight: 600;
  color: #262626;
  text-decoration: none;
}

.post-options {
  position: relative;
}

.post-options .btn {
  background: rgba(0, 0, 0, 0.5);
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.options-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 8px 16px;
  background: none;
  border: none;
  text-align: left;
  color: #262626;
  white-space: nowrap;
}

.dropdown-item:hover {
  background-color: #f8f9fa;
}

.post-caption {
  padding: 16px;
  border-bottom: 1px solid #e9ecef;
}

.caption-content {
  margin-bottom: 8px;
}

.post-date {
  color: #8e8e8e;
  font-size: 12px;
  text-transform: uppercase;
}

.comments-section {
  flex: 1;
  overflow-y: auto;
  padding: 0 16px;
}

.comments-list {
  max-height: 100%;
}

.comment-item {
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-content {
  margin-bottom: 4px;
}

.comment-username {
  font-weight: 600;
  color: #262626;
  text-decoration: none;
  margin-right: 8px;
}

.comment-username:hover {
  text-decoration: underline;
}

.comment-text {
  color: #262626;
  word-wrap: break-word;
}

.comment-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.comment-date {
  color: #8e8e8e;
  font-size: 12px;
}

.comment-options {
  position: relative;
}

.comment-options .btn {
  color: #8e8e8e;
  font-size: 12px;
}

.comment-options-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
}

.no-comments {
  padding: 20px 0;
}

.post-actions {
  border-top: 1px solid #e9ecef;
  padding: 16px;
}

.action-buttons {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.btn-link {
  color: #262626;
  text-decoration: none;
}

.btn-link:hover {
  color: #007bff;
}

.btn-link.text-danger {
  color: #dc3545 !important;
}

.btn-link:disabled {
  opacity: 0.6;
  pointer-events: none;
}

.likes-count {
  color: #262626;
  margin-bottom: 12px;
}

.comment-form {
  margin-top: 12px;
}

.comment-input {
  border: none;
  outline: none;
  box-shadow: none;
  font-size: 14px;
  border-radius: 0;
}

.comment-input:focus {
  box-shadow: none;
  border-color: transparent;
}

.input-group .btn {
  border: none;
  color: #007bff;
  font-weight: 600;
  font-size: 14px;
}

.input-group .btn:disabled {
  color: #8e8e8e;
}

.input-group .btn:not(:disabled):hover {
  color: #0056b3;
}

@media (max-width: 768px) {
  .post-modal {
    width: 95%;
    height: 95vh;
  }

  .post-container {
    flex-direction: column;
  }

  .post-sidebar {
    width: 100%;
    height: 50%;
  }

  .post-image-container {
    height: 50%;
  }
  
  .modal-header {
    position: fixed;
    top: 10px;
    right: 10px;
  }
}
</style>