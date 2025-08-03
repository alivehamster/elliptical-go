<script setup lang="ts">
import { context, newRoom, currentRoom } from '@/assets/store';
import { sendMessage } from '@/assets/socket';

function joinRoom(room: { title: string; roomid: string }) {
  currentRoom.title = room.title;
  sendMessage({ type: 'JoinRoom', string: room.roomid });
}

</script>

<template>
  <div class="w-64 bg-gray-800 rounded-lg">
    <div class="p-4 overflow-y-auto overflow-x-hidden"
      :class="context.rooms.length === 0 ? '' : 'h-[calc(100%-6.5rem)]'">

      <details :open="context.rooms.length !== 0">
        <summary class="font-bold cursor-pointer">
          Public Rooms
        </summary>

        <ul class="h-full scrollbar-transparent" v-if="context.rooms.length !== 0">
          <li v-for="(room, index) in context.rooms" :key="index" class="my-2 rounded-lg">
            <button @click="joinRoom(room)"
              class="w-full px-4 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-600 hover:cursor-pointer">
              {{ room.title }}
            </button>
          </li>
        </ul>
      </details>
    </div>

    <p v-if="context.rooms.length === 0" class="ml-4 text-gray-300 text-sm">
      No active rooms, create one below
    </p>

    <div class="absolute bottom-0 left-0 right-0 p-2 w-64 bg-gray-800 rounded-lg">
      <button @click="newRoom.open = true"
        class="w-full px-5 py-2.5 bg-green-500 hover:bg-green-600 rounded-lg text-sm font-bold hover:cursor-pointer">
        Create Room
      </button>
    </div>
  </div>
</template>