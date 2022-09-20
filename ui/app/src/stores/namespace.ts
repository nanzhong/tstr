import { ref, Ref } from "vue";
import { defineStore } from "pinia";

export const useNamespaceStore = defineStore('namespace', () => {
  const namespaces: Ref<string[]> = ref([]);

  function updateNamespaces(newNamespaces: string[]) {
    namespaces.value = newNamespaces;
  }

  return { namespaces, updateNamespaces }
})
