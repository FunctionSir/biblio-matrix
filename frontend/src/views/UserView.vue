<template>
  <div class="user-view">
    <header>
      <h2>用户您好，欢迎登录</h2>
      <button @click="logout" class="logout-btn">退出系统</button>
    </header>
    
    <UserPanel />
  </div>
</template>

<script setup>
import UserPanel from '@/components/UserPanel.vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const logout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<style scoped>
.user-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 20px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

header h2 {
  font-size: 1.8rem;
  margin: 0;
  background: linear-gradient(45deg, #a1c4fd, #c2e9fb);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.logout-btn {
  padding: 10px 20px;
  background: linear-gradient(45deg, #ff9a9e, #fad0c4);
  color: #1a237e;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(255, 154, 158, 0.3);
}

.logout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(255, 154, 158, 0.4);
}

.logout-btn:hover {
  background: #d32f2f;
}
</style>