<template>
  <div class="post-card" @click="openPost">
    <div class="post-thumbnail">
      <img 
        :src="getImageUrl(post.ImageData)" 
        :alt="post.Caption"
        class="post-image"
        @load="onImageLoad"
        @error="onImageError"
      >
      
      <!-- Overlay con statistiche (visibile solo on hover) -->
      <div class="post-overlay">
        <div class="post-stats">
          <span class="stat-item">
            <i class="fas fa-heart"></i>
            {{ post.Nlike || 0 }}
          </span>
          <span class="stat-item">
            <i class="fas fa-comment"></i>
            {{ post.Ncomment || 0 }}
          </span>
        </div>
      </div>
      
      <!-- Indicatore se è un post multiplo (se implementi in futuro) -->
      <div v-if="post.IsMultiple" class="multiple-indicator">
        <i class="fas fa-clone"></i>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PostCard',
  props: {
    post: {
      type: Object,
      required: true
    }
  },
  emits: ['click'],
  data() {
    return {
      imageLoaded: false,
      imageError: false
    }
  },
  methods: {
    openPost() {
      this.$emit('click', this.post);
    },

    onImageLoad() {
      this.imageLoaded = true;
      this.imageError = false;
    },

    onImageError() {
      this.imageError = true;
      this.imageLoaded = false;
    },

    getImageUrl(imageData) {
      if (!imageData) {
        return '/placeholder-image.jpg';
      }
      
      if (typeof imageData === 'string' && imageData.startsWith('http')) {
        return imageData;
      }
      
      // Se imageData è una stringa base64
      if (typeof imageData === 'string' && imageData.startsWith('data:image')) {
        return imageData;
      }
      
      // Se imageData è una stringa base64 senza prefisso
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
.post-card {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  background-color: #f8f9fa;
}

.post-thumbnail {
  position: relative;
  aspect-ratio: 1;
  overflow: hidden;
}

.post-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.2s ease;
}

.post-card:hover .post-image {
  transform: scale(1.05);
}

.post-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.post-card:hover .post-overlay {
  opacity: 1;
}

.post-stats {
  color: white;
  text-align: center;
  font-weight: 600;
}

.stat-item {
  display: inline-block;
  margin: 0 12px;
  font-size: 16px;
}

.stat-item i {
  margin-right: 6px;
}

.multiple-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
  color: white;
  font-size: 16px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
}

/* Loading state */
.post-image[src="/placeholder-image.jpg"] {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: loading 1.5s infinite;
}

@keyframes loading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* Responsive */
@media (max-width: 768px) {
  .stat-item {
    margin: 0 8px;
    font-size: 14px;
  }
  
  .multiple-indicator {
    font-size: 14px;
  }
}
</style>