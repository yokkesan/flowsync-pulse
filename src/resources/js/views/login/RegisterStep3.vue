<template>
    <div class="register-step3">
        <h1 class="register-step3__title">新規登録（Step 3）</h1>
        <p class="register-step3__desc">
            任意情報を入力してください。（スキップ可能）
        </p>

        <form @submit.prevent="submitForm">
            <div class="register-step3__field">
                <label for="github">GitHub アカウント</label>
                <input id="github" v-model="form.github" type="text" placeholder="@username" />
            </div>

            <div class="register-step3__field">
                <label for="repository">アサインしている GitHub リポジトリ</label>
                <input id="repository" v-model="form.repository" type="text"
                    placeholder="example-org/repository-name" />
            </div>

            <div class="register-step3__actions">
                <button type="button" class="register-step3__btn--back" @click="backToStep2" :disabled="loading">
                    ← 戻る
                </button>
                <button type="submit" class="register-step3__btn--submit" :disabled="loading">
                    {{ loading ? '登録中...' : '登録完了' }}
                </button>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { router } from '@inertiajs/vue3'

const form = ref({
    github: '',
    repository: ''
})

const loading = ref(false)

const submitForm = async () => {
    loading.value = true
    try {
        const res = await axios.post('/register/step3', form.value)
        if (res.data?.nextStep) {
            // ✅ Inertia構成に統一：router.visitで遷移
            alert('登録が完了しました。ログイン画面へ移動します。')
            router.visit(res.data.nextStep)
        }
    } catch (error) {
        console.error(error)
        alert('登録処理中にエラーが発生しました。')
    } finally {
        loading.value = false
    }
}

const backToStep2 = () => {
    router.visit('/register/step2')
}

import '@/../css/registerStep3.scss'
</script>