<template>
  <div class="container">
    <el-card class="form-card">
      <template #header>
        <div style="text-align: center">
          <h2>用户余额修改</h2>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="80px"
        style="max-width: 450px; margin: 0 auto"
      >
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入用户手机号" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入用户密码" />
        </el-form-item>

        <el-form-item label="充值金额" prop="rechargeAmount">
          <el-input v-model="formData.rechargeAmount" placeholder="请输入修改后的充值金额" />
        </el-form-item>

        <el-form-item label="赠送金额" prop="giftAmount">
          <el-input v-model="formData.giftAmount" placeholder="请输入修改后的赠送金额" />
        </el-form-item>

        <div style="text-align: center">
          <el-button type="primary" @click="submitForm" :loading="loading">提交</el-button>
          <el-button @click="resetForm">重置</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { modifyUserBalance } from '@/api/userBalance'

const formRef = ref()
const loading = ref(false)

const formData = reactive({
  phone: '',
  password: '',
  rechargeAmount: '',
  giftAmount: ''
})

// 手机号验证规则
const phoneValidator = (rule, value, callback) => {
  if (!value) {
    callback(new Error('用户手机号不能为空'))
  } else if (!/^1[3-9]\d{9}$/.test(value)) {
    callback(new Error('请输入正确的手机号格式'))
  } else {
    callback()
  }
}

// 密码验证规则
const passwordValidator = (rule, value, callback) => {
  if (!value) {
    callback(new Error('用户密码不能为空'))
  } else {
    callback()
  }
}

// 数字验证规则 - 充值金额和赠送金额都不能为空
const numberValidator = (rule, value, callback) => {
  if (value === '' || value === null || value === undefined) {
    callback(new Error('金额不能为空'))
  } else if (isNaN(value)) {
    callback(new Error('请输入有效的数字'))
  } else if (Number(value) < 0) {
    callback(new Error('金额不能为负数'))
  } else if (!/^\d+(\.\d{1,2})?$/.test(String(value))) {
    callback(new Error('金额最多支持小数点后2位'))
  } else {
    callback()
  }
}

// 表单验证规则
const rules = {
  phone: [{ validator: phoneValidator, trigger: 'blur' }],
  password: [{ validator: passwordValidator, trigger: 'blur' }],
  rechargeAmount: [{ validator: numberValidator, trigger: 'blur' }],
  giftAmount: [{ validator: numberValidator, trigger: 'blur' }]
}

// 提交表单
const submitForm = async () => {
  try {
    await formRef.value.validate()
    loading.value = true

    // 调用API接口
    await modifyUserBalance({
      phone: formData.phone,
      password: formData.password,
      rechargeAmount: Number(formData.rechargeAmount),
      giftAmount: Number(formData.giftAmount)
    })
    
    loading.value = false
    ElMessage.success('用户余额修改成功！')
    resetForm()

  } catch (error) {
    loading.value = false
    // 错误信息已由拦截器显示，这里不需要再显示
  }
}

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
  formData.phone = ''
  formData.password = ''
  formData.rechargeAmount = ''
  formData.giftAmount = ''
}
</script>

<style scoped>
.container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  background: #f5f7fa;
  padding: 60px 20px 20px;
  overflow-y: auto;
}

.form-card {
  width: 100%;
  max-width: 380px;
  flex-shrink: 0;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

:deep(.el-card) {
  border-radius: 12px;
}

:deep(.el-card__header) {
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
  border-radius: 12px 12px 0 0;
}

:deep(.el-card__body) {
  border-radius: 0 0 12px 12px;
}

h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

:deep(.el-input__wrapper) {
  border-radius: 6px;
}

:deep(.el-input) {
  border-radius: 6px;
}

:deep(.el-button) {
  border-radius: 6px;
}
</style>
