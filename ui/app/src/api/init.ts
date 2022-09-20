import { inject } from "vue";
import { useRoute } from "vue-router";
import { InitReq } from "./fetch.pb";

export const useInitReq = (): InitReq => {
  const route = useRoute();

  const req: InitReq = {
    pathPrefix: inject('apiPathPrefix'),
  };

  if (route.params.namespace !== "") {
    req.headers = {
      "Grpc-Metadata-namespace": route.params.namespace as string,
    }
  }

  return req;
}
