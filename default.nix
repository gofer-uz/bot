# For more, refer to:
# https://github.com/NixOS/nixpkgs/blob/master/doc/languages-frameworks/rust.section.md
{pkgs ? import <nixpkgs> {}}: let
  lib = pkgs.lib;
  getLibFolder = pkg: "${pkg}/lib";
in
  pkgs.buildGoModule rec {
    pname = "gobot";
    version = "0.0.1";

    src = pkgs.lib.cleanSource ./.;

    vendorHash = "sha256-0GAeVDcuPuqjeLyFflBjMMBKcj/joy/ipwh5DGBAYi4=";

    meta = with lib; {
      homepage = "https://gopher.uz";
      description = "Telegram bot of Go Uzbekistan community";
      license = with lib.licenses; [mit];
      platforms = with platforms; linux ++ darwin;
      maintainers = [lib.maintainers.orzklv];
    };
  }
