<template>
  <el-container class="login-container">
    <el-form class="login-form" :model="form" :rules="rules">
      <div class="header">登录</div>
      <el-form-item prop="username">
        <el-input v-model="form.username" placeholder="请输入用户名" />
      </el-form-item>

      <el-form-item prop="password">
        <el-input v-model="form.password" type="password" placeholder="请输入密码" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleLogin" class="login-button"> 登录 </el-button>
      </el-form-item>
    </el-form>
  </el-container>
</template>

<script setup>
import { user_service } from '@/apis/user_service';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const form = ref({
  username: '',
  password: '',
});

const rules = ref({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
});

const handleLogin = () => {
  user_service
    .login(form.value)
    .then((res) => {
      ElMessage.success(res.data.message);
      localStorage.setItem('token', res.data.data);
      router.push('/room');
    })
    .catch((err) => {
      ElMessage.error(err.message);
    });
};
</script>

<style scoped lang="scss">
.login-container {
  display: flex;
  justify-content: center;
  margin-top: 15vh;

  .login-form {
    width: 20%;
    padding: 24px;
    border-radius: 8px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);

    .header {
      text-align: center;
      padding: 10px;
    }

    .login-button {
      width: 100%;
      margin-top: 16px;
      padding: 12px;
      background-color: #409eff;
      color: white;
      border-radius: 4px;
      transition: background-color 0.3s;
    }

    .login-button:hover {
      background-color: #2b6bd3;
    }
  }
}

.el-form-item {
  margin-bottom: 16px;
}

.el-input {
  border-radius: 4px;
  height: 40px;
}

.el-form-item__label {
  font-weight: 500;
  color: #303133;
}
</style>
