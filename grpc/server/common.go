package server

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nanzhong/tstr/grpc/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func namespaceFromContext(ctx context.Context) (string, error) {
	ns, err := auth.NamespaceFromContext(ctx)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, "request metadata missing namespace")
	}
	return ns, nil
}

func matchNamespace(nsSel []string, ns string) (bool, error) {
	for _, nsSel := range nsSel {
		re, err := regexp.Compile(nsSel)
		if err != nil {
			return false, fmt.Errorf("matching namespace: %w", err)
		}
		if re.MatchString(ns) {
			return true, nil
		}
	}
	return false, nil
}
