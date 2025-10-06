<template>
    <div class="register-step2">
        <h1 class="register-step2__title">新規登録（Step 2）</h1>
        <p class="register-step2__desc">
            ご登録のメールアドレスに送信された認証コードを入力してください。
        </p>

        <form @submit.prevent="goNext">
            <div class="register-step2__field">
                <label for="code">認証コード</label>
                <input id="code" v-model="form.code" type="text" required maxlength="6" />
                <p v-if="errors.code" class="register-step2__error">{{ errors.code[0] }}</p>
            </div>

            <button type="submit" class="register-step2__btn" :disabled="loading">
                {{ loading ? '確認中...' : '次へ' }}
            </button>
        </form>

        <p class="register-step2__back">
            <a href="/register/step1" @click.prevent="backToStep1">← 戻る</a>
        </p>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { router } from '@inertiajs/vue3'

const form = ref({ code: '' })
const errors = ref({})
const loading = ref(false)

const goNext = async () => {
    errors.value = {}
    loading.value = true

    try {
        const res = await axios.post('/register/step2', form.value)

        if (res.data?.nextStep) {
            router.visit(res.data.nextStep) // ✅ Inertia経由で遷移（Step3へ）
        }
    } catch (e) {
        if (e.response?.status === 422) {
            // Laravelバリデーションエラー対応
            errors.value = e.response.data.errors
        } else {
            console.error(e)
            alert('認証コードが無効、または期限切れです。')
        }
    } finally {
        loading.value = false
    }
}

const backToStep1 = () => {
    router.visit('/register/step1')
}

import '@/../css/registerStep2.scss'
</script>