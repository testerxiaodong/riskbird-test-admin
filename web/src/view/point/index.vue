<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openPointDialog">
          修改用户积分
        </el-button>
      </div>
    </div>

    <!-- 修改积分抽屉 -->
    <el-drawer
      v-model="pointDialogVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">修改用户积分</span>
          <div>
            <el-button @click="closePointDialog">取 消</el-button>
            <el-button type="primary" @click="submitPoint" :loading="loading">确 定</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="pointFormRef"
        :model="pointForm"
        :rules="pointRules"
        label-width="100px"
      >
        <el-form-item label="手机号" prop="phone" required>
          <el-input v-model="pointForm.phone" placeholder="请输入用户手机号" />
        </el-form-item>

        <el-form-item label="密码" prop="password" required>
          <el-input 
            v-model="pointForm.password" 
            type="password" 
            show-password 
            placeholder="请输入用户密码" 
          />
        </el-form-item>

        <el-form-item label="修改积分" prop="pointAmount" required>
          <el-input 
            v-model="pointForm.pointAmount" 
            placeholder="请输入修改后的积分数（必须是5的倍数）" 
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useAppStore } from '@/pinia'
import { modifyUserPoint } from '@/api/userPoint'

const appStore = useAppStore()
const pointDialogVisible = ref(false)
const pointFormRef = ref()
const loading = ref(false)

const pointForm = reactive({
  phone: '',
  password: '',
  pointAmount: ''
})

// 手机号验证
const phoneValidator = (rule, value, callback) => {
  if (!value) {
    callback(new Error('用户手机号不能为空'))
  } else if (!/^1[3-9]\d{9}$/.test(value)) {
    callback(new Error('请输入正确的手机号格式'))
  } else {
    callback()
  }
}

// 密码验证
const passwordValidator = (rule, value, callback) => {
  if (!value) {
    callback(new Error('用户密码不能为空'))
  } else {
    callback()
  }
}

// 积分验证
const pointValidator = (rule, value, callback) => {
  if (value === '' || value === null || value === undefined) {
    callback(new Error('积分不能为空'))
  } else if (isNaN(value)) {
    callback(new Error('请输入有效的数字'))
  } else if (Number(value) < 0) {
    callback(new Error('积分不能为负数'))
  } else if (Number(value) % 5 !== 0) {
    callback(new Error('积分必须是5的倍数'))
  } else if (!Number.isInteger(Number(value))) {
    callback(new Error('积分必须是整数'))
  } else {
    callback()
  }
}

const pointRules = {
  phone: [{ validator: phoneValidator, trigger: 'blur' }],
  password: [{ validator: passwordValidator, trigger: 'blur' }],
  pointAmount: [{ validator: pointValidator, trigger: 'blur' }]
}

// 打开积分修改对话框
const openPointDialog = () => {
  pointForm.phone = ''
  pointForm.password = ''
  pointForm.pointAmount = ''
  pointDialogVisible.value = true
}

// 关闭积分修改对话框
const closePointDialog = () => {
  pointFormRef.value?.clearValidate()
  pointDialogVisible.value = false
}

// 提交积分修改
const submitPoint = async () => {
  try {
    await pointFormRef.value.validate()
    loading.value = true

    await modifyUserPoint({
      phone: pointForm.phone,
      password: pointForm.password,
      pointAmount: Number(pointForm.pointAmount)
    })

    loading.value = false
    ElMessage.success('用户积分修改成功')
    pointDialogVisible.value = false
  } catch (error) {
    loading.value = false
  }
}
</script>

<style lang="scss">
</style>
