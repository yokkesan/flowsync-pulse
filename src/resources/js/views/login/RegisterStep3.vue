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
                <label for="bio">自己紹介</label>
                <textarea id="bio" v-model="form.bio" rows="4" placeholder="簡単な自己紹介を入力してください"></textarea>
            </div>

            <div class="register-step3__actions">
                <button type="button" class="register-step3__btn--back" @click="backToStep2">← 戻る</button>
                <button type="submit" class="register-step3__btn--submit">登録完了</button>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'

const form = ref({
    github: '',
    bio: ''
})

const submitForm = () => {
    router.visit('/login')
}

const backToStep2 = () => {
    router.visit('/register/step2')
}

import '@/../css/registerStep3.scss'
</script>