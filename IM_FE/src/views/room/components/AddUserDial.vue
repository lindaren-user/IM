<template>
  <el-dialog
    :model-value="props.dialogVisible"
    title="添加好友"
    width="500"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-input
      v-model="keyword"
      style="max-width: 600px"
      placeholder="请输入关键字"
      class="input-with-select"
    >
      <template #prepend>
        <el-select v-model="select" placeholder="账号" style="width: 115px">
          <el-option label="账号" value="0" />
          <el-option label="昵称" value="1" />
        </el-select>
      </template>
      <template #append>
        <el-button :icon="Search" @click="search" :disabled="!keyword" />
      </template>
    </el-input>

    <div class="userList" v-loading="loading">
      <div class="userInfo" v-for="(user, index) in userList" :key="user.id">
        {{ user.nickname }}
      </div>
      <el-empty description="暂无结果" v-if="!userList || userList.length === 0" />
    </div>

    <template #footer>
      <div class="dialog-footer"></div>
    </template>
  </el-dialog>
</template>

<script setup>
import { user_service } from '@/apis/user_service';
import { Search } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { ref } from 'vue';

let controller = null; // 用于中断请求

const props = defineProps({
  dialogVisible: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:dialogVisible']);
const handleClose = () => {
  emit('update:dialogVisible', false);
  if (controller) controller.abort();
};

const loading = ref(false);

const select = ref('0');
const keyword = ref('');

const userList = ref([]);

const search = () => {
  // 每次请求前，先取消上一次请求
  if (controller) controller.abort();
  controller = new AbortController(); // 每次请求都是新的实例

  loading.value = true;

  user_service
    .search(select.value, keyword.value, controller.signal)
    .then((res) => {
      userList.value = res.data.data;
      ElMessage.success('查询成功');
    })
    .catch((err) => {
      // 判断是否为取消错误
      if (err.name === 'AbortError') {
        ElMessage.info('请求已取消');
      } else {
        ElMessage.error(err.message);
      }
    })
    .finally(() => (loading.value = false));
};
</script>

<style scoped lang="scss">
.userList {
  border: 1px solid lightgray;
  margin-top: 5%;
  height: 40vh;
  display: flex;
  flex-direction: column;
  gap: 5px;
  overflow-y: auto;
  // padding: 10px;

  .userInfo {
    border: 1px solid lightgray;
    display: flex;
    justify-content: space-between;
    height: 15%;
  }
}
</style>
