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
    quicktemplate = {
      url = "github:valyala/quicktemplate";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, pggen, dbmate, overmind, quicktemplate }:
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

                postgresql_14

                devTools.pggen
                devTools.dbmate
                devTools.overmind
                devTools.quicktemplate
              ];
            };
        }
    );
}
