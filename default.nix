# For more, refer to:
# https://github.com/NixOS/nixpkgs/blob/master/doc/languages-frameworks/go.section.md
{pkgs ? import <nixpkgs> {}}: let
  lib = pkgs.lib;
  getLibFolder = pkg: "${pkg}/lib";
in
  pkgs.buildGoModule rec {
    pname = "gobot";
    version = "0.1.0";

    src = pkgs.lib.cleanSource ./.;

    vendorHash = "sha256-SJ5J7X0Q6IAylyCZP3wjG/TwxrYWsCgSXsKbej5OIDg=";

    meta = with lib; {
      homepage = "https://gopher.uz";
      description = "Telegram bot of Go Uzbekistan community";
      license = with lib.licenses; [mit];
      platforms = with platforms; linux ++ darwin;
      maintainers = [lib.teams.uzinfocom.bahrom04];
    };
  }
