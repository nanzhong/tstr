{
  description = "tstr devleopment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    pggen = {
      url = "github:jschaf/pggen";
      flake = false;
    };
    dbmate = {
      url = "github:amacneil/dbmate";
      flake = false;
    };
    overmind = {
      url = "github:DarthSim/overmind";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, pggen, dbmate, overmind }:
    flake-utils.lib.eachDefaultSystem(system:
      let
        pkgs = import nixpkgs { inherit system; };
        devTools = {
          pggen = pkgs.buildGoModule {
            name = "pggen";
            src = pggen;
            subPackages = [ "cmd/pggen" ];
            doCheck = false;
            vendorSha256 = "sha256-WLoFpwOP97160WfmfbCUUlhqGC0qiEPWDg0qL/DrzIA=";
          };
          dbmate = pkgs.buildGoModule {
            name = "dbmate";
            src = dbmate;
            subPackages = [ "." ];
            doCheck = false;
            vendorSha256 = "sha256-U9VTS0rmLHxweFiIcFyoybHMBihy5ezloDC2iLc4IMc=";
          };
          overmind = pkgs.buildGoModule {
            name = "overmind";
            src = overmind;
            subPackages = [ "." ];
            dbCheck = false;
            vendorSha256 = "sha256-KDMzR6qAruscgS6/bHTN6RnHOlLKCm9lxkr9k3oLY+Y=";
          };
        };
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

                postgresql_14

                devTools.pggen
                devTools.dbmate
                devTools.overmind
              ];
            };
        }
    );
}
