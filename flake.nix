{
  description = "tstr devleopment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem(system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
        {
          devShell = with pkgs;
            mkShell {
              buildInputs = [
                go_1_18
                go-tools
                gopls
                protobuf
                protoc-gen-go
                protoc-gen-go-grpc
              ];
            };
        }
    );
}
