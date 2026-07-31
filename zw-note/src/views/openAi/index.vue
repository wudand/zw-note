<template>
    <div class="chat-container">
        <!-- 头部标题 -->
        <header class="chat-header">
            <h1>DeepSeek AI助手</h1>
            <p class="subtitle">与AI进行智能对话</p>
        </header>

        <!-- 消息显示区域 -->
        <div class="messages-container" ref="messagesContainer">
            <div 
                v-for="(message, index) in messages" 
                :key="index"
                :class="['message', message.role]"
            >
                <div class="message-avatar">
                    <div v-if="message.role === 'user'" class="avatar user-avatar">
                        <span>我</span>
                    </div>
                    <div v-else class="avatar ai-avatar">
                        <span>AI</span>
                    </div>
                </div>
                <div class="message-content">
                    <div class="message-header">
                        <span class="sender">{{ message.role === 'user' ? '我' : 'DeepSeek助手' }}</span>
                        <span class="time">{{ message.time }}</span>
                    </div>
                    <div class="message-text" v-html="formatMessage(message.content)"></div>
                </div>
            </div>
            
            <!-- 加载指示器 -->
            <div v-if="loading" class="message ai">
                <div class="message-avatar">
                    <div class="avatar ai-avatar">
                        <span>AI</span>
                    </div>
                </div>
                <div class="message-content">
                    <div class="message-header">
                        <span class="sender">DeepSeek助手</span>
                    </div>
                    <div class="message-text">
                        <div class="typing-indicator">
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-container">
            <div class="input-wrapper">
                <textarea
                    v-model="userInput"
                    @keydown.enter.exact.prevent="sendMessage"
                    @keydown.enter.shift.exact.prevent="userInput += '\n'"
                    placeholder="输入您的问题..."
                    rows="1"
                    ref="inputRef"
                    class="message-input"
                ></textarea>
                <div class="input-actions">
                    <button 
                        @click="sendMessage" 
                        :disabled="loading || !userInput.trim()"
                        class="send-button"
                    >
                        <svg v-if="!loading" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M22 2L11 13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M22 2L15 22L11 13L2 9L22 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                        <div v-else class="loading-spinner"></div>
                    </button>
                </div>
            </div>
            <div class="input-hint">
                按 Enter 发送，Shift+Enter 换行
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import OpenAI from "openai";
// sk-b6410926bb0c40d69a7f50026ecc2689

// 消息接口
interface Message {
  role: 'user' | 'assistant'
  content: string
  time: string
}

const openai = new OpenAI({
    baseURL: 'https://api.deepseek.com',
    apiKey: 'sk-b6410926bb0c40d69a7f50026ecc2689',
    dangerouslyAllowBrowser: true
});

// 响应式数据
const userInput = ref('')
const messages = ref<Message[]>([
  {
    role: 'assistant',
    content: '您好！我是DeepSeek AI助手，很高兴为您服务。请问有什么可以帮您的吗？',
    time: getCurrentTime()
  }
])
const loading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLTextAreaElement | null>(null)

// 获取当前时间
function getCurrentTime(): string {
  const now = new Date()
  return `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
}

// 格式化消息内容（简单的换行处理）
function formatMessage(content: string): string {
  return content.replace(/\n/g, '<br>')
}

// 自动调整输入框高度
function adjustTextareaHeight() {
  if (!inputRef.value) return
  
  inputRef.value.style.height = 'auto'
  inputRef.value.style.height = Math.min(inputRef.value.scrollHeight, 120) + 'px'
}

// 滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 模拟AI回复
async function simulateAIResponse(userMessage: string): Promise<string> {
    const completion = await openai.chat.completions.create({
        messages: [{ role: "user", content: userMessage }],
        model: "deepseek-chat",
    });
console.log(completion.choices[0].message.content)
  if (!completion.choices[0].message.content) {
    return '抱歉，我暂时无法回答您的问题。请稍后再试。';
  }
  return completion.choices[0].message.content;
}

// 发送消息
async function sendMessage() {
  const content = userInput.value.trim()
  if (!content || loading.value) return
  
  // 添加用户消息
  const userMessage: Message = {
    role: 'user',
    content,
    time: getCurrentTime()
  }
  messages.value.push(userMessage)
  
  // 清空输入框并重置高度
  userInput.value = ''
  if (inputRef.value) {
    inputRef.value.style.height = 'auto'
  }
  
  // 滚动到底部
  scrollToBottom()
  
  // 显示加载状态
  loading.value = true
  
  try {
    // 模拟AI回复
    const aiResponse = await simulateAIResponse(content)
    
    // 添加AI回复
    const aiMessage: Message = {
      role: 'assistant',
      content: aiResponse,
      time: getCurrentTime()
    }
    messages.value.push(aiMessage)
    
    // 滚动到底部
    scrollToBottom()
  } catch (error) {
    console.error('获取AI回复失败:', error)
    
    // 添加错误消息
    const errorMessage: Message = {
      role: 'assistant',
      content: '抱歉，我暂时无法回答您的问题。请稍后再试。',
      time: getCurrentTime()
    }
    messages.value.push(errorMessage)
  } finally {
    loading.value = false
  }
}

// 监听输入框变化调整高度
watch(userInput, () => {
  adjustTextareaHeight()
})

// 组件挂载时聚焦输入框
onMounted(() => {
    
  if (inputRef.value) {
    inputRef.value.focus()
  }
  scrollToBottom()
})
</script>

<style scoped lang="scss">
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 800px;
  margin: 0 auto;
  background-color: #f8f9fa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
}

.chat-header {
  text-align: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  
  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
  }
  
  .subtitle {
    margin: 8px 0 0;
    font-size: 14px;
    opacity: 0.9;
  }
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  scroll-behavior: smooth;
  
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: #f1f1f1;
  }
  
  &::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 3px;
  }
}

.message {
  display: flex;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease-out;
  
  &.user {
    flex-direction: row-reverse;
    
    .message-avatar {
      margin-left: 12px;
      margin-right: 0;
    }
    
    .message-content {
      align-items: flex-end;
    }
    
    .message-text {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-radius: 18px 4px 18px 18px;
    }
  }
  
  &.ai {
    .message-text {
      background-color: white;
      color: #333;
      border-radius: 4px 18px 18px 18px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-avatar {
  margin-right: 12px;
  flex-shrink: 0;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  
  &.user-avatar {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }
  
  &.ai-avatar {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    color: white;
  }
}

.message-content {
  display: flex;
  flex-direction: column;
  max-width: 70%;
}

.message-header {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
}

.sender {
  font-weight: 600;
  font-size: 14px;
  color: #555;
}

.time {
  font-size: 12px;
  color: #999;
  margin-left: 8px;
}

.message-text {
  padding: 12px 16px;
  line-height: 1.5;
  font-size: 15px;
  word-break: break-word;
}

.input-container {
  padding: 20px;
  background-color: white;
  border-top: 1px solid #eaeaea;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}

.input-wrapper {
  display: flex;
  align-items: flex-end;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  background-color: white;
  padding: 8px;
  transition: border-color 0.2s;
  
  &:focus-within {
    border-color: #667eea;
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
  }
}

.message-input {
  flex: 1;
  border: none;
  outline: none;
  resize: none;
  font-size: 15px;
  line-height: 1.5;
  padding: 8px;
  max-height: 120px;
  font-family: inherit;
  
  &::placeholder {
    color: #999;
  }
}

.input-actions {
  display: flex;
  align-items: center;
  margin-left: 8px;
}

.send-button {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  
  &:hover:not(:disabled) {
    transform: scale(1.05);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.typing-indicator {
  display: flex;
  align-items: center;
  height: 20px;
  
  span {
    height: 8px;
    width: 8px;
    background-color: #bbb;
    border-radius: 50%;
    display: inline-block;
    margin-right: 4px;
    animation: bounce 1.4s infinite ease-in-out both;
    
    &:nth-child(1) {
      animation-delay: -0.32s;
    }
    
    &:nth-child(2) {
      animation-delay: -0.16s;
    }
  }
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

.input-hint {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

@media (max-width: 768px) {
  .chat-container {
    height: 100vh;
  }
  
  .message-content {
    max-width: 85%;
  }
  
  .chat-header h1 {
    font-size: 20px;
  }
}
</style>