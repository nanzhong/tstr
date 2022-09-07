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
    protoc-gen-grpc-gateway-ts = {
      url = "github:grpc-ecosystem/protoc-gen-grpc-gateway-ts/v1.1.2";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, nix-npm-buildpackage, sqlc, dbmate, overmind, quicktemplate, grpc-gateway, protoc-gen-grpc-gateway-ts }:
    flake-utils.lib.eachDefaultSystem(system:
      let
        pkgs = import nixpkgs { inherit system; };
        bp = pkgs.callPackage nix-npm-buildpackage {};
        version = "0.0.1";
        tstr = pkgs.buildGo119Module {
          pname = "tstr";
          inherit version;
          src = ./.;
          subPackages = [
            "cmd/tstr"
          ];
          vendorSha256 = null;
          buildInputs = [ tstr-ui ];
          preBuild = ''
            rm -rf ui/app/dist
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
          sqlc = pkgs.buildGoModule {
            name = "sqlc";
            src = sqlc;
            subPackages = [ "cmd/sqlc" ];
            doCheck = false;
            vendorSha256 = "sha256-QTNzhlphJNq918450WNXGDI/y6D0QHDoTwWYnH+NkbM=";
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
            vendorSha256 = "sha256-QIKyLknPvmt8yiUCSCIqha8h9ozDGeQnKSM9Vwus0uY=";
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
          protoc-gen-grpc-gateway-ts = pkgs.buildGoModule {
            name = "protoc-gen-grpc-gateway-ts";
            src = protoc-gen-grpc-gateway-ts;
            doCheck = false;
            vendorSha256 = "sha256-2Kytwh7jQulrrsYqHAsQPNFWFe3zIXuS4kza3mnkyDs=";
          };
        };
      in
        rec {
          packages = flake-utils.lib.flattenTree {
            tstr = tstr;
            image = pkgs.dockerTools.buildLayeredImage {
              name = "nanzhong/tstr";
              contents = [ tstr pkgs.busybox pkgs.cacert ];
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
                go_1_19
                gopls
                grpcurl
                kustomize
                nodePackages.vls
                postgresql_14
                protobuf
                protoc-gen-go
                protoc-gen-go-grpc
                protoc-gen-validate
                yarn
                nodejs

                devTools.sqlc
                devTools.dbmate
                devTools.overmind
                devTools.grpc-gateway
                devTools.protoc-gen-grpc-gateway-ts
              ];
            };
        }
    );
}
