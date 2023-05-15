{
  description = "tstr devleopment";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nix-npm-buildpackage.url = "github:serokell/nix-npm-buildpackage";
    sqlc = {
      url = "github:kyleconroy/sqlc/v1.18.0";
      flake = false;
    };
    dbmate = {
      url = "github:amacneil/dbmate/v2.3.0";
      flake = false;
    };
    overmind = {
      url = "github:DarthSim/overmind/v2.4.0";
      flake = false;
    };
    grpc-gateway = {
      url = "github:grpc-ecosystem/grpc-gateway/v2.15.2";
      flake = false;
    };
    protoc-gen-grpc-gateway-ts = {
      url = "github:grpc-ecosystem/protoc-gen-grpc-gateway-ts/v1.1.2";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, nix-npm-buildpackage, sqlc, dbmate, overmind, grpc-gateway, protoc-gen-grpc-gateway-ts }:
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
            vendorSha256 = "sha256-gDePB+IZSyVIILDAj+O0Q8hgL0N/0Mwp1Xsrlh3B914=";
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
            vendorSha256 = "sha256-m1Nbu1bE04iOXnxW5kJfI9W95FU87eRKkOzg+YVvRsg=";
          };
          overmind = pkgs.buildGoModule {
            name = "overmind";
            src = overmind;
            doCheck = false;
            vendorSha256 = "sha256-ndgnFBGtVFc++h+EnA37aY9+zNsO5GDrTECA4TEWPN4==";
          };
          grpc-gateway = pkgs.buildGoModule {
            name = "grpc-gateway";
            src = grpc-gateway;
            doCheck = false;
            subPackages = [
              "protoc-gen-grpc-gateway"
              "protoc-gen-openapiv2"
            ];
            vendorSha256 = "sha256-WHwZ0EAJIz7mZr27x+Z7PKLLAiw1z2rQvvNynpMJQDw=";
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
                mockgen
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
