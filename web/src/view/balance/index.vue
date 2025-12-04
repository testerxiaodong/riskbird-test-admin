<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openBalanceDialog">
          修改用户余额
        </el-button>
      </div>
    </div>

    <!-- 修改余额抽屉 -->
    <el-drawer
      v-model="balanceDialogVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">修改用户余额</span>
          <div>
            <el-button @click="closeBalanceDialog">取 消</el-button>
            <el-button type="primary" @click="submitBalance" :loading="loading">确 定</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="balanceFormRef"
        :model="balanceForm"
        :rules="balanceRules"
        label-width="100px"
      >
        <el-form-item label="手机号" prop="phone" required>
          <el-input v-model="balanceForm.phone" placeholder="请输入用户手机号" />
        </el-form-item>

        <el-form-item label="密码" prop="password" required>
          <el-input 
            v-model="balanceForm.password" 
            type="password" 
            show-password 
            placeholder="请输入用户密码" 
          />
        </el-form-item>

        <el-form-item label="充值金额" prop="rechargeAmount" required>
          <el-input 
            v-model="balanceForm.rechargeAmount" 
            placeholder="请输入修改后的充值金额" 
          />
        </el-form-item>

        <el-form-item label="赠送金额" prop="giftAmount" required>
          <el-input 
            v-model="balanceForm.giftAmount" 
            placeholder="请输入修改后的赠送金额" 
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
import { modifyUserBalance } from '@/api/userBalance'

const appStore = useAppStore()
const balanceDialogVisible = ref(false)
const balanceFormRef = ref()
const loading = ref(false)

const balanceForm = reactive({
  phone: '',
  password: '',
  rechargeAmount: '',
  giftAmount: ''
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

// 数字验证 - 金额字段
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

const balanceRules = {
  phone: [{ validator: phoneValidator, trigger: 'blur' }],
  password: [{ validator: passwordValidator, trigger: 'blur' }],
  rechargeAmount: [{ validator: numberValidator, trigger: 'blur' }],
  giftAmount: [{ validator: numberValidator, trigger: 'blur' }]
}

// 打开余额修改抽屉
const openBalanceDialog = () => {
  balanceForm.phone = ''
  balanceForm.password = ''
  balanceForm.rechargeAmount = ''
  balanceForm.giftAmount = ''
  balanceDialogVisible.value = true
}

// 关闭余额修改抽屉
const closeBalanceDialog = () => {
  balanceFormRef.value?.clearValidate()
  balanceDialogVisible.value = false
}

// 提交余额修改
const submitBalance = async () => {
  try {
    await balanceFormRef.value.validate()
    loading.value = true

    await modifyUserBalance({
      phone: balanceForm.phone,
      password: balanceForm.password,
      rechargeAmount: Number(balanceForm.rechargeAmount),
      giftAmount: Number(balanceForm.giftAmount)
    })

    loading.value = false
    ElMessage.success('用户余额修改成功')
    balanceDialogVisible.value = false
  } catch (error) {
    loading.value = false
  }
}
</script>

<style lang="scss">
</style>
