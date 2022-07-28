{
  description = "tstr devleopment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nix-npm-buildpackage.url = "github:serokell/nix-npm-buildpackage";
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
    grpc-gateway = {
      url = "github:grpc-ecosystem/grpc-gateway/v2.10.3";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, nix-npm-buildpackage, sqlc, dbmate, overmind, quicktemplate, grpc-gateway }:
    flake-utils.lib.eachDefaultSystem(system:
      let
        pkgs = import nixpkgs { inherit system; };
        bp = pkgs.callPackage nix-npm-buildpackage {};
        version = "0.0.1";
        tstr = pkgs.buildGo118Module {
          pname = "tstr";
          inherit version;
          src = ./.;
          subPackages = [
            "cmd/tstr"
          ];
          vendorSha256 = null;
          buildInputs = [ tstr-ui ];
          preBuild = ''
            cp -r "${tstr-ui}/dist" ui/app/dist
          '';
        };
        tstr-ui = bp.buildYarnPackage {
          pname = "tstr-ui";
          inherit version;
          src = ./ui/app;
          yarnBuildMore = "yarn vite build";
          postInstall = ''
            cp -r dist $out/
          '';
        };

        devTools = {
          sqlc = pkgs.buildGo118Module {
            name = "sqlc";
            src = sqlc;
            subPackages = [ "cmd/sqlc" ];
            doCheck = false;
            vendorSha256 = "sha256-0Q2HYP3am8H757wT8WqI+jglAuTkmysKPaZFKVQMYFo=";
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
            doCheck = false;
            vendorSha256 = "sha256-KDMzR6qAruscgS6/bHTN6RnHOlLKCm9lxkr9k3oLY+Y=";
          };
          quicktemplate = pkgs.buildGoModule {
            name = "quicktemplate";
            src = quicktemplate;
            doCheck = false;
            vendorSha256 = null;
          };
          grpc-gateway = pkgs.buildGoModule {
            name = "grpc-gateway";
            src = grpc-gateway;
            doCheck = false;
            subPackages = [
              "protoc-gen-grpc-gateway"
              "protoc-gen-openapiv2"
            ];
            vendorSha256 = "sha256-FhiTU9VmDZNCPBWrmCqmQo/kPdDe8Da1T2E06CVN2kw=";
          };
        };
      in
        rec {
          packages = flake-utils.lib.flattenTree {
            tstr = tstr;
            image = pkgs.dockerTools.buildLayeredImage {
              name = "nanzhong/tstr";
              tag = tstr.version;
              contents = [ tstr ];
              config = {
                Cmd = [ "/bin/tstr" ];
              };
            };

            default = tstr;
          };
          defaultPackage = packages.tstr;

          apps = {
            tstr = flake-utils.lib.mkApp {
              drv = packages.tstr;
              exePath = "/bin/tstr";
            };

            default = apps.tstr;
          };

          devShell = with pkgs;
            mkShell {
              buildInputs = [
                buf
                entr
                getopt
                go-tools
                go_1_18
                gopls
                grpcurl
                nodePackages.vls
                postgresql_14
                protobuf
                protoc-gen-go
                protoc-gen-go-grpc
                protoc-gen-validate
                yarn

                devTools.sqlc
                devTools.dbmate
                devTools.overmind
                devTools.grpc-gateway
              ];
            };
        }
    );
}
