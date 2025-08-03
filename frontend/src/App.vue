<script setup lang="ts">
import { onMounted, ref } from "vue";
import { createWebSocket, sendMessage } from "@/assets/socket";
import { context, currentRoom } from "./assets/store";
import Navbar from "@/components/Navbar.vue";
import usernamePicker from "@/components/modals/usernamePicker.vue";
import Messages from "./components/Messages.vue";
import Rooms from "./components/Rooms.vue";
import roomCreation from "./components/modals/roomCreation.vue";
import type { Message } from "@/assets/types";

createWebSocket();
const chat = ref('');

onMounted(() => {
  const username = localStorage.getItem("username") || `Guest ${Math.floor(1000 + Math.random() * 9000)}`;
  localStorage.setItem("username", username);
  context.username = username;
});

function sendChat() {
  if (chat.value.trim() === "") return;
  if (!currentRoom.roomid) {
    console.warn("No room selected to send chat.");
    return;
  }

  const message: Message = {
    type: "SendChat",
    chat: {
      id: currentRoom.roomid,
      msg: `${context.username}: ${chat.value}`,
    },
  };

  sendMessage(message);
  chat.value = "";
}
</script>
<template>
  <div id="main" class="fade-in">
    <roomCreation />

    <usernamePicker />

    <div class="flex flex-col h-screen bg-gray-900 text-white p-4">
      <Navbar />

      <div class="flex flex-1 overflow-hidden mt-4 relative">
        <Rooms />
        <Messages />
      </div>

      <div class="p-4 mt-2 bg-gray-800 rounded-lg text-center" v-if="currentRoom.roomid">
        <form @submit.prevent="sendChat" class="flex items-center">
          <input v-model="chat" autocomplete="off" placeholder="Message" required
            class="flex-grow py-2.5 px-5 text-sm font-medium rounded-lg border border-gray-400 bg-gray-600 outline-none" />
          <button class="px-5 py-2.5 ml-2 bg-blue-500 hover:bg-blue-600 rounded-lg text-sm font-bold">Send</button>
        </form>
      </div>

      <!-- <Admin v-if="path === '/admin'" /> -->
    </div>
  </div>
</template>

<style>
@keyframes pulse {
  0% {
    color: white;
  }

  50% {
    color: gray;
  }

  100% {
    color: white;
  }
}

.animate-pulse {
  animation: pulse 1s infinite;
}

.fade-out {
  transition: opacity 0.5s;
  opacity: 0;
}

.fade-in {
  transition: opacity 0.5s;
  opacity: 1;
}
</style>