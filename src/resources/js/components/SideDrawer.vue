<template>
    <div :class="['drawer', { open: isOpen }]">
        <ul class="menu">
            <li>メインページ</li>
            <li>チャットルーム</li>
            <li>プロジェクトルーム</li>
            <li>管理メニュー</li>
        </ul>
    </div>

    <!-- トグルボタン -->
    <button class="drawer-toggle" @click="$emit('toggle')">
        ☰
    </button>
</template>

<script setup>
defineProps({
    isOpen: {
        type: Boolean,
        default: false,
    },
})
</script>

<style scoped lang="scss" src="../../css/sideDrawer.scss"></style>