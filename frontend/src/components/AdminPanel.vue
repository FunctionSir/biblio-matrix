<template>
  <div class="admin-panel">
    <div class="tabs">
      <button v-for="tab in tabs" :key="tab.id" :class="{ active: currentTab === tab.id }" @click="switchTab(tab.id)">
        {{ tab.label }}
      </button>
    </div>

    <div class="tab-content">
      <div v-if="currentTab === 'addUser'">
        <h3>添加用户</h3>
        <form @submit.prevent="addUser">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="newUser.username" type="text" required>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="newUser.password" type="password" required>
          </div>
          <div class="form-group">
            <label>角色</label>
            <select v-model="newUser.role" required>
              <option value="admin">管理员</option>
              <option value="user">普通用户</option>
            </select>
          </div>
          <button type="submit">添加用户</button>
        </form>
      </div>

      <div v-if="currentTab === 'addBook'">
        <h3>图书入库</h3>
        <form @submit.prevent="addBook">
          <div class="form-group">
            <label>ID</label>
            <input v-model="newBook.id" type="text" required>
          </div>
          <div class="form-group">
            <label>书名</label>
            <input v-model="newBook.name" type="text">
          </div>
          <div class="form-group">
            <label>作者</label>
            <input v-model="newBook.author" type="text">
          </div>
          <div class="form-group">
            <label>数量</label>
            <input v-model="newBook.count" type="number" min="1" required>
          </div>
          <button type="submit">入库</button>
        </form>
      </div>

      <div v-if="currentTab === 'removeBook'">
        <h3>图书出库</h3>
        <form @submit.prevent="removeBook">
          <div class="form-group">
            <label>图书ID</label>
            <input v-model="bookToRemove.id" type="text" required>
          </div>
          <div class="form-group">
            <label>出库数量</label>
            <input v-model="bookToRemove.count" type="number" min="0" required>
          </div>
          <button type="submit">出库</button>
        </form>
      </div>

      <!-- 借阅记录管理 -->
      <div v-if="currentTab === 'borrowRecords'" class="borrow-records">
        <!-- 刷新按钮 -->
        <div class="header-controls">
          <button class="refresh-btn" @click="loadBorrowRecords">刷新</button>
        </div>

        <!-- 借阅记录表格 -->
        <table class="records-table">
          <thead>
            <tr>
              <th>用户名</th>
              <th>图书ID</th>
              <th>借阅日期</th>
              <th>归还日期</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in borrowRecords" :key="index">
              <td>{{ record.username }}</td>
              <td>{{ record.bookId }}</td>
              <td>{{ record.borrowDate }}</td>
              <td>{{ record.returnDate || '未归还' }}</td>
            </tr>
            <tr v-if="borrowRecords.length === 0">
              <td colspan="4" class="no-records">暂无借阅记录</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 管理员面板组件
 * 提供添加用户、图书入库、图书出库和借阅记录管理功能
 */
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

// 获取认证信息
// eslint-disable-next-line no-unused-vars
const authStore = useAuthStore()

// 当前选中的标签页
const currentTab = ref('addUser')

// 标签页配置
const tabs = [
  { id: 'addUser', label: '添加用户' },
  { id: 'addBook', label: '图书入库' },
  { id: 'removeBook', label: '图书出库' },
  { id: 'borrowRecords', label: '借阅记录' }
]

// 新用户表单数据
const newUser = ref({
  username: '',
  password: '',
  name: '',
  role: 'user'
})

// 新图书表单数据
const newBook = ref({
  name: '',
  author: '',
  id: '',
  price: 0,
  count: 1
})

// 图书出库表单数据
const bookToRemove = ref({
  id: '',
  count: 1
})

// 借阅记录相关
const borrowRecords = ref([])

// 在组件挂载时加载借阅记录
onMounted(() => {
  if (currentTab.value === 'borrowRecords') {
    loadBorrowRecords()
  }
})

/**
 * 切换标签页
 * @param {string} tabId - 标签页ID
 */
const switchTab = (tabId) => {
  currentTab.value = tabId

  // 当切换到借阅记录标签页时，加载数据
  if (tabId === 'borrowRecords') {
    loadBorrowRecords()
  }
}

/**
 * 添加用户
 * 使用真实的后端API调用
 */
const addUser = async () => {
  try {
    const formData = new FormData()
    formData.append('username', newUser.value.username)
    formData.append('passwd', newUser.value.password)
    formData.append('name', newUser.value.name || newUser.value.username)

    const endpoint = newUser.value.role === 'admin' ? '/new/admin' : '/new/reader'

    const response = await fetch(endpoint, {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (response.ok) {
      alert(`${newUser.value.role === 'admin' ? '管理员' : '用户'} ${newUser.value.username} 添加成功`)
      newUser.value = { username: '', password: '', name: '', role: 'user' }
    } else {
      const errorText = await response.text()
      alert(`添加失败: ${errorText}`)
    }
  } catch (error) {
    console.error('添加用户出错:', error)
    alert('添加用户时发生错误，请检查网络连接')
  }
}

/**
 * 添加图书
 * 使用真实的后端API调用
 */
const addBook = async () => {
  try {
    const formData = new FormData()
    formData.append('book', newBook.value.id)
    formData.append('name', newBook.value.name)
    formData.append('author', newBook.value.author)
    formData.append('price', newBook.value.price)
    formData.append('count', newBook.value.count)

    const response = await fetch('/add', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (response.ok) {
      alert(`图书 ${newBook.value.name} 入库成功`)
      newBook.value = { name: '', author: '', id: '', price: 0, count: 1 }
    } else {
      const errorText = await response.text()
      alert(`入库失败: ${errorText}`)
    }
  } catch (error) {
    console.error('添加图书出错:', error)
    alert('添加图书时发生错误，请检查网络连接')
  }
}

/**
 * 图书出库
 * 使用真实的后端API调用
 */
const removeBook = async () => {
  try {
    const formData = new FormData()
    formData.append('book', bookToRemove.value.id)
    formData.append('count', -bookToRemove.value.count)

    const response = await fetch('/add', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (response.ok) {
      alert(`图书ID ${bookToRemove.value.id} 出库 ${bookToRemove.value.count} 本成功`)
      bookToRemove.value = { id: '', count: 1 }
    } else {
      const errorText = await response.text()
      alert(`出库失败: ${errorText}`)
    }
  } catch (error) {
    console.error('图书出库出错:', error)
    alert('图书出库时发生错误，请检查网络连接')
  }
}

/**
 * 加载借阅记录
 * 使用真实的后端API调用
 */
const loadBorrowRecords = async () => {
  try {
    const response = await fetch('/list/records', {
      method: 'POST',
      credentials: 'include'
    })

    if (response.ok) {
      const records = await response.json()

      const formattedRecords = records.map(record => ({
        username: record.username,
        bookId: record.id,
        borrowDate: new Date(record.borrowed).toLocaleDateString(),
        returnDate: new Date(record.return).toLocaleDateString(),
        returned: new Date(record.return) < new Date()
      }))

      borrowRecords.value = formattedRecords
    } else {
      console.error('获取借阅记录失败:', await response.text())
      alert('获取借阅记录失败，请稍后再试')
    }
  } catch (error) {
    console.error('加载借阅记录出错:', error)
    alert('加载借阅记录时发生错误，请检查网络连接')
    useMockBorrowRecords()
  }
}

/**
 * 使用模拟数据（当API请求失败时的备用方案）
 */
const useMockBorrowRecords = () => {
  const allRecords = [
    {
      username: 'user1',
      bookId: '001',
      bookTitle: 'Vue.js设计与实现',
      borrowDate: '2023-05-15',
      returnDate: '2023-06-15',
      returned: true
    },
    {
      username: 'user2',
      bookId: '002',
      bookTitle: 'JavaScript高级程序设计',
      borrowDate: '2023-06-01',
      returnDate: null,
      returned: false
    },
    {
      username: 'user1',
      bookId: '003',
      bookTitle: '深入浅出Node.js',
      borrowDate: '2023-06-10',
      returnDate: null,
      returned: false
    }
  ]

  // 直接显示所有记录
  borrowRecords.value = allRecords
}

// markAsReturned函数已删除，因为管理员界面不再允许标记归还操作
</script>

<style scoped>
.admin-panel {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.08));
  border-radius: 28px;
  backdrop-filter: blur(25px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3),
    0 8px 30px rgba(0, 0, 0, 0.15),
    inset 0 1px 2px rgba(255, 255, 255, 0.3);
  border: 2px solid rgba(255, 255, 255, 0.25);
  position: relative;
  overflow: hidden;
}

.tabs {
  display: flex;
  margin-bottom: 3rem;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.06));
  border-radius: 18px;
  padding: 6px;
  backdrop-filter: blur(15px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1), inset 0 1px 2px rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.15);
  gap: 10px;
}

.tabs button {
  flex: 1;
  padding: 16px 32px;
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  border-radius: 14px;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  font-weight: 600;
  font-size: 1.05rem;
  position: relative;
  overflow: hidden;
  font-family: 'Poppins', sans-serif;
}

.tabs button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left 0.5s;
}

.tabs button:hover::before {
  left: 100%;
}

.tabs button:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.08));
  color: rgba(255, 255, 255, 0.95);
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.tabs button.active {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.25), rgba(255, 255, 255, 0.15)) !important;
  color: #fff !important;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2), inset 0 1px 2px rgba(255, 255, 255, 0.3) !important;
  transform: translateY(-1px) !important;
  border: 1px solid rgba(255, 255, 255, 0.3) !important;
}

.tab-content {
  padding: 20px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.04));
  border-radius: 20px;
  backdrop-filter: blur(15px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1), inset 0 1px 2px rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #c5cae9;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  max-width: 100%;
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: white;
  font-size: 14px;
  transition: all 0.3s ease;
  box-sizing: border-box;
  margin: 0;
}

.form-group input:focus,
.form-group select:focus {
  border-color: rgba(124, 77, 255, 0.6);
  box-shadow: 0 0 0 2px rgba(124, 77, 255, 0.2);
  outline: none;
}

button[type="submit"],
.refresh-btn {
  padding: 12px 24px;
  background: linear-gradient(45deg, #7c4dff, #448aff);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(124, 77, 255, 0.3);
}

button[type="submit"]:hover,
.refresh-btn:hover {
  background: linear-gradient(45deg, #651fff, #2979ff);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(124, 77, 255, 0.4);
}

/* 借阅记录样式 */
.filter-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background: rgba(0, 0, 0, 0.1);
  padding: 15px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.filter-controls .form-group {
  margin-bottom: 0;
  min-width: 200px;
}

.records-table {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

thead {
  background: rgba(124, 77, 255, 0.3);
}

th,
td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

th {
  color: white;
  font-weight: 500;
  letter-spacing: 0.5px;
}

tbody tr {
  background: rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

tbody tr:hover {
  background: rgba(0, 0, 0, 0.3);
}

tbody td {
  color: #fff;
  font-weight: 500;
}

.status-active {
  display: inline-block;
  padding: 6px 12px;
  background: rgba(76, 175, 80, 0.15);
  color: #2e7d32;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  border: 1px solid rgba(76, 175, 80, 0.2);
  backdrop-filter: blur(3px);
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.1);
}

.status-returned {
  display: inline-block;
  padding: 6px 12px;
  background: rgba(33, 150, 243, 0.15);
  color: #1565c0;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  border: 1px solid rgba(33, 150, 243, 0.2);
  backdrop-filter: blur(3px);
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.1);
}

.action-btn {
  padding: 6px 12px;
  background: linear-gradient(45deg, #a1c4fd, #c2e9fb);
  color: #1a237e;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(161, 196, 253, 0.4);
}

.no-records {
  text-align: center;
  padding: 40px;
  color: rgba(0, 0, 0, 0.6);
  font-style: italic;
  font-weight: 500;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 12px;
  margin-top: 20px;
  border: 1px dashed rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(3px);
  font-size: 1.1rem;
}

/* 已在上面定义，这里删除重复样式 */

button {
  background-color: #4CAF50;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

/* 借阅记录样式 */
.filter-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.filter-controls .form-group {
  flex: 1;
  margin-right: 15px;
  margin-bottom: 0;
}

.refresh-btn {
  padding: 8px 15px;
}

.records-table {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

table th,
table td {
  border: 1px solid #dddddd24;
  padding: 8px;
  text-align: left;
}

table th {
  background-color: #f2f2f224;
  font-weight: bold;
}

table tr:nth-child(even) {
  background-color: #f9f9f924;
}

table tr:hover {
  background-color: #f1f1f124;
}

.status-active {
  color: #e74c3c;
  font-weight: bold;
}

.status-returned {
  color: #2ecc71;
  font-weight: bold;
}

.action-btn {
  padding: 5px 10px;
  font-size: 0.9em;
  background-color: #3498db;
}

.action-btn:hover {
  background-color: #2980b9;
}

.no-records {
  text-align: center;
  padding: 20px;
  color: #7f8c8d;
  font-style: italic;
}
</style>