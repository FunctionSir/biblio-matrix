<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-title">
        <h1>瀚典云枢</h1>
        <h2 class="subtitle">图书管理系统</h2>
        <p class="slogan">知识的海洋，智慧的殿堂</p>
      </div>
      <div class="login-form-wrapper">
        <LoginForm />
      </div>
    </div>
  </div>
</template>

<script setup>
import LoginForm from '@/components/LoginForm.vue'
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 20px;
}

.login-content {
  display: flex;
  width: 100%;
  max-width: 1000px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.login-title {
  flex: 1;
  padding: 60px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  text-align: left;
}

.login-title h1 {
  font-size: 4rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  font-size: 2rem;
  color: #c5cae9;
  margin-bottom: 2rem;
  font-weight: 500;
}

.slogan {
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.7);
  font-style: italic;
  margin-top: 2rem;
}

.login-form-wrapper {
  flex: 1;
  padding: 40px;
  background: rgba(255, 255, 255, 0.05);
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (max-width: 768px) {
  .login-content {
    flex-direction: column;
  }
  
  .login-title {
    padding: 30px 20px;
    text-align: center;
  }
  
  .login-title h1 {
    font-size: 3rem;
  }
}
</style>