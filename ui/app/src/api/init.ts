import { inject } from "vue";
import { storeToRefs } from "pinia";
import { useNamespaceStore } from "../stores/namespace";
import { InitReq } from "./fetch.pb";

export const useInitReq = (): InitReq => {
  const nsStore = useNamespaceStore();
  const { currentNamespace } = storeToRefs(nsStore);

  const req: InitReq = {
    pathPrefix: inject('apiPathPrefix'),
  };

  console.log(`current namespace: ${currentNamespace.value}`);

  if (currentNamespace.value !== "") {
    req.headers = {
      "Grpc-Metadata-namespace": currentNamespace.value,
    }
  }

  return req;
}
