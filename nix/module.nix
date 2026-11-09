self:
{ config, lib, pkgs, ... }:
let cfg = config.programs.ryomarket;
in {
  options.programs.ryomarket = {
    enable = lib.mkEnableOption "RyoMarket application discovery and management";
    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.ryomarket;
      description = "RyoMarket package to install.";
    };
    edition = lib.mkOption {
      type = lib.types.enum [ "default" "hyperion" ];
      default = "default";
      description = ''
        Catalog edition. "hyperion" sets RYOMARKET_EDITION=hyperion so the Ryoku section installs
        from github:Ryoku-dev/Ryoku-on-NixOS/feature/hyperion-v2-nixos.
      '';
    };
  };
  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
    environment.sessionVariables = lib.mkIf (cfg.edition == "hyperion") { RYOMARKET_EDITION = "hyperion"; };
  };
}
