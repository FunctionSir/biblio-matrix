<template>
  <div class="login-form-container">
    <div class="login-form">
      <h2>图书管理系统</h2>
      
      <!-- 角色选择器 -->
      <div class="role-selector">
        <button 
          :class="{ active: selectedRole === 'reader' }" 
          @click="selectRole('reader')"
        >
          用户登录
        </button>
        <button 
          :class="{ active: selectedRole === 'admin' }" 
          @click="selectRole('admin')"
        >
          管理员登录
        </button>
      </div>
      
      <!-- 登录表单 -->
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="username">用户名</label>
          <input 
            type="text" 
            id="username" 
            v-model="username" 
            placeholder="请输入用户名"
            required
          >
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            placeholder="请输入密码"
            required
          >
        </div>
        
        <!-- 错误消息显示区域 - 始终存在但内容可能为空 -->
        <div class="error-container">
          <div v-if="errorMessage" class="error-message">
            {{ errorMessage }}
          </div>
          <div v-else class="error-placeholder"></div>
        </div>
        
        <button type="submit" class="login-btn">登录</button>
      </form>
    </div>
  </div>
</template>

<script setup>
/**
 * 登录表单组件
 * 提供用户和管理员登录功能，与后端API交互进行身份验证
 */
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// 获取路由和认证存储
const router = useRouter()
const authStore = useAuthStore()

// 表单数据
const selectedRole = ref('reader') // 默认选择用户角色
const username = ref('')
const password = ref('')
const errorMessage = ref('')

/**
 * 选择角色
 * @param {string} role - 角色类型（'user'或'admin'）
 */
const selectRole = (role) => {
  selectedRole.value = role
}

/**
 * 登录处理
 * 使用真实的后端API调用
 */
const login = async () => {
  // 清除之前的错误信息
  errorMessage.value = ''
  
  // 表单验证
  if (!username.value || !password.value) {
    errorMessage.value = '用户名和密码不能为空'
    return
  }
  
  try {
    const formData = new FormData()
    formData.append('username', username.value)
    formData.append('passwd', password.value)
    formData.append('role', selectedRole.value)
    
    const response = await fetch('/auth', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })
    
    if (response.ok) {
      const data = await response.json()
      const userRole = data.admin ? 'admin' : 'user'
      
      authStore.login(username.value, userRole)
      
      if (userRole === 'admin') {
        router.push('/admin')
      } else {
        router.push('/user')
      }
    } else {
      const errorText = await response.text()
      errorMessage.value = `登录失败: ${errorText}`
    }
  } catch (error) {
    errorMessage.value = '登录失败，请检查网络连接'
    console.error('登录错误:', error)
  }
}
</script>

<style scoped>
.login-form-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(20px);
  position: relative;
}

.login-form {
  width: 100%;
  max-width: 520px;
  margin: 0 auto;
  padding: 3rem;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.85));
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3), 
              0 8px 30px rgba(0, 0, 0, 0.15),
              inset 0 1px 2px rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  min-height: 680px;
  border: 2px solid rgba(255, 255, 255, 0.4);
  animation: formAppear 0.8s ease-out forwards;
  box-sizing: border-box;
  overflow: hidden;
  position: relative;
}

@keyframes formAppear {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-form h2 {
  text-align: center;
  margin-bottom: 3.5rem;
  background: linear-gradient(135deg, #667eea, #764ba2, #f093fb);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 800;
  font-size: 2.8rem;
  text-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  letter-spacing: -0.02em;
  position: relative;
  font-family: 'Poppins', sans-serif;
}

.login-form h2::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 70px;
  height: 3px;
  background: linear-gradient(90deg, #4a6fa5, #2c3e50);
  border-radius: 3px;
}

.role-selector {
  display: flex;
  margin-bottom: 3.5rem;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  background: rgba(240, 240, 240, 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.role-selector button {
  flex: 1;
  padding: 1.2rem;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  font-weight: 600;
  font-size: 1rem;
  color: #555;
  position: relative;
  overflow: hidden;
}

.role-selector button::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  opacity: 0;
  transition: opacity 0.4s ease;
  z-index: -1;
}

.role-selector button.active {
  color: white;
  transform: scale(1.02);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.role-selector button.active::before {
  opacity: 1;
}

.role-selector button:hover:not(.active) {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.form-group {
  margin-bottom: 2.5rem; /* 增加表单组间距使布局更均匀 */
}

.form-group label {
  display: block;
  margin-bottom: 0.75rem;
  font-weight: 600;
  color: #333;
  font-size: 1.1rem;
}

.form-group input {
  width: 100%;
  max-width: 100%;
  padding: 1rem 1.2rem;
  border: 2px solid rgba(200, 200, 200, 0.6);
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.9));
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08), inset 0 1px 2px rgba(255, 255, 255, 0.8);
  color: #2c3e50;
  box-sizing: border-box;
  margin: 0;
  backdrop-filter: blur(5px);
  font-family: 'Inter', sans-serif;
}

.form-group input::placeholder {
  color: #a0aec0;
  transition: all 0.3s ease;
}

.form-group input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.2), 0 8px 25px rgba(102, 126, 234, 0.15);
  outline: none;
  transform: translateY(-3px) scale(1.01);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.98), rgba(255, 255, 255, 0.95));
}

.form-group input:focus::placeholder {
  opacity: 0.7;
  transform: translateX(5px);
}

/* 错误消息容器 */
.error-container {
  height: 70px; /* 进一步增加固定高度 */
  margin-bottom: 2.5rem; /* 增加与登录按钮的间距使布局更均匀 */
  display: flex;
  align-items: center;
  position: relative;
}

.error-message {
  width: 100%;
  color: #e74c3c;
  background-color: #fdeaea;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-size: 1rem;
  border-left: 4px solid #e74c3c;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(231, 76, 60, 0.2);
  animation: fadeIn 0.3s ease-in-out;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  z-index: 10; /* 确保错误消息显示在最上层 */
}

.error-placeholder {
  width: 100%;
  height: 100%;
  visibility: hidden;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 表单样式调整 */
form {
  position: relative;
}

.login-btn {
  width: 100%;
  padding: 1.4rem;
  margin-top: 2.5rem;
  background: linear-gradient(135deg, #667eea, #764ba2, #f093fb);
  color: white;
  border: none;
  border-radius: 16px;
  font-size: 1.1rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4), 
              0 4px 15px rgba(118, 75, 162, 0.3),
              inset 0 1px 2px rgba(255, 255, 255, 0.3);
  letter-spacing: 0.5px;
  position: relative;
  overflow: hidden;
  text-transform: uppercase;
  z-index: 5;
  font-family: 'Poppins', sans-serif;
}

.login-btn:hover {
  background: linear-gradient(135deg, #7c8cfa, #8a5fc7, #ff6bcb);
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 35px rgba(102, 126, 234, 0.6), 
              0 8px 25px rgba(118, 75, 162, 0.4),
              inset 0 1px 3px rgba(255, 255, 255, 0.4);
}

.login-btn:active {
  transform: translateY(-1px) scale(0.98);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
  background: linear-gradient(135deg, #5a6fd8, #6a4f9a, #e055a8);
}

/* 添加按钮点击波纹效果 */
.login-btn::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 5px;
  height: 5px;
  background: rgba(255, 255, 255, 0.5);
  opacity: 0;
  border-radius: 100%;
  transform: scale(1, 1) translate(-50%);
  transform-origin: 50% 50%;
}

.login-btn:focus:not(:active)::after {
  animation: ripple 1s ease-out;
}

@keyframes ripple {
  0% {
    transform: scale(0, 0);
    opacity: 0.5;
  }
  100% {
    transform: scale(100, 100);
    opacity: 0;
  }
}
</style>