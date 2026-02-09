<template>
  <div class="series-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" class="search-form" @submit.prevent>
        <el-form-item label="名称">
          <el-input v-model="searchForm.name" placeholder="请输入名称" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择" clearable style="width: 120px" :teleported="false">
            <el-option label="全部" value="" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item class="search-actions">
          <HaloButton type="primary" size="small" content="搜索" :icon="Search" @click.prevent="handleSearch" />
          <HaloButton type="default" size="small" content="重置" :icon="Refresh" @click.prevent="handleReset" />
        </el-form-item>
      </el-form>

      <!-- 操作栏 -->
      <div class="action-bar">
        <HaloButton type="primary" size="medium" content="新增系列" :icon="Plus" @click="handleCreate" />
      </div>

      <!-- 卡片列表 -->
      <div v-loading="loading" class="series-grid">
        <div v-for="series in seriesList" :key="series.id" class="series-card">
          <div class="card-cover">
            <img v-if="series.cover" :src="series.cover" :alt="series.name" />
            <div v-else class="default-cover">{{ series.icon }}</div>
          </div>
          <div class="card-content">
            <h3>{{ series.name }}</h3>
            <p class="description">{{ series.description }}</p>
            <div class="meta">
              <el-tag :type="series.status === 1 ? 'success' : 'info'" size="small">
                {{ series.status === 1 ? '启用' : '禁用' }}
              </el-tag>
              <span class="info">
                <el-icon>
                  <Document />
                </el-icon>
                {{ series.sections?.length || 0 }} 章节
              </span>
            </div>
            <div class="actions">
              <el-button type="primary" link size="small" @click="handleView(series)">查看</el-button>
              <el-button type="primary" link size="small" @click="handleEdit(series)">编辑</el-button>
              <el-button type="danger" link size="small" @click="handleDelete(series.id)">删除</el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
          :page-sizes="[12, 24, 48, 96]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadList" @current-change="loadList" background />
      </div>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="系列名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入系列名称" />
        </el-form-item>
        <el-form-item label="系列标识" prop="slug">
          <el-input v-model="form.slug" placeholder="请输入系列标识(URL别名)" />
        </el-form-item>
        <el-form-item label="系列图标" prop="icon">
          <el-input v-model="form.icon" placeholder="请输入图标（emoji）" />
        </el-form-item>
        <el-form-item label="封面图片" prop="cover">
          <el-input v-model="form.cover" placeholder="请输入封面图片URL" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sortOrder">
          <el-input-number v-model="form.sortOrder" :min="0" placeholder="请输入排序" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="dialogTitle" width="1200px" class="detail-dialog" @close="closeDetail">
      <div v-if="currentSeries" class="series-detail">
        <!-- 顶部信息区域 -->
        <div class="detail-header">
          <div class="cover-section">
            <div class="cover-image">
              <img v-if="currentSeries.cover" :src="currentSeries.cover" :alt="currentSeries.name" />
              <div v-else class="cover-placeholder">
                <span class="placeholder-icon">📚</span>
              </div>
            </div>
          </div>

          <div class="info-section">
            <div class="title-row">
              <h2>{{ currentSeries.name }}</h2>
              <div class="tags">
                <el-tag :type="currentSeries.status === 1 ? 'success' : 'info'" size="small">
                  {{ currentSeries.status === 1 ? '已启用' : '已禁用' }}
                </el-tag>
              </div>
              <div class="stats">
                <span class="stat">{{ currentSeries.sections?.length || 0 }} 章节</span>
                <span class="stat">{{ getTotalArticles() }} 文章</span>
              </div>
              <div class="actions">
                <HaloButton type="default" size="small" content="" :icon="Edit" @click="editSeriesDetail" />
              </div>
            </div>
            <p class="description">{{ currentSeries.description || '暂无描述' }}</p>
            <div class="bottom-row">
              <span class="time">{{ formatDate(currentSeries.created_at) }}</span>
              <div class="actions">
                <HaloButton type="primary" size="small" content="添加章节" :icon="Plus" @click="openSectionDialog()" />
                <HaloButton type="default" size="small" content="导出" :icon="Download" @click="exportSeries" />
              </div>
            </div>
          </div>
        </div>

        <!-- 章节区域 -->
        <div class="chapters-section">
          <div class="section-header">
            <h3>章节结构</h3>
            <div class="section-actions">
              <el-button text @click="expandAll">展开全部</el-button>
              <el-button text @click="collapseAll">折叠全部</el-button>
            </div>
          </div>

          <div class="chapters-container">
            <!-- 章节列表 -->
            <div v-if="currentSeries.sections && currentSeries.sections.length > 0" class="chapters-list">
              <div v-for="(section, sectionIndex) in currentSeries.sections" :key="section.id" class="chapter-item"
                :class="{ 'expanded': expandedSections.includes(section.id) }">
                <!-- 章节头部 -->
                <div class="chapter-header" @click="toggleSection(section.id)">
                  <div class="chapter-left">
                    <span class="chapter-num">第{{ sectionIndex + 1 }}章</span>
                    <div class="chapter-info">
                      <h4>{{ section.name }}</h4>
                      <span class="chapter-meta">{{ section.subchapters?.length || 0 }} 个子章节</span>
                    </div>
                  </div>
                  <div class="chapter-right">
                    <el-button text size="small" @click.stop="openSubChapterDialog(section)">
                      <el-icon><Plus /></el-icon>
                    </el-button>
                    <el-button text size="small" @click.stop="editSection(section)">
                      <el-icon><Edit /></el-icon>
                    </el-button>
                    <el-button text size="small" type="danger" @click.stop="deleteSectionOperate(section.id)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                    <el-icon class="expand-icon">
                      <ArrowDown v-if="expandedSections.includes(section.id)" />
                      <ArrowRight v-else />
                    </el-icon>
                  </div>
                </div>

                <!-- 子章节列表 -->
                <div v-if="expandedSections.includes(section.id)" class="subchapters-list">
                  <div v-for="(sub, subIndex) in section.subchapters" :key="sub.id" class="subchapter-item">
                    <div class="subchapter-left">
                      <span class="subchapter-num">{{ sectionIndex + 1 }}.{{ subIndex + 1 }}</span>
                      <div class="subchapter-info">
                        <h5>{{ sub.name }}</h5>
                        <span class="subchapter-meta">{{ sub.articleIds?.length || 0 }} 篇文章</span>
                      </div>
                    </div>
                    <el-dropdown @command="handleSubChapterCommand($event, section, sub)">
                      <el-button text size="small">
                        <el-icon><Operation /></el-icon>
                      </el-button>
                      <template #dropdown>
                          <el-dropdown-menu>
                            <el-dropdown-item command="addArticle">
                              <el-icon>
                                <Plus />
                              </el-icon>
                              添加文章
                            </el-dropdown-item>
                            <el-dropdown-item command="viewArticles">
                              <el-icon>
                                <View />
                              </el-icon>
                              查看文章
                            </el-dropdown-item>
                            <el-dropdown-item command="edit" divided>
                              <el-icon>
                                <Edit />
                              </el-icon>
                              编辑
                            </el-dropdown-item>
                            <el-dropdown-item command="delete">
                              <el-icon>
                                <Delete />
                              </el-icon>
                              删除
                            </el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                      </el-dropdown>
                  </div>
                </div>

                <!-- 空状态 -->
                <div v-if="!section.subchapters || section.subchapters.length === 0" class="empty-hint">
                  暂无子章节
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-else class="empty-state">
              <el-icon class="empty-icon"><Folder /></el-icon>
              <p>暂无章节</p>
              <el-button type="primary" size="small" @click="openSectionDialog()">
                <el-icon><Plus /></el-icon>
                创建章节
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 章节/子章节对话框 -->
    <el-dialog v-model="chapterDialogVisible" :title="chapterDialogTitle" width="500px">
      <el-form :model="chapterForm" :rules="chapterRules" ref="chapterFormRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="chapterForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="chapterForm.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sortOrder">
          <el-input-number v-model="chapterForm.sortOrder" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="chapterDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleChapterSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 文章选择对话框 -->
    <el-dialog v-model="articleDialogVisible" title="选择文章" width="700px">
      <div class="article-selector">
        <el-form :model="articleSearch" inline class="article-search">
          <el-form-item label="搜索">
            <el-input v-model="articleSearch.title" placeholder="输入文章标题" clearable style="width: 200px"
              @clear="loadArticleList" />
            <el-button type="primary" @click="loadArticleList">搜索</el-button>
          </el-form-item>
        </el-form>
        <el-checkbox-group v-model="selectedArticles" v-loading="articleLoading">
          <div v-for="article in articleList" :key="article.id" class="article-item">
            <el-checkbox :label="article.id">
              <span class="article-title">{{ article.title }}</span>
              <span class="article-date">{{ article.created_at }}</span>
            </el-checkbox>
          </div>
        </el-checkbox-group>
        <div v-if="articleList.length === 0" class="empty-tip">
          <el-empty description="暂无可用文章" />
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <div class="selected-info">
            <el-tag type="info" size="small">已选 {{ selectedArticles.length }} 篇</el-tag>
          </div>
          <div class="footer-buttons">
            <el-button @click="articleDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleArticleSubmit">确定</el-button>
          </div>
        </div>
        <!-- 文章列表分页器 -->
        <div class="article-pagination-wrapper">
          <el-pagination v-model:current-page="articlePagination.page" v-model:page-size="articlePagination.pageSize"
            :page-sizes="[10, 20, 50]" :total="articlePagination.total" layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadArticleList" @current-change="loadArticleList" background small />
        </div>
      </template>
    </el-dialog>

    <!-- 查看文章列表对话框 -->
    <el-dialog v-model="articleListDialogVisible" title="已选文章" width="600px">
      <el-table :data="currentArticles" border>
        <el-table-column prop="title" label="文章标题" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="100">
          <template #default="{ row, $index }">
            <el-button type="danger" link size="small" @click="removeArticle($index)">移除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({
  name: 'SeriesList'
})

import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Document, Folder, Edit, Delete, View,
  Camera, Sort, Calendar, Download, Expand, Fold,
  ArrowDown, ArrowRight, Operation, Search, Refresh
} from '@element-plus/icons-vue'
import HaloButton from '@/layout/components/HaloButton.vue'
import {
  getSeriesList,
  getSeriesDetail,
  createSeries,
  updateSeries,
  deleteSeries,
  createSection,
  updateSection,
  deleteSection,

  createSubchapter,
  updateSubchapter,
  deleteSubchapter,
  getSubchapterArticles,
  addArticleToSubchapter
} from '@/api/series'


import { getArticleList } from '@/api/article'

const loading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const seriesList = ref([])
const currentSeries = ref(null)
const chapterDialogVisible = ref(false)
const chapterFormRef = ref(null)
const articleDialogVisible = ref(false)
const articleListDialogVisible = ref(false)
const articleLoading = ref(false)

const searchForm = reactive({
  name: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})
// 新增的响应式数据
const expandedSections = ref([])

// 新增的计算属性
const dialogTitle = computed(() => {
  return currentSeries.value ? `${currentSeries.value.name} - 系列管理` : '系列管理'
})

// 新增的方法
const getTotalArticles = () => {
  if (!currentSeries.value?.sections) return 0
  let total = 0
  currentSeries.value.sections.forEach(section => {
    if (section.subchapters) {
      section.subchapters.forEach(sub => {
        total += sub.articleIds?.length || 0
      })
    }
  })
  return total
}

// const dialogTitle = ref('新增系列')

const formatDate = (dateStr) => {
  if (!dateStr) return '未记录'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const toggleSection = (sectionId) => {
  const index = expandedSections.value.indexOf(sectionId)
  if (index > -1) {
    expandedSections.value.splice(index, 1)
  } else {
    expandedSections.value.push(sectionId)
  }
}

const expandAll = () => {
  if (!currentSeries.value?.sections) return
  expandedSections.value = currentSeries.value.sections.map(section => section.id)
}

const collapseAll = () => {
  expandedSections.value = []
}

const handleSubChapterCommand = (command, section, subchapter) => {
  switch (command) {
    case 'addArticle':
      openArticleDialog(section, subchapter)
      break
    case 'viewArticles':
      viewArticles(section, subchapter)
      break
    case 'edit':
      editSubChapter(section, subchapter)
      break
    case 'delete':
      deleteSubChapter(section, subchapter)
      break
  }
}

const editSeriesDetail = () => {
  handleEdit(currentSeries.value)
}

const exportSeries = () => {
  // 导出系列结构的实现
  ElMessage.info('导出功能开发中...')
}

// 监听详情对话框打开
watch(detailVisible, (newVal) => {
  if (newVal && currentSeries.value?.sections) {
    // 默认展开所有章节
    expandedSections.value = currentSeries.value.sections.map(section => section.id)
  } else {
    expandedSections.value = []
  }
})

const form = reactive({
  id: null,
  name: '',
  slug: '',
  icon: '📚',
  cover: '',
  description: '',
  sortOrder: 0,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入系列名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入系列标识', trigger: 'blur' }]
}

// 章节/子章节相关
const chapterDialogTitle = ref('')
const chapterForm = reactive({
  id: null,
  name: '',
  sortOrder: 0,
  description: '',
  sectionId: null,
  seriesId: null,
  type: 'section'
})

const chapterRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

// 文章选择相关
const currentSection = ref(null)
const currentSubChapter = ref(null)
const articleList = ref([])
const selectedArticles = ref([])
const articleSearch = reactive({
  title: ''
})
const currentArticles = ref([])

// 文章分页
const articlePagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadList = async () => {
  loading.value = true
  try {
    const res = await getSeriesList({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    seriesList.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadList()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.status = ''
  handleSearch()
}

const handleCreate = () => {
  dialogTitle.value = '新增系列'
  Object.assign(form, {
    id: null,
    name: '',
    slug: '',
    icon: '📚',
    cover: '',
    description: '',
    sortOrder: 0,
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑系列'
  Object.assign(form, {
    id: row.id,
    name: row.name,
    slug: row.slug,
    icon: row.icon || '📚',
    cover: row.cover || '',
    description: row.description || '',
    sortOrder: row.sortOrder || 0,
    status: row.status || 1
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const dataToSubmit = {
        name: form.name,
        slug: form.slug,
        icon: form.icon,
        description: form.description,
        cover: form.cover,
        sortOrder: form.sortOrder,
        status: form.status
      }
      console.log('提交的系列数据:', dataToSubmit)

      if (form.id) {
        await updateSeries(form.id, dataToSubmit)
        ElMessage.success('更新成功')
      } else {
        await createSeries(dataToSubmit)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadList()
    } catch (error) {
      console.error('操作失败:', error)
      ElMessage.error('操作失败: ' + (error.message || error.response?.data?.message || error))
    } finally {
      submitting.value = false
    }
  })
}

const handleView = async (row) => {
  try {
    const res = await getSeriesDetail(row.id)
    currentSeries.value = res.data
    detailVisible.value = true
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该系列吗？', '提示', { type: 'warning' })
    await deleteSeries(id)
    ElMessage.success('删除成功')
    loadList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除系列失败:', error)
      ElMessage.error('删除失败: ' + (error.message || error.response?.data?.message || error))
    }
  }
}

// 章节管理
const openSectionDialog = (section = null) => {
  chapterDialogTitle.value = section ? '编辑章节' : '添加章节'
  Object.assign(chapterForm, {
    id: section?.id || null,
    name: section?.name || '',
    sortOrder: section?.sortOrder || 0,
    description: section?.description || '',
    sectionId: null,
    seriesId: currentSeries.value?.id || null,
    type: 'section'
  })
  chapterDialogVisible.value = true
}

const openSubChapterDialog = (section) => {
  chapterDialogTitle.value = '添加子章节'
  Object.assign(chapterForm, {
    id: null,
    name: '',
    sortOrder: 0,
    description: '',
    sectionId: section.id,
    type: 'subchapter'
  })
  chapterDialogVisible.value = true
}

const editSection = (section) => {
  openSectionDialog(section)
}

const editSubChapter = (section, subchapter) => {
  chapterDialogTitle.value = '编辑子章节'
  Object.assign(chapterForm, {
    id: subchapter.id,
    name: subchapter.name,
    sortOrder: subchapter.sortOrder || 0,
    description: subchapter.description || '',
    sectionId: section.id,
    type: 'subchapter'
  })
  chapterDialogVisible.value = true
}

const handleChapterSubmit = async () => {
  if (!chapterFormRef.value) return

  await chapterFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (chapterForm.type === 'section') {
        // 构建章节请求数据
        const sectionData = {
          name: chapterForm.name,
          description: chapterForm.description,
          sortOrder: chapterForm.sortOrder || 0
        }
        if (chapterForm.id) {
          await updateSection(chapterForm.id, sectionData)
          ElMessage.success('章节更新成功')
        } else {
          await createSection(currentSeries.value.id, sectionData)
          ElMessage.success('章节添加成功')
        }
      } else {
        // 构建子章节数据
        const subchapterData = {
          name: chapterForm.name,
          description: chapterForm.description,
          sortOrder: chapterForm.sortOrder || 0,
          sectionId: chapterForm.sectionId
        }
        if (chapterForm.id) {
          await updateSubchapter(chapterForm.id, subchapterData)
          ElMessage.success('子章节更新成功')
        } else {
          await createSubchapter(subchapterData)
          ElMessage.success('子章节添加成功')
        }
      }
      chapterDialogVisible.value = false
      await handleView(currentSeries.value)
    } catch (error) {
      console.error('操作失败:', error)
      ElMessage.error('操作失败: ' + (error.message || error.response?.data?.message || error))
    }
  })
}

const deleteSectionOperate = async (sectionId) => {
  try {
    await ElMessageBox.confirm('确定删除该章节及其所有子章节吗？', '提示', { type: 'warning' })
    await deleteSection(sectionId)
    ElMessage.success('删除成功')
    await handleView(currentSeries.value)
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const deleteSubChapter = async (section, subchapter) => {
  try {
    await ElMessageBox.confirm('确定删除该子章节吗？', '提示', { type: 'warning' })
    await deleteSubchapter(subchapter.id)
    ElMessage.success('删除成功')
    await handleView(currentSeries.value)
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 文章管理
const openArticleDialog = async (section, subchapter) => {
  currentSection.value = section
  currentSubChapter.value = subchapter
  selectedArticles.value = []
  articlePagination.page = 1
  await loadArticleList()
  articleDialogVisible.value = true
}

const loadArticleList = async () => {
  articleLoading.value = true
  try {
    const res = await getArticleList({
      page: articlePagination.page,
      pageSize: articlePagination.pageSize,
      keyword: articleSearch.title,
      status: 1
    })
    articleList.value = res.data.items || []
    articlePagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载文章列表失败')
  } finally {
    articleLoading.value = false
  }
}

const handleArticleSubmit = async () => {
  try {
    for (const articleId of selectedArticles.value) {
      await addArticleToSubchapter(currentSubChapter.value.id, {
        articleId: articleId,
        sortOrder: 0
      })
    }
    ElMessage.success(`已添加 ${selectedArticles.value.length} 篇文章`)
    articleDialogVisible.value = false
    await handleView(currentSeries.value)
  } catch (error) {
    console.error('添加文章失败:', error)
    ElMessage.error('添加文章失败: ' + (error.message || error.response?.data?.message || error))
  }
}

const viewArticles = async (section, subchapter) => {
  currentSection.value = section
  currentSubChapter.value = subchapter
  try {
    const res = await getSubchapterArticles(subchapter.id)
    currentArticles.value = res.data || []
    articleListDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取文章列表失败')
  }
}

const removeArticle = (index) => {
  currentArticles.value.splice(index, 1)
  ElMessage.success('移除成功')
}

const closeDetail = () => {
  currentSeries.value = null
}

onMounted(() => {
  loadList()
})
</script>

<style scoped lang="scss">
.series-list {
  height: 100%;

  :deep(.el-card) {
    height: 100%;
    display: flex;
    flex-direction: column;

    .el-card__body {
      height: 100%;
      display: flex;
      flex-direction: column;
      padding: 16px;
    }
  }

  .search-form {
    margin-bottom: 12px;
    display: flex;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 12px;
    flex-shrink: 0;

    :deep(.el-form-item) {
      margin-bottom: 0;
    }

    .search-actions {
      margin-left: auto;

      :deep(.el-form-item__content) {
        display: flex;
      }

      .HaloButton {
        margin-left: 8px;
      }
    }
  }

  .action-bar {
    margin-bottom: 12px;
    flex-shrink: 0;
  }

  .series-grid {
    flex: 1;
    overflow-y: auto;
  }

  .series-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;

    .series-card {
      background: #fff;
      border: 1px solid #e1f5fe;
      border-radius: 12px;
      overflow: hidden;
      transition: all 0.3s ease;
      box-shadow: 0 2px 8px rgba(33, 150, 243, 0.08);

      &:hover {
        box-shadow: 0 8px 24px rgba(33, 150, 243, 0.15);
        transform: translateY(-4px);
        border-color: #64b5f6;
      }

      .card-cover {
        height: 140px;
        overflow: hidden;
        background: linear-gradient(135deg, #81d4fa 0%, #4fc3f7 100%);
        display: flex;
        align-items: center;
        justify-content: center;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }

        .default-cover {
          font-size: 48px;
          color: white;
          text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }
      }

      .card-content {
        padding: 16px;

        h3 {
          margin: 0 0 8px;
          font-size: 16px;
          color: #0288d1;
          font-weight: 600;
        }

        .description {
          margin: 0 0 12px;
          font-size: 13px;
          color: #546e7a;
          line-height: 1.5;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
          min-height: 39px;
        }

        .meta {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 12px;

          .info {
            display: flex;
            align-items: center;
            gap: 4px;
            color: #78909c;
            font-size: 12px;
          }
        }

        .actions {
          display: flex;
          gap: 8px;

          .el-button {
            font-size: 12px;
          }
        }
      }
    }
  }

  .series-detail {
    /* 顶部信息区域 - 简约版 */
    .detail-header {
      display: flex;
      gap: 16px;
      padding: 16px;
      background: #fafafa;
      border-radius: 6px;
      margin-bottom: 16px;
      border: 1px solid #f0f0f0;

      .cover-section {
        flex-shrink: 0;

        .cover-image {
          width: 80px;
          height: 100px;
          border-radius: 4px;
          overflow: hidden;
          background: #f0f0f0;
          box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);

          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }
        }

        .cover-placeholder {
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
          background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ec 100%);

          .placeholder-icon {
            font-size: 28px;
            opacity: 0.5;
          }
        }
      }

      .info-section {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 8px;
        min-width: 0;

        .title-row {
          display: flex;
          align-items: center;
          gap: 12px;
          flex-wrap: wrap;

          h2 {
            margin: 0;
            font-size: 16px;
            font-weight: 600;
            color: #333;
            flex-shrink: 0;
          }

          .tags {
            flex-shrink: 0;
          }

          .stats {
            display: flex;
            gap: 12px;
            font-size: 12px;
            color: #666;

            .stat {
              display: flex;
              align-items: center;
            }
          }

          .actions {
            margin-left: auto;
          }
        }

        .description {
          margin: 0;
          color: #999;
          line-height: 1.4;
          font-size: 12px;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          -webkit-box-lines: 2;
          line-clamp: 2;
        }

        .bottom-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-top: 4px;

          .time {
            font-size: 12px;
            color: #999;
          }

          .actions {
            display: flex;
            gap: 8px;
          }
        }
      }
    }



    .chapters-section {
      background: #fff;
      border-radius: 12px;
      box-shadow: 0 4px 20px rgba(33, 150, 243, 0.08);
      overflow: hidden;
      border: 1px solid #e1f5fe;

      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 20px 24px;
        background: linear-gradient(to right, #e1f5fe, #f5fbfe);
        border-bottom: 2px solid #4fc3f7;

        .section-title {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: #0277bd;
        }

        .section-actions {
          display: flex;
          gap: 8px;

          .el-button {
            color: #0288d1;
            border-color: #4fc3f7;

            &:hover {
              background-color: #e1f5fe;
            }
          }
        }
      }

      .chapters-container {
        padding: 0;

        .chapters-list {
          .chapter-item {
            margin-bottom: 8px;
            border: 1px solid #e8e8e8;
            background: #fff;
            transition: all 0.2s;

            &:hover {
              border-color: #409eff;
            }

            &.expanded {
              background: #f5f7fa;
            }

            .chapter-header {
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding: 16px 20px;
              cursor: pointer;
              background: #fff;

              .chapter-left {
                display: flex;
                align-items: center;
                gap: 16px;
                flex: 1;

                .chapter-num {
                  flex-shrink: 0;
                  font-size: 13px;
                  font-weight: 500;
                  color: #909399;
                  padding: 4px 10px;
                  background: #f0f0f0;
                  border-radius: 4px;
                }

                .chapter-info {
                  flex: 1;

                  h4 {
                    margin: 0 0 4px 0;
                    font-size: 15px;
                    font-weight: 500;
                    color: #303133;
                  }

                  .chapter-desc {
                    margin: 0;
                    font-size: 13px;
                    color: #546e7a;
                    display: -webkit-box;
                    -webkit-line-clamp: 1;
                    -webkit-box-orient: vertical;
                    overflow: hidden;
                  }
                }

                .chapter-meta {
                  display: flex;
                  align-items: center;
                  gap: 12px;
                  flex-shrink: 0;

                  .sort-order {
                    font-size: 12px;
                    color: #78909c;
                  }
                }
              }

              .chapter-actions {
                display: flex;
                align-items: center;
                gap: 8px;

                .el-button {
                  color: #0288d1;

                  &--primary {
                    color: #0277bd;

                    &:hover {
                      background-color: #e1f5fe;
                    }
                  }

                  &--danger {
                    color: #f44336;
                  }
                }

                .expand-icon {
                  font-size: 16px;
                  color: #4fc3f7;
                  transition: transform 0.3s ease;
                  margin-left: 8px;
                }
              }
            }

            .subchapters-list {
              padding: 16px 20px 16px 116px;
              background: #fafdfe;

              .subchapter-item {
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 12px 16px;
                background: #fff;
                border-radius: 8px;
                border: 1px solid #e1f5fe;
                margin-top: 8px;
                transition: all 0.2s ease;

                &:hover {
                  border-color: #81d4fa;
                  background: #f8fcff;
                  box-shadow: 0 2px 8px rgba(129, 212, 250, 0.2);
                }

                .subchapter-content {
                  display: flex;
                  align-items: flex-start;
                  gap: 16px;
                  flex: 1;

                  .subchapter-number {
                    flex-shrink: 0;
                    width: 40px;
                    text-align: center;
                    font-size: 12px;
                    font-weight: 600;
                    color: white;
                    padding: 4px 6px;
                    background: linear-gradient(135deg, #81d4fa, #4fc3f7);
                    border-radius: 4px;
                    margin-top: 2px;
                    box-shadow: 0 1px 3px rgba(129, 212, 250, 0.4);
                  }

                  .subchapter-main {
                    flex: 1;

                    .subchapter-header {
                      display: flex;
                      align-items: center;
                      gap: 12px;
                      margin-bottom: 4px;

                      h5 {
                        margin: 0;
                        font-size: 14px;
                        font-weight: 600;
                        color: #0277bd;
                      }

                      .article-count {
                        flex-shrink: 0;
                        background: #e8f5e9;
                        border-color: #81c784;
                        color: #388e3c;
                      }
                    }

                    .subchapter-desc {
                      margin: 0;
                      font-size: 13px;
                      color: #546e7a;
                      line-height: 1.4;
                    }
                  }
                }

                .subchapter-actions {
                  flex-shrink: 0;

                  .el-button {
                    color: #0288d1;

                    &:hover {
                      background-color: #e1f5fe;
                    }
                  }
                }
              }

              .empty-subchapters {
                margin-top: 16px;
                background: #f8fcff;
                border-radius: 8px;
                border: 1px dashed #bbdefb;
                padding: 20px;

                :deep(.el-empty__description) {
                  margin-top: 8px;
                  color: #78909c;
                }

                .el-button {
                  background: linear-gradient(135deg, #81d4fa, #4fc3f7);
                  border: none;
                  color: white;

                  &:hover {
                    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
                  }
                }
              }
            }
          }
        }

        .empty-chapters {
          padding: 60px 0;
          background: #f8fcff;
          border-radius: 10px;
          border: 2px dashed #bbdefb;

          :deep(.el-empty__image) {
            opacity: 0.6;
          }

          .el-button {
            background: linear-gradient(135deg, #4fc3f7, #29b6f6);
            border: none;
            color: white;

            &:hover {
              background: linear-gradient(135deg, #29b6f6, #0288d1);
            }
          }
        }
      }
    }
  }

  .article-selector {
    max-height: 400px;
    overflow-y: auto;
    margin-bottom: 16px;

    .article-search {
      margin-bottom: 16px;
      padding-bottom: 16px;
      border-bottom: 1px solid #e1f5fe;

      .el-input {
        :deep(.el-input__wrapper) {
          border-color: #4fc3f7;

          &:hover {
            border-color: #29b6f6;
          }
        }
      }
    }

    .article-item {
      padding: 10px;
      border-bottom: 1px solid #f0f8ff;
      display: flex;
      align-items: center;
      transition: all 0.2s;

      &:hover {
        background: #f8fcff;
        border-radius: 6px;
      }

      .article-title {
        flex: 1;
        margin-left: 10px;
        color: #37474f;
      }

      .article-date {
        color: #78909c;
        font-size: 12px;
      }
    }

    .empty-tip {
      text-align: center;
      padding: 40px 0;
      color: #78909c;
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;

    .selected-info {
      display: flex;
      align-items: center;

      .el-tag {
        background: #e3f2fd;
        border-color: #90caf9;
        color: #0277bd;
      }
    }

    .footer-buttons {
      display: flex;
      gap: 8px;
    }
  }

  .article-pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #e1f5fe;

    :deep(.el-pagination) {
      .el-pager li {
        border-radius: 4px;
        font-weight: 500;
        margin: 0 4px;

        &.is-active {
          background: linear-gradient(135deg, #4fc3f7, #29b6f6);
        }

        &:hover {
          color: #0288d1;
        }
      }

      .btn-prev,
      .btn-next {
        border-radius: 4px;
      }

      button:hover {
        color: #0288d1;
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid #e1f5fe;

    :deep(.el-pagination) {
      .el-pager li {
        border-radius: 4px;
        font-weight: 500;
        margin: 0 4px;

        &.is-active {
          background: linear-gradient(135deg, #4fc3f7, #29b6f6);
        }

        &:hover {
          color: #0288d1;
        }
      }

      .btn-prev,
      .btn-next {
        border-radius: 4px;
        border: 1px solid #e1f5fe;
      }

      button:hover {
        color: #0288d1;
      }
    }
  }
}

/* 暗黑模式适配 */
html.dark {
  .series-list {
    .series-grid .series-card {
      background: #1e1e1e;
      border-color: #37474f;

      &:hover {
        border-color: #0288d1;
      }

      .card-content {
        h3 {
          color: #81d4fa;
        }

        .description {
          color: #b0bec5;
        }

        .meta .info {
          color: #90a4ae;
        }
      }
    }

    .series-detail {
      .detail-header {
        background: linear-gradient(135deg, #0d47a1 0%, #01579b 100%);
        border-color: #0288d1;

        .info-container {
          .series-title h2 {
            color: #e3f2fd;
          }

          .series-description p {
            color: #bbdefb;
            background: rgba(255, 255, 255, 0.1);
            border-left-color: #29b6f6;
          }

          .series-stats {
            background: rgba(255, 255, 255, 0.1);
            border-color: rgba(33, 150, 243, 0.3);

            .stat-item .stat-info {
              .stat-value {
                color: #e3f2fd;
              }

              .stat-label {
                color: #bbdefb;
              }
            }
          }
        }
      }

      .chapters-section {
        background: #1e1e1e;
        border-color: #37474f;

        .section-header {
          background: linear-gradient(to right, #0d47a1, #1a237e);
          border-bottom-color: #0288d1;

          .section-title {
            color: #bbdefb;
          }
        }

        .chapters-container {
          .chapters-list .chapter-item {
            background: #263238;
            border-color: #37474f;

            &:hover {
              border-color: #0288d1;
            }

            &.is-expanded {
              background: #1c313a;
            }

            .chapter-header {
              background: linear-gradient(to right, #1c313a, #263238);

              &:hover {
                background: linear-gradient(to right, #0d47a1, #1c313a);
              }

              .chapter-info {
                .chapter-main {
                  h4 {
                    color: #bbdefb;
                  }

                  .chapter-desc {
                    color: #90a4ae;
                  }
                }

                .chapter-meta .sort-order {
                  color: #78909c;
                }
              }
            }

            .subchapters-list {
              background: #1c313a;

              .subchapter-item {
                background: #263238;
                border-color: #37474f;

                &:hover {
                  border-color: #0288d1;
                  background: #1c313a;
                }

                .subchapter-content .subchapter-main {
                  h5 {
                    color: #bbdefb;
                  }

                  .subchapter-desc {
                    color: #90a4ae;
                  }
                }
              }
            }
          }

          .empty-chapters {
            background: #263238;
            border-color: #37474f;
          }
        }
      }
    }
  }
}

/* 响应式调整 */
@media (max-width: 992px) {
  .series-detail {
    .detail-header .header-content {
      flex-direction: column;
      align-items: center;
      text-align: center;

      .cover-container .cover-wrapper {
        margin: 0 auto;
      }

      .info-container {
        width: 100%;

        .series-stats {
          grid-template-columns: repeat(2, 1fr);
        }
      }
    }

    .chapters-container .chapters-list .chapter-item {
      .chapter-header .chapter-info {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;

        .chapter-number {
          align-self: flex-start;
        }

        .chapter-meta {
          width: 100%;
          justify-content: flex-start;
        }
      }

      .subchapters-list {
        padding-left: 20px;
      }
    }
  }
}

@media (max-width: 576px) {
  .series-detail {
    .detail-header .header-content {
      padding: 16px;

      .cover-container .cover-wrapper .cover-image {
        width: 140px;
        height: 180px;
      }

      .info-container .series-stats {
        grid-template-columns: 1fr;
      }
    }

    .chapters-section {
      .section-header {
        flex-direction: column;
        gap: 12px;
        align-items: flex-start;
      }

      .chapters-container .chapters-list .chapter-item {
        .chapter-header {
          flex-direction: column;
          gap: 12px;

          .chapter-actions {
            width: 100%;
            justify-content: flex-end;
          }
        }

        .subchapters-list .subchapter-item {
          flex-direction: column;
          gap: 12px;
          align-items: flex-start;

          .subchapter-actions {
            align-self: flex-end;
          }
        }
      }
    }
  }
}
</style>