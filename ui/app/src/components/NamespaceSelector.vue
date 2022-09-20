<script setup lang="ts">
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import { storeToRefs } from "pinia";
import { Listbox, ListboxButton, ListboxLabel, ListboxOption, ListboxOptions } from "@headlessui/vue";
import { CheckIcon, ChevronUpDownIcon } from "@heroicons/vue/24/solid";
import { useNamespaceStore } from "../stores/namespace";
import { useInitReq } from "../api/init";
import { IdentityService } from "../api/identity/v1/identity.pb";

const nsStore = useNamespaceStore();
const { namespaces } = storeToRefs(nsStore);
const { updateNamespaces } = nsStore;

const route = useRoute();
const initReq = useInitReq();
const identityRes = (await IdentityService.Identity({}, initReq));

updateNamespaces(identityRes.accessibleNamespaces || []);

const namespace = computed(() => {
  return route.params.namespace;
});

const selected = ref(namespaces.value.find(ns => ns === namespace.value));

</script>
<!-- This example requires Tailwind CSS v2.0+ -->
<template>
  <Listbox as="div" v-model="selected" class="flex items-baseline">
    <ListboxLabel class="inline-block mr-2 text-sm font-medium text-gray-300">Namespace</ListboxLabel>
    <div class="relative">
      <ListboxButton class="relative w-full cursor-default rounded-md border border-gray-300 bg-white py-1 pl-3 pr-10 text-left shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 text-sm">
        <span class="block truncate">{{ selected }}</span>
        <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
          <ChevronUpDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
        </span>
      </ListboxButton>

      <transition leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100" leave-to-class="opacity-0">
        <ListboxOptions class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none text-sm">
          <router-link v-for="namespace in namespaces" :key="namespace" :to="{ name: 'dashboard', params: { namespace: namespace } }" custom v-slot="{ navigate }">
            <ListboxOption as="template" :key="namespace" :value="namespace" v-slot="{ active, selected }" @click="navigate">
              <li :class="[active ? 'text-white bg-indigo-600' : 'text-gray-900', 'text-sm relative cursor-default select-none py-1 pl-3 pr-9']">
                <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{ namespace }}</span>
                <span v-if="selected" :class="[active ? 'text-white' : 'text-indigo-600', 'absolute inset-y-0 right-0 flex items-center pr-4']">
                  <CheckIcon class="h-5 w-5" aria-hidden="true" />
                </span>
              </li>
            </ListboxOption>
          </router-link>
        </ListboxOptions>
      </transition>
    </div>
  </Listbox>
</template>
