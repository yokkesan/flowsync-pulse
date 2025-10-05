<template>
    <div :class="['drawer', { open: isOpen }]">
        <ul class="menu">
            <li>メインページ</li>
            <li>チャットルーム</li>
            <li>プロジェクトルーム</li>
            <li>管理メニュー</li>
        </ul>
    </div>

    <!-- ハンバーガー -->
    <button class="drawer-toggle" :class="{ 'drawer-toggle--shifted': isOpen }" @click="$emit('toggle')">
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