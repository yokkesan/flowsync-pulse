<template>
  <div class="register-step1">
    <h1 class="register-step1__title">新規登録（Step 1）</h1>

    <form @submit.prevent="goNext">
      <div class="register-step1__field">
        <label for="name">ユーザー名</label>
        <input id="name" v-model="form.name" type="text" required />
        <p v-if="errors.name" class="register-step1__error">{{ errors.name[0] }}</p>
      </div>

      <div class="register-step1__field">
        <label for="email">メールアドレス</label>
        <input id="email" v-model="form.email" type="email" required />
        <p v-if="errors.email" class="register-step1__error">{{ errors.email[0] }}</p>
      </div>

      <div class="register-step1__field">
        <label for="password">パスワード</label>
        <input id="password" v-model="form.password" type="password" required />
        <p v-if="errors.password" class="register-step1__error">{{ errors.password[0] }}</p>
      </div>

      <button type="submit" class="register-step1__btn" :disabled="loading">
        {{ loading ? '送信中...' : '次へ' }}
      </button>
    </form>

    <p class="register-step1__back">
      <a href="/" @click.prevent="backToTop">← 戻る</a>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { router } from '@inertiajs/vue3'

const form = ref({
  name: '',
  email: '',
  password: ''
})

const errors = ref({})
const loading = ref(false)

const goNext = async () => {
  errors.value = {}
  loading.value = true

  try {
    const res = await axios.post('/register/step1', form.value)

    // ✅ コントローラから返る nextStep に遷移
    if (res.data?.nextStep) {
      router.visit(res.data.nextStep)
    }
  } catch (e) {
    if (e.response?.status === 422) {
      errors.value = e.response.data.errors
    } else {
      console.error(e)
    }
  } finally {
    loading.value = false
  }
}

const backToTop = () => {
  router.visit('/')
}
</script>