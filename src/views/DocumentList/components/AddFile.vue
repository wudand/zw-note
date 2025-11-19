<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建新文档"
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
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { createDocument, updateDocument } from '@/service/api/documentList'
import { useRouter } from 'vue-router'

const router = useRouter()

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
const formRef = ref<FormInstance>()
const submitting = ref(false)

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

// 打开弹窗
function open(type: 'create' | 'edit' = 'create', data?: DocumentFormData) {
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

// 关闭弹窗
function handleClose() {
  dialogVisible.value = false
  resetForm()
}

// 重置表单
function resetForm() {
  formData.id = ''
  formData.title = ''
  formData.author = ''
  formData.description = ''
  formRef.value?.clearValidate()
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return
  submitting.value = true

  try {
    await formRef.value.validate()

    const api = formData.id ? updateDocument : createDocument

    const res = await api(formData)
    console.log("res", res)
    if (!formData.id) {
      router.push({ name: 'document-edit', params: { id: formData.id } })
    }
    emit('refresh')
    handleClose()
  } catch (error) {
    console.log('表单验证失败', error)
  } finally {
    submitting.value = false
  }
}

// 暴露方法供父组件调用
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

