<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { SparklesIcon, MenuIcon, XIcon } from '@heroicons/vue/outline';

const navigation = [
  { name: 'Dashboard', route: 'dashboard' },
  { name: 'Tests', route: 'tests' },
  { name: 'Runs', route: 'runs' },
  { name: 'Runners', route: 'runners' },
];

const currentRoute = computed((path) => {
  const route = useRoute();
  return route;
});

</script>

<template>
  <div class="min-h-full">
    <div class="bg-gray-800 pb-32">
      <Disclosure as="nav" class="bg-gray-800" v-slot="{ open }">
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
          <div class="border-b border-gray-700">
            <div class="flex items-center justify-between h-16 px-4 sm:px-0">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <SparklesIcon class="h-8 w-8 text-white" />
                </div>
                <div class="flex-shrink-0">
                  <span class="ml-4 text-white text-xl font-extrabold">tstr</span>
                </div>
                <div class="hidden md:block">
                  <div class="ml-10 flex items-baseline space-x-4">
                    <router-link v-for="item in navigation" :key="item.route" :to="{ name: item.route }" custom v-slot="{ href, isActive, navigate }">
                      <a :href="href" @click="navigate"
                        :class="[isActive ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white', 'px-3 py-2 rounded-md text-sm font-medium']"
                        :aria-current="isActive ? 'page' : undefined">{{ item.name }}</a>
                    </router-link>
                  </div>
                </div>
              </div>
              <div class="-mr-2 flex md:hidden">
                <!-- Mobile menu button -->
                <DisclosureButton
                  class="bg-gray-800 inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white">
                  <span class="sr-only">Open main menu</span>
                  <MenuIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
                  <XIcon v-else class="block h-6 w-6" aria-hidden="true" />
                </DisclosureButton>
              </div>
            </div>
          </div>
        </div>

        <DisclosurePanel class="border-b border-gray-700 md:hidden">
          <div class="px-2 py-3 space-y-1 sm:px-3">
            <router-link v-for="item in navigation" :key="item.name" :to="{ name: item.route }" custom
              v-slot="{ href, navigate, isActive }">
              <DisclosureButton as="a" :href="href"
                :class="[isActive ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white', 'block px-3 py-2 rounded-md text-base font-medium']"
                :aria-current="isActive ? 'page' : undefined" @click="navigate">{{ item.name }}</DisclosureButton>
            </router-link>
          </div>
        </DisclosurePanel>
      </Disclosure>
      <header class="py-10">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <router-view name="header"></router-view>
        </div>
      </header>
    </div>

    <main class="-mt-32">
      <router-view v-slot="{ Component, route }">
        <template v-if="Component">
          <Transition name="fade">
            <Suspense>
              <component :is="Component" :key="route.path" />

              <template #fallback>
                <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
                  <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
                    <div class="flex justify-center items-center">
                      <div role="status">
                        <svg aria-hidden="true"
                          class="mr-2 w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600"
                          viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
                          <path
                            d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                            fill="currentColor" />
                          <path
                            d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                            fill="currentFill" />
                        </svg>
                      </div>
                      <span>Loading...</span>
                    </div>
                  </div>
                </div>
              </template>
            </Suspense>
          </Transition>
        </template>
      </router-view>
    </main>
  </div>
</template>
