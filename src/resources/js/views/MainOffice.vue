<template>
    <div class="office-bg">
        <SideDrawer :isOpen="isOpen" @toggle="toggleDrawer" />
        <div class="office-layer">
            <!-- 後でユーザーアイコン等をここに追加していく -->
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue"
import SideDrawer from "../components/SideDrawer.vue"

const isOpen = ref(false)
const toggleDrawer = () => {
  isOpen.value = !isOpen.value
}
</script>

<style lang="scss" scoped src="../../css/mainOffice.scss"></style>