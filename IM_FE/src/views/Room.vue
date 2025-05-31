<template>
  <el-container class="main">
    <el-aside width="30%" style="border-right: 1px solid #ccc">
      <div class="asideHeader">聊天列表</div>
      <div class="aside">
        <div
          class="chat-item"
          v-for="(chat, index) in chatList"
          :key="index"
          @click="handSelectedChat(index)"
        >
          <div>{{ chat.name }}</div>
          <div>{{ chat.host }}</div>
        </div>
      </div>
    </el-aside>
    <el-container>
      <el-header class="header">{{ selectedChat?.name }}</el-header>
      <el-main>
        <div class="chat-view">
          <div
            v-for="(chatMsg, index) in messageList"
            :key="index"
            :class="chatMsg.otherMsg ? 'other-msg' : 'self-msg'"
            class="chat-msg"
          >
            <div class="chat-avatar">{{ chatMsg.id }}</div>
            <div class="chat-message">{{ chatMsg.content }}</div>
          </div>
        </div>
        <div class="chat-input">
          <el-input v-model="content" placeholder="输入文本" />
          <el-button type="success" @click="sendMessage">发送</el-button>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { getGroupsInfoMsg, getTextMsg } from '@/utils/messageHandler';
import { ref, onMounted } from 'vue';

const chatList = ref([]);
const selectedChat = ref(null);

const content = ref('');

const messageList = ref([]);

const handSelectedChat = (index) => {
  selectedChat.value = chatList.value[index];
};

let ws = null;

onMounted(() => startWS());

const startWS = () => {
  ws = new WebSocket('ws://localhost:8080/ws');

  ws.onopen = () => {
    ElMessage.success('连接成功');
  };

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      console.log(message);

      switch (message.type) {
        case 'groupsInfo':
          chatList.value = getGroupsInfoMsg(message);
          break;

        case 'text':
          messageList.value.push(getTextMsg(message));
          break;

        case 'image':
          break;

        case 'audio':
          break;

        case 'video':
          break;

        default:
          ElMessage.error('消息类型出错');
          break;
      }
    } catch (err) {
      ElMessage.error('消息接收出错');
      console.log(err);
    }
  };

  ws.onerror = () => {
    ElMessage.error('连接失败');
  };

  ws.onclose = () => {
    ElMessage.error('连接关闭');
  };
};

const sendMessage = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(
      JSON.stringify({
        type: 'text',
        target: selectedChat.value.id,
        content: {
          text: content.value,
        },
        toGroup: true,
        // createdAt: Date.now(),
      }),
    );
    content.value = '';
  } else {
    ElMessage.error('没有连接成功 ws!');
    return;
  }
};
</script>

<style scoped>
.main {
  border: 1px solid #ccc;
  width: 70vw;
  height: 80vh;
  margin: 5vh auto;
  display: flex;
}

.asideHeader {
  text-align: center;
  height: 5%;
  border-bottom: 1px solid #ccc;
}

.aside {
  height: 94%;
  overflow-y: auto;
}

.header {
  border-bottom: 1px solid #ccc;
}

.chat-item {
  width: 100%;
  padding: 10px;
  box-sizing: border-box;
  border-bottom: 1px solid #ccc;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
}

.chat-view {
  height: 80%;
  border-bottom: 1px solid #ccc;
  margin-bottom: 10px;
  overflow-y: auto;
}

.chat-input {
  display: flex;
  justify-content: center;
  gap: 10px;
}

.chat-msg {
  display: flex;
  align-items: flex-start;
  margin: 10px;
  padding: 10px;
  border-radius: 10px;
}

.self-msg {
  justify-content: right;
  background-color: #e1f5fe;
  flex-direction: row-reverse;
}

.other-msg {
  justify-content: flex-start;
  background-color: #f5f5f5;
}

.chat-avatar {
  width: 40px;
  height: 40px;
  border: 1px solid #ccc;
  border-radius: 50%;
}

.chat-message {
  max-width: 70%;
  padding: 10px;
  border-radius: 10px;
  word-wrap: break-word;
}

/* 响应式设计 */
@media screen and (max-width: 600px) {
  .main {
    flex-direction: column;
    width: 90vw;
  }

  .aside {
    width: 100%;
    height: auto;
  }

  .chat-view {
    height: 60vh;
  }
}
</style>
