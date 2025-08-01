<script setup lang="ts">
import { context, usernameModal } from "@/assets/store";
</script>

<template>
  <div class="p-4 px-6 bg-gray-800 flex justify-between items-center rounded-lg">
    <h1 class="text-xl font-bold flex items-center">
      Elliptical v0
      <span class="px-2 py-1 rounded-lg text-xs ml-2" :class="{
        'bg-red-500': context.status.code === 2,
        'bg-yellow-500': context.status.code === 1,
        'bg-green-500': context.status.code === 0
      }">
        {{ context.status.text }}
      </span>
    </h1>

    <div class="flex items-center">
      
      <button @click="usernameModal.open = true"
        class="w-full px-5 py-2.5 mr-4 bg-blue-500 hover:bg-blue-600 hover:cursor-pointer rounded-lg text-sm font-bold">
        {{ context.username || 'Set Username' }}
      </button>

      <h3 class="text-nowrap">Users: {{ context.online }}</h3>
    </div>
  </div>
</template>