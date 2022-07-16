{
  description = "tstr devleopment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    sqlc = {
      url = "github:kyleconroy/sqlc";
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
    quicktemplate = {
      url = "github:valyala/quicktemplate";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, sqlc, dbmate, overmind, quicktemplate }:
    flake-utils.lib.eachDefaultSystem(system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = "0.0.1";
        tstr = pkgs.buildGoModule {
          pname = "tstr";
          inherit version;
          src = ./.;
          subPackages = [
            "cmd/tstr"
            "cmd/tstrctl"
          ];
          vendorSha256 = null;
        };

        devTools = {
          sqlc = pkgs.buildGo118Module {
            name = "sqlc";
            src = sqlc;
            subPackages = [ "cmd/sqlc" ];
            doCheck = false;
            vendorSha256 = "sha256-mxDrO23FuoEi06Q0xvwKXVPpXDfB4HQzYPL2e6CtFIM=";
            proxyVendor = true;
            buildInputs = [
              pkgs.xxHash
              pkgs.libpg_query
              pkgs.postgresql_14
            ];
          };
          dbmate = pkgs.buildGoModule {
            name = "dbmate";
            src = dbmate;
            doCheck = false;
            vendorSha256 = "sha256-U9VTS0rmLHxweFiIcFyoybHMBihy5ezloDC2iLc4IMc=";
          };
          overmind = pkgs.buildGoModule {
            name = "overmind";
            src = overmind;
            dbCheck = false;
            vendorSha256 = "sha256-KDMzR6qAruscgS6/bHTN6RnHOlLKCm9lxkr9k3oLY+Y=";
          };
          quicktemplate = pkgs.buildGoModule {
            name = "quicktemplate";
            src = quicktemplate;
            dbCheck = false;
            vendorSha256 = null;
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
                protoc-gen-validate

                entr
                grpcurl

                postgresql_14

                devTools.sqlc
                devTools.dbmate
                devTools.overmind

                yarn
              ];
            };
        }
    );
}
