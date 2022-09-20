<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useNamespaceStore } from "../stores/namespace";
import { useInitReq } from "../api/init";
import { IdentityService } from "../api/identity/v1/identity.pb";

const nsStore = useNamespaceStore();
const { namespaces } = storeToRefs(nsStore);
const { updateNamespaces } = nsStore;

const initReq = useInitReq();
const identityRes = (await IdentityService.Identity({}, initReq));

updateNamespaces(identityRes.accessibleNamespaces || []);
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <div class="col-span-1 md:col-span-3 bg-white rounded-lg shadow px-5 py-6 sm:px-6">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div v-for="namespace in namespaces" :key="namespace"
            class="rounded-lg border border-gray-300 bg-white px-6 py-5 shadow-sm focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 hover:border-gray-400">
            <div class="min-w-0 flex-1">
              <div class="flex align-baseline">
                <p class="grow text-lg font-bold text-gray-900">{{ namespace }}</p>
                <router-link :to="{ name: 'dashboard', params: { namespace: namespace } }" class="focus:outline-none">
                  <button type="button" class="inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-2 py-1 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">Select</button>
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
