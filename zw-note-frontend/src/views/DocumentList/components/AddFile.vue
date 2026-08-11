<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="80px"
      label-position="top"
    >
      <el-form-item label="文档名称" prop="title">
        <el-input
          v-model="formData.title"
          placeholder="请输入文档名称"
          maxlength="50"
          show-word-limit
          clearable
        />
      </el-form-item>

      <el-form-item label="作者" prop="author">
        <el-input
          v-model="formData.author"
          placeholder="请输入作者名称"
          maxlength="20"
          show-word-limit
          clearable
        />
      </el-form-item>

      <el-form-item label="文档描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="4"
          placeholder="请输入文档描述"
          maxlength="200"
          show-word-limit
          resize="none"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="success" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useDocumentStore } from '@/store/documentStore'

const router = useRouter()
const store = useDocumentStore()

interface DocumentFormData {
  id?: string
  title: string
  author: string
  description: string
}

interface Emits {
  (e: 'refresh'): void
}

const emit = defineEmits<Emits>()

const dialogVisible = ref(false)
const mode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const submitting = ref(false)

const dialogTitle = computed(() => (mode.value === 'edit' ? '文档设置' : '创建新文档'))

const formData = reactive<DocumentFormData>({
  id: '',
  title: '',
  author: '',
  description: '',
})

const rules: FormRules<DocumentFormData> = {
  title: [
    { required: true, message: '请输入文档名称', trigger: 'blur' },
    { min: 1, max: 50, message: '文档名称长度在 1 到 50 个字符', trigger: 'blur' },
  ],
  author: [
    { required: true, message: '请输入作者名称', trigger: 'blur' },
    { min: 1, max: 20, message: '作者名称长度在 1 到 20 个字符', trigger: 'blur' },
  ],
  description: [
    { max: 200, message: '文档描述不能超过 200 个字符', trigger: 'blur' },
  ],
}

/** 打开弹窗：create=新建；edit=列表「设置」修改基础信息 */
function open(type: 'create' | 'edit' = 'create', data?: DocumentFormData) {
  mode.value = type
  if (type === 'edit') {
    formData.id = data?.id
    formData.title = data?.title || ''
    formData.author = data?.author || ''
    formData.description = data?.description || ''
  } else {
    resetForm()
  }
  dialogVisible.value = true
}

function handleClose() {
  dialogVisible.value = false
  resetForm()
}

function resetForm() {
  mode.value = 'create'
  formData.id = ''
  formData.title = ''
  formData.author = ''
  formData.description = ''
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    const payload = {
      title: formData.title,
      author: formData.author,
      description: formData.description,
    }

    if (mode.value === 'edit' && formData.id) {
      // 对应列表页下拉「设置」
      await store.editDocument({ id: formData.id, ...payload })
      ElMessage.success('保存成功')
      emit('refresh')
      handleClose()
      return
    }

    const doc = await store.addDocument(payload)
    ElMessage.success('创建成功')
    handleClose()
    router.push({ name: 'document-edit', params: { id: doc.id } })
  } catch (error: any) {
    ElMessage.error(error?.message || '操作失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

defineExpose({
  open,
})
</script>

<style scoped lang="scss">
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #24292e;
}

:deep(.el-input__wrapper) {
  border-radius: 4px;
}

:deep(.el-textarea__inner) {
  border-radius: 4px;
  font-family: inherit;
}
</style>
