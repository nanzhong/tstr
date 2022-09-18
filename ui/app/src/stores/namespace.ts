import { ref, Ref } from "vue";
import { defineStore } from "pinia";

export const useNamespaceStore = defineStore('namespace', () => {
  const namespaces: Ref<string[]> = ref([]);
  const currentNamespace: Ref<string> = ref("");

  function updateNamespaces(newNamespaces: string[]) {
    namespaces.value = newNamespaces;
  }

  function setCurrentNamespace(ns: string) {
    currentNamespace.value = ns;
    console.log(`updated current namespace to ${currentNamespace.value}`);
  }

  return { namespaces, currentNamespace, updateNamespaces, setCurrentNamespace }
})
