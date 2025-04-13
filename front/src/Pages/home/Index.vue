<script setup lang="ts">
import Foo from "../../components/Foo.vue"

defineProps<{ name: string }>()
</script>

<template>
  <h1>Hello, {{ name }}</h1>
  <Foo />
</template>
n
