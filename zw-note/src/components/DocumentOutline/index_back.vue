<template>
  <div class="document-outline">
    <div class="document-outline__header">
      <div class="document-outline__header_left">
        <!-- 返回上一页 -->
        <div class="return-btn" @click="handleReturn">
          <el-icon :size="size" :color="color">
            <Back />
          </el-icon>
        </div>
        <h3 class="document-outline__title">目录</h3>        
      </div>
      <el-button
        type="success"
        size="small"
        :icon="Plus"
        @click="handleAdd"
      ></el-button>
    </div>
    <div class="document-outline__content">
      <draggable
        v-model="outlineItems"
        :animation="200"
        handle=".document-outline-item__drag-handle"
        @end="handleDragEnd"
      >
        <template #item="{ element, index }">
          <div
            class="document-outline-item"
            :class="{ 'is-active': activeIndex === index }"
            @click="handleItemClick(index, element)"
          >
            <div class="document-outline-item__drag-handle">
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="3" cy="3" r="1" fill="#9ca3af"/>
                <circle cx="9" cy="3" r="1" fill="#9ca3af"/>
                <circle cx="3" cy="6" r="1" fill="#9ca3af"/>
                <circle cx="9" cy="6" r="1" fill="#9ca3af"/>
                <circle cx="3" cy="9" r="1" fill="#9ca3af"/>
                <circle cx="9" cy="9" r="1" fill="#9ca3af"/>
              </svg>
            </div>
            <div class="document-outline-item__content">
              <input
                v-model="element.title"
                class="document-outline-item__input"
                type="text"
                :placeholder="`章节 ${index + 1}`"
                @click.stop
                @change="handleTitleChange(index, element)"
              />
            </div>
            <div class="document-outline-item__actions">
              <el-button
                class="px-1"
                type="text"
                size="small"
                @click.stop="handleDelete(index)"
              >
                删除
              </el-button>
            </div>
          </div>
        </template>
      </draggable>
      <div v-if="outlineItems.length === 0" class="document-outline__empty">
        <p>暂无目录，点击"添加章节"创建</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Back, Plus } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import draggable from 'vuedraggable'
import { ElMessageBox } from 'element-plus'

const router = useRouter()

export interface OutlineItem {
  id: string
  title: string
  level: number
  anchor?: string
}

interface Props {
  modelValue: OutlineItem[]
}

interface Emits {
  (e: 'update:modelValue', value: OutlineItem[]): void
  (e: 'item-click', item: OutlineItem, index: number): void
  (e: 'item-change', items: OutlineItem[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const outlineItems = ref<OutlineItem[]>([...props.modelValue])
const activeIndex = ref<number>(-1)

function handleAdd() {
  const newItem: OutlineItem = {
    id: `item-${Date.now()}`,
    title: '',
    level: 1,
  }
  outlineItems.value.push(newItem)
  emit('update:modelValue', outlineItems.value)
  emit('item-change', outlineItems.value)
}

function handleDelete(index: number) {
  ElMessageBox.confirm(
    '确定删除该章节吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(() => {
      outlineItems.value.splice(index, 1)
      emit('update:modelValue', outlineItems.value)
      emit('item-change', outlineItems.value)
      if (activeIndex.value === index) {
        activeIndex.value = -1
      } else if (activeIndex.value > index) {
        activeIndex.value--
      }
    })
    .catch(() => {
    })
}

function handleItemClick(index: number, item: OutlineItem) {
  activeIndex.value = index
  emit('item-click', item, index)
}

function handleTitleChange(index: number, item: OutlineItem) {
  emit('update:modelValue', outlineItems.value)
  emit('item-change', outlineItems.value)
}

function handleDragEnd() {
  emit('update:modelValue', outlineItems.value)
  emit('item-change', outlineItems.value)
}

watch(
  () => props.modelValue,
  (newVal) => {
    outlineItems.value = [...newVal]
  },
  { deep: true }
)

const handleReturn = () => {
  router.back()
}
</script>

<style scoped lang="scss">
.document-outline {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #ffffff;
  // border: 1px solid #e5e7eb;
  // border-radius: 4px;
  overflow: hidden;

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    // padding: 12px 16px;
    border-bottom: 1px solid #ddd;
    height: 40px;
    padding: 2px 8px;
    &_left {
      display: flex;
      align-items: center;
      gap: 5px;
      .return-btn {
        width: 28px;
        height: 28px;
        border: 1px solid #ddd;
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
    color: #374151;
  }

  &__content {
    flex: 1;
    overflow-y: auto;
    padding: 6px;
  }

  &__empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: #9ca3af;
    font-size: 13px;
  }
}

.document-outline-item {
  display: flex;
  align-items: center;
  padding: 4px;
  margin-bottom: 4px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: #ffffff;

  &:hover {
    border-color: #67c23a;
    background-color: #f9fafb;
  }

  &.is-active {
    border-color: #67c23a;
    background-color: #f0f9ff;
  }

  &__drag-handle {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    cursor: move;
    opacity: 0.5;
    transition: opacity 0.2s ease;

    &:hover {
      opacity: 1;
    }
  }

  &__content {
    flex: 1;
    min-width: 0;
  }

  &__input {
    width: 100%;
    padding: 4px 8px;
    box-sizing: border-box;
    border: none;
    background: transparent;
    font-size: 14px;
    font-weight: 500;
    color: #374151;
    outline: none;

    &::placeholder {
      color: #9ca3af;
    }

    &:focus {
      background-color: #ffffff;
      border: 1px solid #67c23a;
      border-radius: 2px;
    }
  }

  &__actions {
    opacity: 0;
    transition: opacity 0.2s ease;
  }

  &:hover &__actions {
    opacity: 1;
  }
}
.px-1 {
  padding-left: 2px!important;
  padding-right: 2px!important;
}
</style>

