<template>
  <div class="container">
    <h1>多行文本流式处理</h1>
    <textarea
        v-model="textInput"
        placeholder="请输入多行文字..."
        rows="10"
        class="text-input"
    ></textarea>
    <button @click="sendText" class="send-button">发送并处理</button>
    <div class="progress">
      <p v-for="line in processedLines" :key="line">{{ line }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

// 定义输入和处理结果
const textInput = ref('https://localhost\nhttps://localhost\nhttps://localhost\nhttps://localhost\nhttps://localhost\n');
const processedLines = ref([]);

// 发送数据并处理流式返回
const sendText = async () => {
  processedLines.value = []; // 清空之前的处理结果
  if (!textInput.value.trim()) {
    alert('请输入内容后再发送！');
    return;
  }

  try {
    const response = await fetch('/api/process-text', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ text: textInput.value }),
    });

    if (!response.ok) {
      throw new Error(`HTTP 错误: ${response.status}`);
    }

    // 处理后端返回的流式数据
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let result = '';

    while (true) {
      const { value, done } = await reader.read();
      if (done) break; // 读取结束
      result += decoder.decode(value, { stream: true });

      // 按行解析返回的数据
      const lines = result.split('\n');
      result = lines.pop(); // 未完成的一行保留到下次
      processedLines.value.push(...lines.filter(line => line.trim()));
    }
  } catch (error) {
    console.error('处理失败:', error);
    alert('处理失败，请稍后重试！');
  }
};
</script>

<style scoped>
.container {
  max-width: 600px;
  margin: 50px auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

h1 {
  text-align: center;
  margin-bottom: 20px;
}

.text-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  line-height: 1.5;
  margin-bottom: 20px;
}

.send-button {
  display: block;
  width: 100%;
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  text-align: center;
}

.send-button:hover {
  background-color: #369f7d;
}

.progress {
  margin-top: 20px;
  font-size: 14px;
  color: #333;
}

.progress p {
  margin: 5px 0;
}
</style>
