class StudexCli < Formula
  desc "A command line interface for Studex platform"
  homepage "https://github.com/itshivams/Studex-CLI"
  url "https://github.com/itshivams/Studex-CLI/archive/refs/tags/v1.0.1.tar.gz"
  version "1.0.1"
  sha256 ""
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"studex-cli"), "main.go"
  end

  test do
    system "#{bin}/studex-cli", "--version"
  end
end
