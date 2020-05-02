class PygmyGo < Formula
  desc "Amazee.io's local development helper tool"
  homepage "https://github.com/fubarhouse/pygmy-go"
  url "https://github.com/fubarhouse/pygmy-go/releases/download/v0.2.0/pygmy-go-darwin"
  sha256 "8d7780ad9183d613313140aab0e088a61919e7efff69c5d0ab09005c2671d67c"
  version "v0.2.0"

  def install
    libexec.install Dir["*"]
    chmod(0755, "#{libexec}/pygmy-go-darwin")
    bin.install_symlink("#{libexec}/pygmy-go-darwin" => "pygmy")
  end
end
