<template>
  <div class="profile-stats">
    <div class="row text-center">
      <div class="col-4">
        <div class="stat-item" @click="$emit('showPosts')" role="button">
          <div class="stat-number">{{ profile.PostsCount || 0 }}</div>
          <div class="stat-label">{{ profile.PostsCount === 1 ? 'post' : 'posts' }}</div>
        </div>
      </div>
      
      <div class="col-4">
        <div 
          class="stat-item" 
          @click="isOwnProfile ? $emit('showFollowers') : null" 
          :class="{ 'clickable': isOwnProfile }"
          role="button"
        >
          <div class="stat-number">{{ formatNumber(profile.Follower || 0) }}</div>
          <div class="stat-label">{{ profile.Follower === 1 ? 'follower' : 'followers' }}</div>
        </div>
      </div>
      
      <div class="col-4">
        <div 
          class="stat-item" 
          @click="isOwnProfile ? $emit('showFollowing') : null"
          :class="{ 'clickable': isOwnProfile }"
          role="button"
        >
          <div class="stat-number">{{ formatNumber(profile.Following || 0) }}</div>
          <div class="stat-label">following</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProfileStats',
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
  methods: {
    formatNumber(num) {
      if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M'
      }
      if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K'
      }
      return num.toString()
    }
  }
}
</script>

<style scoped>
.profile-stats {
  margin: 20px 0;
}

.stat-item {
  transition: opacity 0.2s;
}

.stat-item.clickable {
  cursor: pointer;
}

.stat-item.clickable:hover {
  opacity: 0.7;
}

.stat-number {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
}

.stat-label {
  font-size: 14px;
  color: #8e8e8e;
  margin-top: 2px;
}

@media (max-width: 768px) {
  .stat-number {
    font-size: 16px;
  }
  
  .stat-label {
    font-size: 12px;
  }
}
</style>