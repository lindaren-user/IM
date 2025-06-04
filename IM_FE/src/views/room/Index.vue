<template>
  <el-container class="main">
    <el-aside width="30%" style="border-right: 1px solid #ccc">
      <div class="asideHeader">
        <text>聊天列表</text>
        <img src="@/assets/imgs/add.svg" class="add-icon" @click="openAddUser = true" />
      </div>
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

  <AddUserDial :dialogVisible="openAddUser" @update:dialogVisible="openAddUser = false" />
</template>

<script setup>
import { ref, onMounted } from 'vue';
import AddUserDial from './components/addUserDial.vue';

// 聊天列表
const chatList = ref([]);

// 选择的消息
const selectedChat = ref(null);

// 消息内容
const content = ref('');

// 消息列表
const messageList = ref([]);

// 是否打开添加好友弹窗
const openAddUser = ref(false);

// 处理选择的聊天
const handSelectedChat = (index) => {
  selectedChat.value = chatList.value[index];
};

let ws = null;

onMounted(() => startWS());

const startWS = () => {
  const token = localStorage.getItem('token');
  ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

  ws.onopen = () => {
    ElMessage.success('连接成功');
  };

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      messageList.push(message);
      console.log(message);
    } catch (err) {
      ElMessage.error('消息接收出错');
      console.log(err);
    }
  };

  ws.onerror = (err) => {
    ElMessage.error('连接失败');
    console.log(err);
  };

  ws.onclose = () => {
    ElMessage.error('连接关闭');
  };
};

const sendMessage = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(
      JSON.stringify({
        to_id: 1,
        chat_type: 'private',
        content_type: 'text',
        content: content.value,
        create_at: Date.now(),
      }),
    );
    content.value = '';
  } else {
    ElMessage.error('没有连接成功 ws!');
    return;
  }
};
</script>

<style lang="scss" scoped>
.main {
  border: 1px solid #ccc;
  width: 70vw;
  height: 80vh;
  margin: 5vh auto;
  display: flex;

  .asideHeader {
    text-align: center;
    height: 5%;
    border-bottom: 1px solid #ccc;
    display: flex;
    align-items: center;

    text {
      flex: 1;
    }

    .add-icon {
      width: 2rem;
      cursor: pointer;
    }
  }

  .aside {
    height: 94%;
    overflow-y: auto;
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

  .header {
    border-bottom: 1px solid #ccc;
  }

  .chat-view {
    height: 80%;
    border-bottom: 1px solid #ccc;
    margin-bottom: 10px;
    overflow-y: auto;

    .chat-msg {
      display: flex;
      align-items: flex-start;
      margin: 10px;
      padding: 10px;
      border-radius: 10px;

      &.self-msg {
        justify-content: right;
        background-color: #e1f5fe;
        flex-direction: row-reverse;
      }

      &.other-msg {
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
    }
  }

  .chat-input {
    display: flex;
    justify-content: center;
    gap: 10px;
  }
}

/* 响应式设计 */
@media screen and (max-width: 600px) {
  .main {
    flex-direction: column;
    width: 90vw;

    .aside {
      width: 100%;
      height: auto;
    }

    .chat-view {
      height: 60vh;
    }
  }
}
</style>
