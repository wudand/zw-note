<template>
    <div class="outline">
        <div class="outline__header">
            <div class="outline__header_left">
                <button type="button" class="outline__back" aria-label="返回" @click="handleReturn">
                  <el-icon :size="16"><Back /></el-icon>
                </button>
                <h3 class="outline__title" :title="props.title || '目录'">{{ props.title || '目录' }}</h3>
            </div>
            <button
              v-if="props.showEdit"
              type="button"
              class="outline__add"
              aria-label="新建目录"
              @click="handleAdd"
            >
              <el-icon :size="16"><Plus /></el-icon>
            </button>
        </div>
        <div class="outline__content">
            <nav class="outline__nav" ref="navRef">
              <!-- 预览：只读列表 -->
              <template v-if="!props.showEdit">
                <div
                  v-for="(item, index) in outline"
                  :key="item.id"
                  class="outline__group"
                >
                  <div
                    class="outline__item outline__item--parent"
                    :class="{
                      'outline__item--active': activeItemId === item.id,
                      'outline__item--expanded': expandedItems.has(item.id)
                    }"
                    @click="handleParentClick(item, index)"
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

                  <div
                    v-if="item.children && item.children.length > 0 && expandedItems.has(item.id)"
                    class="outline__children"
                  >
                    <div
                      v-for="child in item.children"
                      :key="child.id"
                      class="outline__item outline__item--child"
                      :class="{ 'outline__item--active': activeItemId === child.id }"
                      @click="handleChildClick(child, item, index)"
                    >
                      <span class="outline__item-icon-placeholder"></span>
                      <span class="outline__item-text" :title="child.title">{{ child.title }}</span>
                    </div>
                  </div>
                </div>
              </template>

              <!-- 编辑：可拖拽双层列表 -->
              <draggable
                v-else
                v-model="outline"
                item-key="id"
                group="outline"
                data-level="root"
                class="outline__draggable outline__draggable--root"
                :animation="180"
                :disabled="reordering"
                :delay="150"
                :delay-on-touch-only="true"
                :force-fallback="true"
                handle=".outline__drag-handle"
                ghost-class="outline__ghost"
                chosen-class="outline__chosen"
                drag-class="outline__drag"
                :move="checkMove"
                @start="onDragStart"
                @end="onDragEnd"
              >
                <template #item="{ element, index }">
                  <div class="outline__group">
                    <div
                      class="outline__item outline__item--parent"
                      :class="{
                        'outline__item--active': activeItemId === element.id,
                        'outline__item--expanded': expandedItems.has(element.id)
                      }"
                      @click="handleParentClick(element, index)"
                      @contextmenu.prevent="handleContextMenu($event, element, index)"
                    >
                      <span
                        class="outline__drag-handle"
                        title="拖动排序"
                        @click.stop
                      >
                        <el-icon :size="14"><Rank /></el-icon>
                      </span>
                      <span
                        class="outline__item-icon"
                        v-if="hasChildren(element)"
                        @click.stop="toggleExpand(element.id)"
                      >
                        <el-icon>
                          <ArrowRight v-if="!expandedItems.has(element.id)" />
                          <ArrowDown v-else />
                        </el-icon>
                      </span>
                      <span v-else class="outline__item-icon-placeholder"></span>
                      <span class="outline__item-text" :title="element.title">{{ element.title }}</span>
                    </div>

                    <draggable
                      v-model="element.children"
                      item-key="id"
                      group="outline"
                      data-level="child"
                      class="outline__draggable outline__draggable--child outline__children"
                      :class="{ 'outline__children--collapsed': !isDragging && (!hasChildren(element) || !expandedItems.has(element.id)) }"
                      :animation="180"
                      :disabled="reordering"
                      :delay="150"
                      :delay-on-touch-only="true"
                      :force-fallback="true"
                      handle=".outline__drag-handle"
                      ghost-class="outline__ghost"
                      chosen-class="outline__chosen"
                      drag-class="outline__drag"
                      :move="checkMove"
                      @start="onDragStart"
                      @end="onDragEnd"
                    >
                      <template #item="{ element: child, index: childIndex }">
                        <div
                          class="outline__item outline__item--child"
                          :class="{ 'outline__item--active': activeItemId === child.id }"
                          @click="handleChildClick(child, element, index)"
                          @contextmenu.prevent="handleContextMenu($event, child, index, childIndex, element)"
                        >
                          <span
                            class="outline__drag-handle"
                            title="拖动排序"
                            @click.stop
                          >
                            <el-icon :size="14"><Rank /></el-icon>
                          </span>
                          <span class="outline__item-icon-placeholder"></span>
                          <span class="outline__item-text" :title="child.title">{{ child.title }}</span>
                        </div>
                      </template>
                    </draggable>
                  </div>
                </template>
              </draggable>
            </nav>
        </div>

        <!-- 右键菜单：仅一级目录可「新建目录」 -->
        <div
            v-if="contextMenuVisible"
            class="context-menu"
            :style="{ top: contextMenuPosition.y + 'px', left: contextMenuPosition.x + 'px' }"
            @click.stop
        >
            <div
              v-if="contextIsRootLevel"
              class="context-menu__item"
              @click="handleNewItem"
            >
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

        <!-- 编辑对话框：上级目录仅可选一级节点 -->
        <el-dialog
            v-model="editDialogVisible"
            :title="editingItem.id ? '编辑目录' : '新建目录'"
            width="400px"
            @opened="focusTitleInput"
            @close="handleEditDialogClose"
        >
            <el-form :model="editingItem" label-width="80px">
                <el-form-item label="目录标题">
                    <el-input ref="titleInputRef" v-model="editingItem.title" placeholder="请输入目录标题" />
                </el-form-item>
                <el-form-item label="上级目录">
                    <el-select
                      v-model="editingItem.parentId"
                      clearable
                      filterable
                      placeholder="不选则为一级目录"
                      style="width: 100%"
                      :disabled="editingHasChildren"
                    >
                      <el-option
                        v-for="opt in parentSelectOptions"
                        :key="opt.id"
                        :label="opt.title"
                        :value="opt.id"
                        :disabled="opt.id === editingItem.id"
                      />
                    </el-select>
                    <div v-if="editingHasChildren" class="outline__hint">
                      已有子目录的一级节点不能改为二级
                    </div>
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="editDialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="saving" @click="handleSaveEdit">确定</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { InputInstance } from 'element-plus'
import { Back, Plus, Edit, Delete, ArrowRight, ArrowDown, Rank } from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { useDocumentStore } from '@/store/documentStore'

const router = useRouter()
const store = useDocumentStore()

export interface OutlineItem {
  id: string
  title: string
  parentId?: string|number
  children?: OutlineItem[]
  anchor?: string
}

interface ReorderItem {
  id: number
  parent_id: number | null
  sort_order: number
}

interface Props {
  modelValue: OutlineItem[],
  showEdit?: boolean
  /** 外部控制的当前选中节点（预览页跳转子目录时同步高亮） */
  activeId?: string
  /** 目录上方标题：不传时展示"目录"（预览页保持不变，编辑页可传当前文档名称） */
  title?: string
}

interface Emits {
  (e: 'item-click', item: OutlineItem, index: number): void
  (e: 'update:modelValue', value: OutlineItem[]): void
  (e: 'item-change', items: OutlineItem[]): void
}

const emit = defineEmits<Emits>()
const props = withDefaults(defineProps<Props>(), {
  showEdit: false,
  activeId: '',
})

function cloneOutline(items: OutlineItem[]): OutlineItem[] {
  return items.map((item) => ({
    ...item,
    children: cloneOutline(item.children ?? []),
  }))
}

function normalizeOutline(items: OutlineItem[]): OutlineItem[] {
  return items.map((item) => ({
    ...item,
    children: normalizeOutline(item.children ?? []),
  }))
}

function flattenOutline(tree: OutlineItem[]): ReorderItem[] {
  const items: ReorderItem[] = []
  tree.forEach((root, i) => {
    items.push({ id: Number(root.id), parent_id: null, sort_order: i })
    ;(root.children ?? []).forEach((child, j) => {
      items.push({
        id: Number(child.id),
        parent_id: Number(root.id),
        sort_order: j,
      })
    })
  })
  return items
}

function assertValidTwoLevelTree(tree: OutlineItem[]): boolean {
  for (const root of tree) {
    for (const child of root.children ?? []) {
      if (child.children && child.children.length > 0) return false
    }
  }
  return true
}

// 目录数据
const outline = ref<OutlineItem[]>(normalizeOutline(props.modelValue))
const saving = ref(false)
const reordering = ref(false)
const outlineSnapshot = ref<OutlineItem[]>([])
const dragSignature = ref('')
/** 是否正处于拖拽中：拖拽期间即使子列表为空也要显示出来充当放置区，
 * 但不写入 expandedItems，避免污染持久的展开状态（历史上出现过污染后
 * 无法清理、导致某个空目录之后一直露出占位空白的问题） */
const isDragging = ref(false)
const navRef = ref<HTMLElement>()

/**
 * 跨一级/二级列表拖拽时，SortableJS 会直接操作 DOM 添加 ghost/chosen/drag 样式类；
 * 元素在两个独立的 Sortable 实例间搬运后，个别情况下这些类不会被正常移除
 * （表现为拖过的项残留浅绿色 ghost 背景、看起来"间隔变宽"）。拖拽结束后统一清一遍作为兜底。
 */
function cleanupDragArtifacts() {
  nextTick(() => {
    const root = navRef.value
    if (!root) return
    const staleClasses = [
      'outline__ghost',
      'outline__chosen',
      'outline__drag',
      'sortable-ghost',
      'sortable-chosen',
      'sortable-drag',
      'sortable-fallback',
    ]
    root.querySelectorAll(`.${staleClasses.join(', .')}`).forEach((el) => {
      el.classList.remove(...staleClasses)
    })
  })
}

function hasItemId(items: OutlineItem[], id: string): boolean {
  for (const item of items) {
    if (item.id === id) return true
    if (item.children && hasItemId(item.children, id)) return true
  }
  return false
}

// 监听外部数据变化（如 store 重新拉取树 / 排序落库后校准）
watch(() => props.modelValue, (newVal) => {
  outline.value = normalizeOutline(newVal)
  newVal.forEach((item) => {
    if (item.children && item.children.length > 0) {
      expandedItems.value.add(item.id)
    }
  })
  if (!hasItemId(newVal, activeItemId.value)) {
    activeItemId.value = newVal[0]?.id || ''
  }
}, { deep: true })

// 展开的父级目录ID集合
const expandedItems = ref<Set<string>>(new Set())

// 当前激活的目录ID
const activeItemId = ref<string>(props.activeId || props.modelValue[0]?.id || '')

watch(
  () => props.activeId,
  (id) => {
    if (!id) return
    activeItemId.value = id
    // 若选中的是子节点，展开其父级
    for (const item of outline.value) {
      if (item.children?.some((child) => child.id === id)) {
        expandedItems.value.add(item.id)
        break
      }
    }
  },
)

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
const titleInputRef = ref<InputInstance>()

/** 弹窗打开动画结束后聚焦标题输入框，并全选默认文案方便直接覆盖输入 */
function focusTitleInput() {
  nextTick(() => {
    titleInputRef.value?.focus()
    titleInputRef.value?.select()
  })
}

/** 约定仅两级：上级目录只能选一级（根）节点 */
const parentSelectOptions = computed(() =>
  outline.value.map((item) => ({ id: item.id, title: item.title })),
)

/** 当前右键是否点在一级目录上（二级不可「新建目录」） */
const contextIsRootLevel = computed(() => currentContextChildIndex.value === -1)

/** 编辑中的节点若已有子节点，则必须保持为一级，不能挂到其他节点下 */
const editingHasChildren = computed(() => {
  if (!editingItem.value.id) return false
  const node = outline.value.find((item) => item.id === editingItem.value.id)
  return !!(node?.children && node.children.length > 0)
})

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

/** 拖拽过程中禁止：有子的一级降为二级、拖进自己的子列表 */
function checkMove(evt: {
  to?: HTMLElement
  draggedContext: { element: OutlineItem }
  relatedContext: { list?: OutlineItem[] }
}): boolean {
  const dragged = evt.draggedContext.element
  const toLevel = evt.to?.dataset?.level
  const relatedList = evt.relatedContext.list

  if (relatedList && dragged.children && relatedList === dragged.children) {
    return false
  }

  if (toLevel === 'child' && hasChildren(dragged)) {
    return false
  }

  return true
}

function onDragStart() {
  outlineSnapshot.value = cloneOutline(outline.value)
  dragSignature.value = JSON.stringify(flattenOutline(outline.value))
  // 拖拽期间让空子列表也显示出来充当放置区（仅样式层面，不写 expandedItems）
  isDragging.value = true
  closeContextMenu()
}

/** 一级节点拿到子节点后自动标记为展开，方便用户立刻看到刚拖进去的项 */
function autoExpandParents() {
  outline.value.forEach((item) => {
    if (hasChildren(item)) expandedItems.value.add(item.id)
  })
}

async function onDragEnd() {
  isDragging.value = false

  const nextSignature = JSON.stringify(flattenOutline(outline.value))
  if (nextSignature === dragSignature.value) {
    cleanupDragArtifacts()
    return
  }

  if (!assertValidTwoLevelTree(outline.value)) {
    outline.value = cloneOutline(outlineSnapshot.value)
    cleanupDragArtifacts()
    ElMessage.warning('目录最多两级，已有子目录的节点不能降为二级')
    return
  }

  // 同步本地 parentId，便于后续编辑弹窗
  outline.value.forEach((root) => {
    root.parentId = undefined
    ;(root.children ?? []).forEach((child) => {
      child.parentId = root.id
      child.children = child.children ?? []
    })
  })
  autoExpandParents()

  const items = flattenOutline(outline.value)
  if (items.some((item) => !Number.isFinite(item.id))) {
    outline.value = cloneOutline(outlineSnapshot.value)
    cleanupDragArtifacts()
    ElMessage.error('目录数据异常，无法排序')
    return
  }

  reordering.value = true
  try {
    await store.reorderOutlineNodes(items)
    emit('update:modelValue', cloneOutline(outline.value))
    emit('item-change', cloneOutline(outline.value))
  } catch (error: any) {
    outline.value = cloneOutline(outlineSnapshot.value)
    ElMessage.error(error?.message || '排序失败，已恢复原顺序')
  } finally {
    cleanupDragArtifacts()
    reordering.value = false
    dragSignature.value = ''
  }
}

onMounted(() => {
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

// 新建目录：仅一级节点右键可新建其子目录
function handleNewItem() {
  if (!currentContextItem.value || !contextIsRootLevel.value) {
    ElMessage.warning('仅一级目录下可新建子目录')
    closeContextMenu()
    return
  }

  editingItem.value = {
    id: '',
    title: '新目录',
    parentId: currentContextItem.value.id,
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

// 保存编辑（对接后端创建 / 更新目录节点；约定仅两级）
async function handleSaveEdit() {
  if (!editingItem.value.title.trim()) {
    ElMessage.warning('请输入目录标题')
    return
  }

  const title = editingItem.value.title.trim()
  let parentId = editingItem.value.parentId

  // 上级必须是一级节点；有子节点的一级目录不能降为二级
  if (parentId != null && parentId !== '') {
    const parentIsRoot = outline.value.some((item) => item.id === String(parentId))
    if (!parentIsRoot) {
      ElMessage.warning('上级目录只能选择一级目录')
      return
    }
    if (editingHasChildren.value) {
      ElMessage.warning('已有子目录的节点不能改为二级')
      return
    }
    if (editingItem.value.id && String(parentId) === String(editingItem.value.id)) {
      ElMessage.warning('不能将自己设为上级目录')
      return
    }
  } else {
    parentId = undefined
  }

  saving.value = true
  try {
    if (!editingItem.value.id) {
      await store.createOutlineNode({ title, parentId })
      if (parentId) expandedItems.value.add(String(parentId))
      ElMessage.success('创建成功')
    } else {
      await store.updateOutlineNode({
        id: editingItem.value.id,
        title,
        parentId,
      })
      if (parentId) expandedItems.value.add(String(parentId))
      ElMessage.success('保存成功')
    }

    editDialogVisible.value = false
    closeContextMenu()
  } catch (error: any) {
    ElMessage.error(error?.message || '操作失败，请稍后重试')
  } finally {
    saving.value = false
  }
}

// 关闭编辑对话框
function handleEditDialogClose() {
  closeContextMenu()
  editingItem.value = { id: '', title: '', parentId: undefined }
}

// 删除目录（对接后端，子节点由外键级联删除）
async function handleDeleteItem() {
  if (!currentContextItem.value) return

  const itemId = currentContextItem.value.id
  const title = currentContextItem.value.title
  closeContextMenu()

  try {
    await ElMessageBox.confirm(
      `确定删除目录「${title}」吗？其下子目录也会一并删除。`,
      '提示',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' },
    )
  } catch {
    return
  }

  try {
    await store.removeOutlineNode(itemId)
    if (activeItemId.value === itemId) {
      activeItemId.value = ''
    }
    ElMessage.success('删除成功')
  } catch (error: any) {
    ElMessage.error(error?.message || '删除失败，请稍后重试')
  }
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
/*
  Notion / Linear 目录侧栏
  选中：圆角底填充 + 字重，不用左侧色条
*/
.outline {
  --bg: #f7f7f5;
  --surface: #ffffff;
  --fg: #111111;
  --muted: #6b6b6b;
  --faint: #9b9b9b;
  --line: #e8e8e5;
  --hover: #efefed;
  --active: #e8eee7;
  --accent: #5a9e58;
  --danger: #b91c1c;
  --radius: 6px;
  --ease: 160ms ease;

  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg);
  color: var(--fg);
  font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", ui-sans-serif, system-ui,
    sans-serif;
  -webkit-font-smoothing: antialiased;

  &__header {
    height: 48px;
    padding: 0 10px;
    border-bottom: 1px solid var(--line);
    flex-shrink: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
    background: var(--bg);
  }

  &__header_left {
    display: flex;
    align-items: center;
    gap: 6px;
    min-width: 0;
  }

  &__back,
  &__add {
    width: 32px;
    height: 32px;
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: var(--radius);
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    transition: background-color var(--ease), color var(--ease);

    &:hover {
      background: var(--hover);
      color: var(--fg);
    }

    &:focus-visible {
      box-shadow: 0 0 0 2px color-mix(in srgb, var(--accent) 28%, transparent);
    }
  }

  &__add:hover {
    color: var(--accent);
    background: color-mix(in srgb, var(--accent) 12%, transparent);
  }

  &__title {
    margin: 0;
    font-size: 13px;
    font-weight: 600;
    letter-spacing: -0.01em;
    color: var(--fg);
    min-width: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &__hint {
    margin-top: 6px;
    font-size: 12px;
    color: var(--muted);
    line-height: 1.4;
  }

  &__content {
    flex: 1;
    overflow-y: auto;
    padding: 8px 8px 16px;
  }

  &__nav {
    display: flex;
    flex-direction: column;
    gap: 1px;
  }

  &__draggable {
    display: flex;
    flex-direction: column;
    gap: 1px;
    min-height: 4px;
  }

  &__group {
    display: flex;
    flex-direction: column;
    gap: 1px;
  }

  &__item {
    display: flex;
    align-items: center;
    gap: 6px;
    min-height: 32px;
    padding: 5px 8px;
    border: none;
    border-radius: var(--radius);
    background: transparent;
    cursor: pointer;
    box-sizing: border-box;
    touch-action: manipulation;
    transition: background-color var(--ease), color var(--ease);

    &:hover {
      background: var(--hover);

      .outline__drag-handle {
        opacity: 1;
      }
    }

    &:active {
      background: color-mix(in srgb, var(--accent) 10%, var(--hover));
    }

    /* 选中：整行浅底，不用左侧色条 */
    &--active {
      background: var(--active);

      .outline__item-text {
        color: var(--fg);
        font-weight: 600;
      }

      .outline__item-icon {
        color: var(--accent);
      }

      &:hover {
        background: color-mix(in srgb, var(--accent) 16%, #fff);
      }
    }

    &--parent {
      font-size: 13px;
    }

    &--child {
      font-size: 13px;
      padding-left: 8px;
      margin-left: 14px;
    }

    &--expanded .outline__item-icon {
      color: var(--muted);
    }
  }

  &__drag-handle {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 18px;
    flex-shrink: 0;
    margin-left: -2px;
    color: var(--faint);
    opacity: 0.35;
    cursor: grab;
    border-radius: 4px;
    transition: opacity var(--ease), color var(--ease), background-color var(--ease);

    &:hover {
      opacity: 1;
      color: var(--muted);
      background: color-mix(in srgb, var(--fg) 6%, transparent);
    }

    &:active {
      cursor: grabbing;
    }
  }

  &__item-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 18px;
    height: 18px;
    flex-shrink: 0;
    border-radius: 4px;
    color: var(--faint);
    transition: color var(--ease), background-color var(--ease);
    cursor: pointer;

    &:hover {
      color: var(--fg);
      background: color-mix(in srgb, var(--fg) 6%, transparent);
    }

    .el-icon {
      font-size: 12px;
    }
  }

  &__item-icon-placeholder {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
  }

  &__item-text {
    color: var(--fg);
    font-size: inherit;
    font-weight: 450;
    line-height: 1.4;
    letter-spacing: -0.01em;
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__item--child .outline__item-text {
    color: var(--muted);
    font-weight: 400;
  }

  &__item--child.outline__item--active .outline__item-text {
    color: var(--fg);
    font-weight: 600;
  }

  &__children {
    display: flex;
    flex-direction: column;
    gap: 1px;
    padding: 1px 0 4px;
    min-height: 8px;

    &--collapsed {
      display: none;
    }
  }
}

.outline__ghost {
  opacity: 0.45;
  background: color-mix(in srgb, var(--accent, #5a9e58) 12%, #fff) !important;
}

.outline__chosen {
  background: var(--hover, #efefed);
}

.outline__drag {
  opacity: 0.92;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
}

.context-menu {
  position: fixed;
  z-index: 9999;
  background: var(--surface, #fff);
  border: 1px solid #e8e8e5;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  padding: 4px;
  min-width: 148px;
  overflow: hidden;

  &__item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 10px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 13px;
    color: #111;
    transition: background-color 160ms ease;

    &:hover {
      background: #efefed;
    }

    &--danger {
      color: #b91c1c;

      &:hover {
        background: #fef2f2;
      }
    }

    .el-icon {
      font-size: 15px;
    }
  }
}

@media (prefers-reduced-motion: reduce) {
  .outline__back,
  .outline__add,
  .outline__item,
  .outline__item-icon,
  .outline__drag-handle,
  .context-menu__item {
    transition: none;
  }
}
</style>
