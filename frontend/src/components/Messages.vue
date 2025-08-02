<script setup lang="ts">
import DOMPurify from 'dompurify';
import { context, currentRoom } from "@/assets/store";

</script>

<template>
  <div v-if="currentRoom.roomid" class="flex-1 p-4 overflow-y-auto rounded-lg">
    <h3 class="font-bold">Welcome to #{{ currentRoom.title }}</h3>

    <ul>
      <li v-for="(message, index) in currentRoom.messages" :key="index">
        <div class="flex items-center justify-between m-2">
          <div v-html="DOMPurify.sanitize(message.msg)" :id="message.id"></div>

        </div>
      </li>
    </ul>
  </div>

  <div v-else class="flex-1 p-4 rounded-lg flex flex-col items-center justify-center text-gray-400">
    <h1>Welcome to Elliptical</h1>
    <h3>{{ context.rooms.length === 0 ? 'Create' : 'Select' }} a room to start chatting...</h3>
  </div>
</template>