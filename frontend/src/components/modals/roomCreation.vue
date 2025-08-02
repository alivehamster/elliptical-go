<script setup lang="ts">
import { ref } from "vue";
import { newRoom } from "@/assets/store";
import { sendMessage } from "@/assets/socket";

const roomTitle = ref("");

function createRoom() {
  if (roomTitle.value.trim() === "") {
    return;
  }
  newRoom.open = false;

  sendMessage({
    type: "CreateRoom",
    string: roomTitle.value,
  });

  roomTitle.value = "";
}
</script>

<template>
  <!--Room creation modal-->
  <div v-if="newRoom.open"
    class="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 bottom-0 z-40 flex justify-center items-center w-full md:inset-0 bg bg-gray-600/50">
    <div class="relative w-full max-w-2xl max-h-full bg-gray-800 rounded-lg shadow p-6 text-white text-center">
      <h3 class="text-xl font-semibold">
        Create a Room
      </h3>
      <div class="py-4 font-semibold text-gray-400 align-middle">
        <input placeholder="Room name" type="text" maxlength="25"
          class="w-80 py-2.5 px-5 m-1 text-sm font-medium rounded-lg border border-gray-400 bg-gray-600 outline-none text-white"
          v-model="roomTitle" />
        <div class="flex justify-center items-center mt-2 mb-1">
          <label>Private</label>
         
        </div>

      </div>
      <div class="flex justify-center">
        <button data-modal-hide="default-modal" type="button"
          class="px-5 py-2.5 bg-blue-500 hover:bg-blue-600 rounded-lg text-sm font-bold"
          @click="createRoom">Create</button>
      </div>
    </div>
  </div>
</template>