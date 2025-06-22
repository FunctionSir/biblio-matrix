<template>
  <div class="user-panel">
    <div class="tabs">
      <button v-for="tab in tabs" :key="tab.id" @click="switchTab(tab.id)" :class="{ active: currentTab === tab.id }">
        {{ tab.label }}
      </button>
    </div>

    <div class="tab-content">
      <!-- 图书查询 -->
      <div v-if="currentTab === 'search'">
        <h3>图书查询</h3>
        <div class="search-form">
          <input v-model="searchQuery" placeholder="输入书名、作者或ID" />
          <button @click="searchBooks">搜索</button>
        </div>
        <div v-if="searchResults.length > 0" class="search-results">
          <div v-for="book in searchResults" :key="book.id" class="book-item">
            <h4>{{ book.title }}</h4>
            <p>ID: {{ book.id }}</p>
            <p>作者: {{ book.author }}</p>
            <p>有库存: {{ book.available ? "是" : "否" }}</p>
            <button @click="borrowBook(book.id)">借阅</button>
          </div>
        </div>
        <p v-else-if="searchPerformed">未找到匹配的图书</p>
      </div>

      <!-- 归还图书 -->
      <div v-if="currentTab === 'return'" class="return-book">
        <h3>归还图书</h3>
        <form @submit.prevent="returnBook" class="return-form">
          <div class="form-group">
            <label for="bookId">图书ID</label>
            <input type="text" id="bookId" v-model="recordToReturn.id" placeholder="请输入图书ID" required>
          </div>
          <button type="submit">归还</button>
        </form>
      </div>

      <!-- 我的借阅记录 -->
      <div v-if="currentTab === 'myRecords'" class="my-records">
        <h3>我的借阅记录</h3>

        <!-- 刷新按钮 -->
        <div class="header-controls">
          <button @click="loadMyBorrowRecords" class="refresh-btn">刷新</button>
        </div>

        <!-- 借阅记录表格 -->
        <table class="records-table">
          <thead>
            <tr>
              <th>图书ID</th>
              <th>借阅日期</th>
              <th>归还日期</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in myBorrowRecords" :key="index">
              <td>{{ record.bookId }}</td>
              <td>{{ record.borrowDate }}</td>
              <td>{{ record.returnDate || '未归还' }}</td>
            </tr>
            <tr v-if="myBorrowRecords.length === 0">
              <td colspan="3" class="no-records">暂无借阅记录</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 用户面板组件
 * 提供借阅图书、归还图书、图书查询和我的借阅记录功能
 */
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

// 获取认证信息
const authStore = useAuthStore()
// 获取当前用户名，用于加载该用户的借阅记录
// eslint-disable-next-line no-unused-vars
const username = authStore.username

// 当前选中的标签页
const currentTab = ref('search')

// 标签页配置
const tabs = [
  { id: 'search', label: '查询图书' },
  { id: 'return', label: '归还图书' },
  { id: 'myRecords', label: '我的借阅' }
]

// 我的借阅记录相关
const myBorrowRecords = ref([])

// 在组件挂载时加载借阅记录
onMounted(() => {
  if (currentTab.value === 'myRecords') {
    loadMyBorrowRecords()
  }
})

const bookToBorrow = ref({
  id: ''
})

const recordToReturn = ref({
  id: ''
})

const searchQuery = ref('')
const searchResults = ref([])
const searchPerformed = ref(false)

/**
 * 借阅图书
 * 使用真实的后端API调用
 */
const borrowBook = async (bookId) => {
  const idToBorrow = bookId || bookToBorrow.value.id

  try {
    const formData = new FormData()
    formData.append('username', username)
    formData.append('book', idToBorrow)

    const response = await fetch('/borrow', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (response.ok) {
      alert(`图书ID ${idToBorrow} 借阅成功`)
      bookToBorrow.value = { id: '' }
      if (currentTab.value === 'myRecords') {
        loadMyBorrowRecords()
      }
    } else {
      const errorText = await response.text()
      alert(`借阅失败: ${errorText}`)
    }
  } catch (error) {
    console.error('借阅图书出错:', error)
    alert('借阅图书时发生错误，请检查网络连接')
  }
}

/**
 * 归还图书
 * 使用真实的后端API调用
 */
const returnBook = async () => {
  try {
    const formData = new FormData()
    formData.append('username', username)
    formData.append('book', recordToReturn.value.id)

    const response = await fetch('/return', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (response.ok) {
      alert(`图书ID ${recordToReturn.value.id} 归还成功`)
      recordToReturn.value = { id: '' }
      if (currentTab.value === 'myRecords') {
        loadMyBorrowRecords()
      }
    } else {
      const errorText = await response.text()
      alert(`归还失败: ${errorText}`)
    }
  } catch (error) {
    console.error('归还图书出错:', error)
    alert('归还图书时发生错误，请检查网络连接')
  }
}

/**
 * 搜索图书
 * 使用真实的后端API调用
 */
const searchBooks = async () => {
  try {
    const response = await fetch('/list/books', {
      method: 'POST',
      credentials: 'include'
    })

    if (response.ok) {
      const books = await response.json()

      if (searchQuery.value) {
        const keyword = searchQuery.value.toLowerCase()
        searchResults.value = books.filter(book =>
          book.name.toLowerCase().includes(keyword) ||
          book.author.toLowerCase().includes(keyword) ||
          book.id.toLowerCase().includes(keyword)
        )
      } else {
        searchResults.value = books
      }

      searchResults.value = searchResults.value.map(book => ({
        id: book.id,
        title: book.name,
        author: book.author,
        available: book.count > 0
      }))

      searchPerformed.value = true
    } else {
      console.error('获取图书列表失败:', await response.text())
      alert('获取图书列表失败，请稍后再试')
      useMockSearchResults()
    }
  } catch (error) {
    console.error('搜索图书出错:', error)
    alert('搜索图书时发生错误，请检查网络连接')
    useMockSearchResults()
  }
}

/**
 * 使用模拟搜索结果（前端测试模式）
 */
const useMockSearchResults = () => {
  searchResults.value = [
    { id: '001', title: 'Vue.js设计与实现', author: '霍春阳', available: true },
    { id: '002', title: 'JavaScript高级程序设计', author: '马特·弗里斯比', available: false },
    { id: '003', title: '深入浅出Node.js', author: '朴灵', available: true },
    { id: '004', title: 'React技术揭秘', author: '卡颂', available: true },
    { id: '005', title: 'TypeScript编程', author: 'Boris Cherny', available: true },
    { id: '006', title: '算法导论', author: 'Thomas H. Cormen', available: false }
  ]

  // 标记已执行搜索
  searchPerformed.value = true
}

/**
 * 加载我的借阅记录
 * 使用真实的后端API调用
 */
const loadMyBorrowRecords = async () => {
  try {
    const formData = new FormData()
    formData.append('username', username)

    const response = await fetch('/list/records', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (response.ok) {
      const records = await response.json()

      const formattedRecords = records.map(record => ({
        bookId: record.id,
        borrowDate: new Date(record.borrowed).toLocaleDateString(),
        returnDate: new Date(record.return).toLocaleDateString(),
        returned: new Date(record.return) < new Date()
      }))

      myBorrowRecords.value = formattedRecords
    } else {
      console.error('获取借阅记录失败:', await response.text())
      alert('获取借阅记录失败，请稍后再试')
      useMockBorrowRecords()
    }
  } catch (error) {
    console.error('加载借阅记录出错:', error)
    alert('加载借阅记录时发生错误，请检查网络连接')
    useMockBorrowRecords()
  }
}

/**
 * 使用模拟借阅记录（前端测试模式）
 */
const useMockBorrowRecords = () => {
  myBorrowRecords.value = [
    {
      id: 'R001',
      bookId: '001',
      bookTitle: 'Vue.js设计与实现',
      borrowDate: '2024/01/15',
      returnDate: '未归还'
    },
    {
      id: 'R002',
      bookId: '003',
      bookTitle: '深入浅出Node.js',
      borrowDate: '2024/01/10',
      returnDate: '2024/01/20'
    },
    {
      id: 'R003',
      bookId: '005',
      bookTitle: 'TypeScript编程',
      borrowDate: '2024/01/18',
      returnDate: '未归还'
    },
    {
      id: 'R004',
      bookId: '004',
      bookTitle: 'React技术揭秘',
      borrowDate: '2024/01/05',
      returnDate: '2024/01/12'
    }
  ]
}

/**
 * 切换标签页
 * @param {string} tabId - 标签页ID
 */
const switchTab = (tabId) => {
  currentTab.value = tabId

  // 当切换到我的借阅记录标签页时，加载数据
  if (tabId === 'myRecords') {
    loadMyBorrowRecords()
  }
}
</script>

<style scoped>
.user-panel {
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

.tab {
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

.tab::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left 0.5s;
}

.tab:hover::before {
  left: 100%;
}

.tab:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.08));
  color: rgba(255, 255, 255, 0.95);
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.active-tab {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.25), rgba(255, 255, 255, 0.15)) !important;
  color: #fff !important;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2), inset 0 1px 2px rgba(255, 255, 255, 0.3) !important;
  transform: translateY(-1px) !important;
  border: 1px solid rgba(255, 255, 255, 0.3) !important;
}

.content {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.04));
  border-radius: 20px;
  padding: 2.5rem;
  backdrop-filter: blur(15px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1), inset 0 1px 2px rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
}

.tabs button {
  padding: 12px 20px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #e8eaf6;
  font-weight: 500;
  transition: all 0.3s ease;
}

.tabs button:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.tabs button.active {
  background: rgba(124, 77, 255, 0.5);
  color: white;
  border-color: rgba(124, 77, 255, 0.8);
  box-shadow: 0 4px 12px rgba(124, 77, 255, 0.3);
}

.tab-content {
  padding: 20px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.05);
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

.search-form {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.search-form input {
  flex: 1;
  max-width: 100%;
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: white;
  box-sizing: border-box;
  margin: 0;
}

.search-results {
  margin-top: 20px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 2rem;
}

.book-item {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.06));
  border-radius: 18px;
  padding: 2rem;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  position: relative;
  overflow: hidden;
}

.book-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.05), transparent);
  opacity: 0;
  transition: opacity 0.4s ease;
  pointer-events: none;
}

.book-item:hover::before {
  opacity: 1;
}

.book-item:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.18), rgba(255, 255, 255, 0.1));
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.25), 0 8px 25px rgba(0, 0, 0, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
}

.book-item h4 {
  margin: 0 0 1rem 0;
  color: #fff;
  font-size: 1.3rem;
  font-weight: 700;
  font-family: 'Poppins', sans-serif;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  line-height: 1.3;
}

.book-item p {
  margin: 5px 0;
  color: rgba(255, 255, 255, 0.7);
}

.book-item button {
  margin-top: 1.5rem;
  width: 100%;
  padding: 12px 24px;
  background: linear-gradient(135deg, rgba(124, 77, 255, 0.8), rgba(124, 77, 255, 0.6));
  border: 1px solid rgba(124, 77, 255, 0.9);
  color: white;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  font-weight: 600;
  font-size: 0.95rem;
  letter-spacing: 0.5px;
  position: relative;
  overflow: hidden;
}

.book-item button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.book-item button:hover::before {
  left: 100%;
}

.book-item button:hover {
  background: linear-gradient(135deg, rgba(124, 77, 255, 0.95), rgba(124, 77, 255, 0.8));
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 8px 25px rgba(124, 77, 255, 0.5), 0 4px 15px rgba(124, 77, 255, 0.3);
  border-color: rgba(124, 77, 255, 1);
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

.status-active {
  display: inline-block;
  padding: 4px 8px;
  background: rgba(76, 175, 80, 0.2);
  color: #4CAF50;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: 500;
  border: 1px solid rgba(76, 175, 80, 0.3);
}

.status-returned {
  display: inline-block;
  padding: 4px 8px;
  background: rgba(33, 150, 243, 0.2);
  color: #2196F3;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: 500;
  border: 1px solid rgba(33, 150, 243, 0.3);
}

.action-btn {
  padding: 6px 12px;
  background: linear-gradient(45deg, #ff9a9e, #fad0c4);
  color: #d32f2f;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: linear-gradient(45deg, #ff8a8e, #f9c0b4);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(211, 47, 47, 0.2);
}

.no-records {
  text-align: center;
  padding: 30px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.7);
  font-style: italic;
  margin-top: 20px;
  border: 1px dashed rgba(255, 255, 255, 0.2);
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

.status-active {
  color: #e74c3c;
  font-weight: bold;
}

.status-returned {
  color: #2ecc71;
  font-weight: bold;
}

.no-records {
  text-align: center;
  padding: 20px;
  color: #7f8c8d;
  font-style: italic;
}

.book-item {
  border: 1px solid #eee;
  padding: 15px;
  margin-bottom: 15px;
  border-radius: 4px;
}

.book-item h4 {
  margin-top: 0;
  color: #42b983;
}

.book-item button {
  background: #42b983;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
}

.book-item button:hover {
  background: #3aa876;
}
</style>