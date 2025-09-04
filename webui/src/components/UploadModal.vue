<template>
  <div class="modal-overlay" @click="close">
    <div class="modal-content upload-modal" @click.stop>
      <div class="modal-header">
        <h5 class="modal-title">Create New Post</h5>
        <button class="btn-close" @click="close">&times;</button>
      </div>
      
      <div class="modal-body">
        <form @submit.prevent="uploadPost">
          <!-- Image Upload -->
          <div class="mb-4">
            <label class="form-label">Choose Photo</label>
            <div 
              class="image-upload-area"
              :class="{ 'has-image': selectedFile }"
              @click="triggerFileInput"
              @dragover.prevent
              @drop.prevent="handleDrop"
            >
              <input 
                ref="fileInput"
                type="file"
                accept="image/*"
                @change="handleFileSelect"
                class="d-none"
              >
              
              <div v-if="!selectedFile" class="upload-placeholder">
                <i class="fas fa-cloud-upload-alt fa-3x mb-3"></i>
                <p class="mb-2">Drag and drop your photo here</p>
                <p class="text-muted">or click to browse</p>
              </div>
              
              <div v-else class="image-preview">
                <img :src="imagePreview" alt="Preview" />
                <button 
                  type="button" 
                  class="btn btn-sm btn-danger remove-image"
                  @click.stop="removeImage"
                >
                  <i class="fas fa-times"></i>
                </button>
              </div>
            </div>
          </div>
          
          <!-- Caption -->
          <div class="mb-4">
            <label for="caption" class="form-label">Caption</label>
            <textarea 
              id="caption"
              v-model="caption"
              class="form-control"
              rows="3"
              placeholder="Write a caption..."
              maxlength="500"
            ></textarea>
            <div class="form-text text-end">
              {{ caption.length }}/500
            </div>
          </div>
          
          <ErrorMsg v-if="error" :msg="error" />
          
          <div class="d-flex justify-content-end gap-2">
            <button type="button" class="btn btn-secondary" @click="close">
              Cancel
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="!selectedFile || uploading"
            >
              <span v-if="uploading">
                <i class="fas fa-spinner fa-spin me-2"></i>
                Uploading...
              </span>
              <span v-else>
                <i class="fas fa-share me-2"></i>
                Share
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UploadModal',
  emits: ['close', 'uploaded'],
  data() {
    return {
      selectedFile: null,
      imagePreview: null,
      caption: '',
      uploading: false,
      error: null
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
    }
  },
  methods: {
    close() {
      this.$emit('close');
    },

    triggerFileInput() {
      this.$refs.fileInput.click();
    },

    handleFileSelect(event) {
      const file = event.target.files[0];
      if (file) {
        this.processFile(file);
      }
    },

    handleDrop(event) {
      const file = event.dataTransfer.files[0];
      if (file && file.type.startsWith('image/')) {
        this.processFile(file);
      }
    },

    processFile(file) {
      if (!file.type.startsWith('image/')) {
        this.error = 'Please select an image file.';
        return;
      }

      if (file.size > 10 * 1024 * 1024) {
        this.error = 'Image size must be less than 10MB.';
        return;
      }

      this.selectedFile = file;
      this.error = null;

      const reader = new FileReader();
      reader.onload = (e) => {
        this.imagePreview = e.target.result;
      };
      reader.readAsDataURL(file);
    },

    removeImage() {
      this.selectedFile = null;
      this.imagePreview = null;
      this.$refs.fileInput.value = '';
    },

    async uploadPost() {
      if (!this.selectedFile || !this.currentUserId) {
        return;
      }

      this.uploading = true;
      this.error = null;

      try {
        const formData = new FormData();
        formData.append('image', this.selectedFile);
        formData.append('caption', this.caption);

        await this.$axios.post(`/profiles/${this.currentUserId}/posts`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });

        // ✅ Emetti solo 'uploaded' senza fare refresh
        this.$emit('uploaded');
      } catch (error) {
        console.error('Upload error:', error);
        if (error.response?.status === 400) {
          this.error = 'Invalid image or caption. Please try again.';
        } else if (error.response?.status === 413) {
          this.error = 'Image file is too large.';
        } else {
          this.error = 'Failed to upload post. Please try again.';
        }
      } finally {
        this.uploading = false;
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

.upload-modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  font-weight: 600;
  margin: 0;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6c757d;
}

.modal-body {
  padding: 20px;
}

.image-upload-area {
  border: 2px dashed #dee2e6;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: #f8f9fa;
}

.image-upload-area:hover {
  border-color: #007bff;
  background-color: #f0f8ff;
}

.image-upload-area.has-image {
  padding: 0;
  border: none;
  background: none;
}

.upload-placeholder {
  color: #6c757d;
}

.image-preview {
  position: relative;
}

.image-preview img {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  border-radius: 8px;
}

.remove-image {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-control:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

@media (max-width: 576px) {
  .upload-modal {
    width: 95%;
    margin: 10px;
  }

  .modal-body {
    padding: 15px;
  }

  .image-upload-area {
    padding: 30px 15px;
  }
}
</style>