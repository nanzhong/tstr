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

        tstr = pkgs.buildGoModule {
          name = "tstr";
          src = ./.;
          subPackages = [
            "cmd/tstr"
            "cmd/tstrctl"
          ];
          vendorSha256 = "sha256-9q36pPzsY3eXSkCq7lY2jB/VQU7lfrxbcHIrj+X/lns=";
        };

        devTools = {
          pggen = pkgs.buildGoModule {
            pname = "pggen";
            src = pggen;
            subPackages = [ "cmd/pggen" ];
            doCheck = false;
            vendorSha256 = "sha256-WLoFpwOP97160WfmfbCUUlhqGC0qiEPWDg0qL/DrzIA=";
          };
          dbmate = pkgs.buildGoModule {
            pname = "dbmate";
            src = dbmate;
            subPackages = [ "." ];
            doCheck = false;
            vendorSha256 = "sha256-U9VTS0rmLHxweFiIcFyoybHMBihy5ezloDC2iLc4IMc=";
          };
          overmind = pkgs.buildGoModule {
            pname = "overmind";
            src = overmind;
            subPackages = [ "." ];
            dbCheck = false;
            vendorSha256 = "sha256-KDMzR6qAruscgS6/bHTN6RnHOlLKCm9lxkr9k3oLY+Y=";
          };
        };
      in
        rec {
          packages = flake-utils.lib.flattenTree {
            tstr = tstr;
          };
          defaultPackage = packages.tstr;

          apps = {
            tstr = flake-utils.lib.mkApp {
              drv = packages.tstr;
              exePath = "/bin/tstr";
            };
            tstrctl = flake-utils.lib.mkApp {
              drv = packages.tstr;
              exePath = "/bin/tstrctl";
            };
          };
          defaultApp = apps.tstr;

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
