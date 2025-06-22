<template>
  <router-view />
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=Poppins:wght@300;400;500;600;700;800&display=swap');

* {
  box-sizing: border-box;
}

body {
  font-family: 'Inter', 'Poppins', Arial, sans-serif;
  margin: 0;
  padding: 0;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 25%, #3b82f6 50%, #6366f1 75%, #8b5cf6 100%);
  background-size: 400% 400%;
  background-attachment: fixed;
  animation: gradientShift 15s ease infinite;
  color: #f5f5f5;
  overflow-x: hidden;
  line-height: 1.6;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

body::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at 20% 80%, rgba(30, 60, 114, 0.3) 0%, transparent 50%),
              radial-gradient(circle at 80% 20%, rgba(139, 92, 246, 0.3) 0%, transparent 50%),
              radial-gradient(circle at 40% 40%, rgba(99, 102, 241, 0.2) 0%, transparent 50%);
  pointer-events: none;
  z-index: -1;
}

#app {
  min-height: 100vh;
  padding: 20px;
  position: relative;
}

/* 全局标题样式 */
h1 {
  font-size: 3rem;
  font-weight: 800;
  background: linear-gradient(135deg, #ff9a9e, #fecfef, #fecfef, #a1c4fd, #c2e9fb);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  margin-bottom: 2rem;
  text-align: center;
  letter-spacing: -0.02em;
  position: relative;
}

h1::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 4px;
  background: linear-gradient(90deg, #ff9a9e, #a1c4fd);
  border-radius: 2px;
  opacity: 0.8;
}

h2 {
  color: #ffffff;
  font-weight: 600;
  font-size: 1.8rem;
  margin-bottom: 1.5rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

h3 {
  color: #f0f4f8;
  font-weight: 500;
  font-size: 1.4rem;
  margin-bottom: 1rem;
}

/* 按钮全局样式 */
button {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.1));
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
  padding: 12px 24px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  backdrop-filter: blur(10px);
  font-weight: 500;
  font-size: 0.95rem;
  letter-spacing: 0.5px;
  position: relative;
  overflow: hidden;
}

button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

button:hover::before {
  left: 100%;
}

button:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.2));
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
  border-color: rgba(255, 255, 255, 0.5);
}

button.active {
  background: linear-gradient(135deg, rgba(124, 77, 255, 0.8), rgba(124, 77, 255, 0.6));
  border-color: rgba(124, 77, 255, 0.9);
  box-shadow: 0 6px 20px rgba(124, 77, 255, 0.4);
}

/* 表单元素全局样式 */
input, select, textarea {
  background: rgba(255, 255, 255, 0.15);
  border: 2px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 12px 16px;
  border-radius: 10px;
  transition: all 0.3s ease;
  backdrop-filter: blur(5px);
  font-size: 0.95rem;
}

input:focus, select:focus, textarea:focus {
  outline: none;
  border-color: rgba(255, 255, 255, 0.5);
  background: rgba(255, 255, 255, 0.2);
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
}

input::placeholder {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 400;
}

label {
  color: #ffffff;
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 8px;
  display: block;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.2));
  border-radius: 4px;
  backdrop-filter: blur(5px);
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4), rgba(255, 255, 255, 0.3));
}

/* 全局动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-in-up {
  animation: fadeInUp 0.6s ease-out;
}
</style>