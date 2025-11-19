<template>
    <div class="outline">
        <div class="outline__header">
            <div class="outline__header_left">
                <!-- 返回上一页 -->
                <div class="return-btn" @click="handleReturn">
                <el-icon :size="16">
                    <Back />
                </el-icon>
                </div>
                <h3 class="outline__title">目录</h3>        
            </div>
            <el-button
              v-if="props.showEdit"
              type="success"
              size="small"
              :icon="Plus"
              @click="handleAdd"
            ></el-button>
        </div>
        <div class="outline__content">
            <nav class="outline__nav">
                <div
                  v-for="(item, index) in outline"
                  :key="item.id"
                  class="outline__group"
                >
                  <!-- 父级目录项 -->
                  <div
                    class="outline__item outline__item--parent"
                    :class="{
                      'outline__item--active': activeItemId === item.id,
                      'outline__item--expanded': expandedItems.has(item.id)
                    }"
                    @click="handleParentClick(item, index)"
                    @contextmenu.prevent="props.showEdit ? handleContextMenu($event, item, index) : null"
                  >
                    <span 
                      class="outline__item-icon"
                      v-if="hasChildren(item)"
                      @click.stop="toggleExpand(item.id)"
                    >
                      <el-icon>
                        <ArrowRight v-if="!expandedItems.has(item.id)" />
                        <ArrowDown v-else />
                      </el-icon>
                    </span>
                    <span v-else class="outline__item-icon-placeholder"></span>
                    <span class="outline__item-text" :title="item.title">{{ item.title }}</span>
                  </div>

                  <!-- 子级目录项 -->
                  <div
                    v-if="item.children && item.children.length > 0 && expandedItems.has(item.id)"
                    class="outline__children"
                  >
                    <div
                      v-for="(child, childIndex) in item.children"
                      :key="child.id"
                      class="outline__item outline__item--child"
                      :class="{
                        'outline__item--active': activeItemId === child.id
                      }"
                      @click="handleChildClick(child, item, index)"
                      @contextmenu.prevent="props.showEdit ? handleContextMenu($event, child, index, childIndex, item) : null"
                    >
                      <span class="outline__item-icon-placeholder"></span>
                      <span class="outline__item-text" :title="child.title">{{ child.title }}</span>
                    </div>
                  </div>
                </div>
            </nav>
        </div>

        <!-- 右键菜单 -->
        <div
            v-if="contextMenuVisible"
            class="context-menu"
            :style="{ top: contextMenuPosition.y + 'px', left: contextMenuPosition.x + 'px' }"
            @click.stop
        >
            <div class="context-menu__item" @click="handleNewItem">
                <el-icon><Plus /></el-icon>
                <span>新建目录</span>
            </div>
            <div class="context-menu__item" @click="handleEditItem">
                <el-icon><Edit /></el-icon>
                <span>编辑</span>
            </div>
            <div class="context-menu__item context-menu__item--danger" @click="handleDeleteItem">
                <el-icon><Delete /></el-icon>
                <span>删除</span>
            </div>
        </div>

        <!-- 编辑对话框 -->
        <el-dialog
            v-model="editDialogVisible"
            :title="editingItem.id ? '编辑目录' : '新建目录'"
            width="400px"
            @close="handleEditDialogClose"
        >
            <el-form :model="editingItem" label-width="80px">
                <el-form-item label="目录标题">
                    <el-input v-model="editingItem.title" placeholder="请输入目录标题" />
                </el-form-item>
                <el-form-item label="上级目录">
                    <el-tree-select
                      v-model="editingItem.parentId"
                      node-key="id"
                      :data="outline"
                      :props="{
                        label: 'title',
                        value: 'id',
                        children: 'children'
                      }"
                      clearable
                      check-strictly
                      default-expand-all
                      placeholder="请选择上级目录"
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="editDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="handleSaveEdit">确定</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Back, Plus, Edit, Delete, ArrowRight, ArrowDown } from '@element-plus/icons-vue'
const router = useRouter()

export interface OutlineItem {
  id: string
  title: string
  parentId?: string|number
  children?: OutlineItem[]
  anchor?: string
}

interface Props {
  modelValue: OutlineItem[],
  showEdit?: {
    type: boolean,
    default: false,
  }
}

interface Emits {
  (e: 'item-click', item: OutlineItem, index: number): void
  (e: 'update:modelValue', value: OutlineItem[]): void
  (e: 'item-change', items: OutlineItem[]): void
}

const emit = defineEmits<Emits>()
const props = defineProps<Props>()

// 目录数据
const outline = ref<OutlineItem[]>([...props.modelValue])

// 监听外部数据变化
watch(() => props.modelValue, (newVal) => {
  outline.value = [...newVal]
  // // 构建树形结构
  // buildTree()
}, { deep: true })

// 展开的父级目录ID集合
const expandedItems = ref<Set<string>>(new Set())

// 当前激活的目录ID
const activeItemId = ref<string>('')

// 右键菜单相关
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const currentContextItem = ref<OutlineItem | null>(null)
const currentContextParent = ref<OutlineItem | null>(null)
const currentContextIndex = ref(-1)
const currentContextChildIndex = ref(-1)

// 编辑对话框相关
const editDialogVisible = ref(false)
const editingItem = ref<OutlineItem>({ id: '', title: '', parentId: undefined })

// 检查是否有子项
function hasChildren(item: OutlineItem): boolean {
  return !!(item.children && item.children.length > 0)
}

// 切换展开/折叠
function toggleExpand(itemId: string) {
  if (expandedItems.value.has(itemId)) {
    expandedItems.value.delete(itemId)
  } else {
    expandedItems.value.add(itemId)
  }
}

// 构建树形结构（如果目录数据是平级的，需要用此方法将扁平数据转换为树形）
function buildTree() {
  const flatData = [...props.modelValue]
  const itemMap = new Map<string, OutlineItem>()
  const rootItems: OutlineItem[] = []

  // 第一遍：创建所有项的映射，保留原有的 children
  flatData.forEach(item => {
    itemMap.set(item.id, { ...item, children: item.children || [] })
  })

  // 第二遍：构建父子关系
  flatData.forEach(item => {
    const mappedItem = itemMap.get(item.id)
    if (!mappedItem) return

    if (item.parentId) {
      const parent = itemMap.get(String(item.parentId))
      if (parent) {
        if (!parent.children) {
          parent.children = []
        }
        // 检查是否已存在
        if (!parent.children.find(child => child.id === mappedItem.id)) {
          parent.children.push(mappedItem)
        }
      } else {
        // 父级不存在，作为根节点
        if (!rootItems.find(root => root.id === mappedItem.id)) {
          rootItems.push(mappedItem)
        }
      }
    } else {
      // 没有父级，作为根节点
      if (!rootItems.find(root => root.id === mappedItem.id)) {
        rootItems.push(mappedItem)
      }
    }
  })

  // 更新 outline 为树形结构
  outline.value = rootItems

  console.log('outline', outline.value)
}

// 点击父级目录
function handleParentClick(item: OutlineItem, index: number) {
  activeItemId.value = item.id
  emit('item-click', item, index)
}

// 点击子级目录
function handleChildClick(child: OutlineItem, _parent: OutlineItem, parentIndex: number) {
  activeItemId.value = child.id
  emit('item-click', child, parentIndex)
}

// 添加目录项（根级）
function handleAdd() {
  editingItem.value = {
    id: '',
    title: '新目录',
    parentId: undefined,
  }
  editDialogVisible.value = true
}

onMounted(() => {
  // buildTree()
  // 默认展开所有父级
  outline.value.forEach(item => {
    if (hasChildren(item)) {
      expandedItems.value.add(item.id)
    }
  })

  if (props.showEdit) {
    document.addEventListener('click', handleClickOutside)
  }
})

// 右键菜单处理
function handleContextMenu(
  event: MouseEvent, 
  item: OutlineItem, 
  index: number, 
  childIndex?: number,
  parent?: OutlineItem
) {
  event.preventDefault()
  contextMenuPosition.value = {
    x: event.clientX,
    y: event.clientY
  }
  currentContextItem.value = item
  currentContextParent.value = parent || null
  currentContextIndex.value = index
  currentContextChildIndex.value = childIndex ?? -1
  contextMenuVisible.value = true
}

// 关闭右键菜单
function closeContextMenu() {
  contextMenuVisible.value = false
  currentContextItem.value = null
  currentContextParent.value = null
  currentContextIndex.value = -1
  currentContextChildIndex.value = -1
}

// 新建目录
function handleNewItem() {
  if (!currentContextItem.value) return
  
  // 如果当前项是父级，新建子级；如果是子级，新建同级
  editingItem.value = {
    id: '',
    title: '新目录',
    parentId: currentContextChildIndex.value === -1 ? currentContextItem.value.id : currentContextParent.value?.id || currentContextItem.value.id,
  }
  editDialogVisible.value = true
  closeContextMenu()
}

// 编辑目录
function handleEditItem() {
  if (!currentContextItem.value) return
  
  editingItem.value = {
    ...currentContextItem.value
  }
  editDialogVisible.value = true
  closeContextMenu()
}

// 保存编辑
function handleSaveEdit() {
  if (!editingItem.value.title.trim()) {
    return
  }

  // 如果是新建
  if (!editingItem.value.id) {
    const newId = `item-${Date.now()}`
    const newItem: OutlineItem = {
      id: newId,
      title: editingItem.value.title,
      parentId: editingItem.value.parentId,
    }

    if (newItem.parentId) {
      // 添加到父级的 children
      const parent = findItemById(outline.value, String(newItem.parentId))
      if (parent) {
        if (!parent.children) {
          parent.children = []
        }
        parent.children.push(newItem)
        // 展开父级
        expandedItems.value.add(String(newItem.parentId))
      }
    } else {
      // 添加到根级
      outline.value.push(newItem)
    }
  } else {
    // 更新现有项
    const item = findItemById(outline.value, editingItem.value.id)
    if (item) {
      item.title = editingItem.value.title
      const oldParentId = item.parentId
      item.parentId = editingItem.value.parentId

      // 如果父级改变，需要移动
      if (oldParentId !== editingItem.value.parentId) {
        moveItem(item, oldParentId, editingItem.value.parentId)
      }
    }
  }

  // // 扁平化数据用于 emit
  // const flatData = flattenTree(outline.value)
  emit('update:modelValue', outline.value)
  emit('item-change', outline.value)
  editDialogVisible.value = false
  closeContextMenu()
}

// 查找项
function findItemById(items: OutlineItem[], id: string): OutlineItem | null {
  for (const item of items) {
    if (item.id === id) return item
    if (item.children) {
      const found = findItemById(item.children, id)
      if (found) return found
    }
  }
  return null
}

// 移动项
function moveItem(item: OutlineItem, oldParentId: string | number | undefined, newParentId: string | number | undefined) {
  // 从旧位置移除
  if (oldParentId) {
    const oldParent = findItemById(outline.value, String(oldParentId))
    if (oldParent && oldParent.children) {
      const index = oldParent.children.findIndex(child => child.id === item.id)
      if (index > -1) {
        oldParent.children.splice(index, 1)
      }
    }
  } else {
    const index = outline.value.findIndex(root => root.id === item.id)
    if (index > -1) {
      outline.value.splice(index, 1)
    }
  }

  // 添加到新位置
  if (newParentId) {
    const newParent = findItemById(outline.value, String(newParentId))
    if (newParent) {
      if (!newParent.children) {
        newParent.children = []
      }
      newParent.children.push(item)
      expandedItems.value.add(String(newParentId))
    }
  } else {
    outline.value.push(item)
  }
}

// 扁平化树形数据
function flattenTree(items: OutlineItem[]): OutlineItem[] {
  const result: OutlineItem[] = []
  function traverse(items: OutlineItem[]) {
    items.forEach(item => {
      const { children, ...itemWithoutChildren } = item
      result.push(itemWithoutChildren)
      if (children && children.length > 0) {
        traverse(children)
      }
    })
  }
  traverse(items)
  return result
}

// 关闭编辑对话框
function handleEditDialogClose() {
  closeContextMenu()
  editingItem.value = { id: '', title: '', parentId: undefined }
}

// 删除目录
function handleDeleteItem() {
  if (!currentContextItem.value) return
  
  const itemId = currentContextItem.value.id
  const isChild = currentContextChildIndex.value !== -1

  if (isChild && currentContextParent.value) {
    // 删除子级
    if (currentContextParent.value.children) {
      const index = currentContextParent.value.children.findIndex(child => child.id === itemId)
      if (index > -1) {
        currentContextParent.value.children.splice(index, 1)
      }
    }
  } else {
    // 删除父级（同时删除所有子级）
    const index = outline.value.findIndex(item => item.id === itemId)
    if (index > -1) {
      outline.value.splice(index, 1)
    }
  }

  // 如果删除的是当前激活项，重置激活ID
  if (activeItemId.value === itemId) {
    activeItemId.value = ''
  }

  // // 扁平化数据用于 emit
  // const flatData = flattenTree(outline.value)
  emit('update:modelValue', outline.value)
  emit('item-change', outline.value)
  
  closeContextMenu()
}

// 点击其他地方关闭右键菜单
function handleClickOutside(event: MouseEvent) {
  if (contextMenuVisible.value) {
    const target = event.target as HTMLElement
    if (!target.closest('.context-menu') && !target.closest('.outline__item')) {
      closeContextMenu()
    }
  }
}

// 返回上一页
function handleReturn() {
  router.back()
}

onBeforeUnmount(() => {
  if (props.showEdit) {
    document.removeEventListener('click', handleClickOutside)
  }
})
</script>

<style lang="scss" scoped>
// 目录样式
.outline {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #ffffff;

  &__header {
    height: 40px;
    padding: 2px 8px;
    border-bottom: 1px solid #ddd;
    flex-shrink: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    &_left {
      display: flex;
      align-items: center;
      gap: 5px;
      .return-btn {
        width: 28px;
        height: 28px;
        border: 1px solid #ddd;
        color: #666;
        border-radius: 50%;
        display: flex;
        justify-content: center;
        align-items: center;
        &:hover {
          border-color: var(--el-color-primary);
          color: var(--el-color-primary);
        }
      }      
    }
  }

  &__title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #24292e;
  }

  &__content {
    flex: 1;
    overflow-y: auto;
    padding: 8px 0;
  }

  &__nav {
    display: flex;
    flex-direction: column;
  }

  &__group {
    display: flex;
    flex-direction: column;
  }

  &__item {
    position: relative;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 10px;
    cursor: pointer;
    transition: all 0.2s ease;
    border-left: 3px solid transparent;
    min-height: 36px;
    box-sizing: border-box;
    
    &:hover {
      background-color: #f6f8fa;
    }

    &--active {
      background-color: #f0f9ff;
      border-left-color: #67c23a;
      
      .outline__item-text {
        color: #67c23a;
        font-weight: 600;
      }
    }

    &--parent {
      font-weight: 500;
      font-size: 14px;
    }

    &--child {
      font-size: 13px;
      font-weight: 400;
    }

    &--expanded {
      .outline__item-icon {
        color: #67c23a;
      }
    }
  }

  &__item-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    color: #6a737d;
    transition: transform 0.2s ease, color 0.2s ease;
    cursor: pointer;

    &:hover {
      color: #67c23a;
    }

    .el-icon {
      font-size: 14px;
    }
  }

  &__item-icon-placeholder {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
  }

  &__item-text {
    color: #24292e;
    font-size: 14px;
    line-height: 1.5;
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: color 0.2s ease;
  }

  &__children {
    display: flex;
    flex-direction: column;
    background-color: #fafbfc;
  }
}

// 右键菜单样式
.context-menu {
  position: fixed;
  z-index: 9999;
  background-color: #ffffff;
  border: 1px solid #e1e4e8;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 4px 0;
  min-width: 140px;
  overflow: hidden;

  &__item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    cursor: pointer;
    font-size: 14px;
    color: #24292e;
    transition: background-color 0.2s ease;

    &:hover {
      background-color: #f6f8fa;
    }

    &--danger {
      color: #f56c6c;

      &:hover {
        background-color: #fef0f0;
      }
    }

    .el-icon {
      font-size: 16px;
    }
  }
}
</style>

