/**
 * 认证状态管理
 * 使用Pinia管理用户的认证状态、角色和用户名
 */
import { defineStore } from 'pinia'

/**
 * 认证状态存储
 * 提供登录和登出功能，并维护用户的认证状态
 */
export const useAuthStore = defineStore('auth', {
  // 状态定义
  state: () => ({
    isAuthenticated: false, // 是否已认证
    userRole: null,         // 用户角色（'user'或'admin'）
    username: null          // 用户名
  }),
  
  // 操作方法
  actions: {
    /**
     * 登录操作
     * 更新认证状态、用户角色和用户名
     * @param {string} username - 用户名
     * @param {string} role - 用户角色
     */
    login(username, role) {
      this.isAuthenticated = true
      this.userRole = role
      this.username = username
    },
    
    /**
     * 登出操作
     * 清除认证状态、用户角色和用户名，并调用后端登出API
     */
    async logout() {
      try {
        // 调用后端登出API
        const response = await fetch('/deauth', {
          method: 'POST',
          credentials: 'include'
        })
        
        // 无论API调用是否成功，都清除本地状态
        this.isAuthenticated = false
        this.userRole = null
        this.username = null
        
        // 如果API调用失败，记录错误但不影响用户体验
        if (!response.ok) {
          console.error('登出API调用失败:', await response.text())
        }
      } catch (error) {
        // 处理网络错误
        console.error('登出时发生错误:', error)
        
        // 仍然清除本地状态
        this.isAuthenticated = false
        this.userRole = null
        this.username = null
      }
    }
  }
})