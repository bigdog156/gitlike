class Gitlike < Formula
  desc "GitLike CLI with Git-like workflow for developers"
  homepage "https://github.com/bigdog156/gitlike"
  version "1.0.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-darwin-arm64"
    sha256 "7e4f3228da3555c21f11fc73b96f13e9ec9c7f1fcb42c7c3c418c0875250e196"
  else
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-darwin-amd64"
    sha256 "2509406fc062c00782b5a292eb104d2d321e4b638984b52fc9d5adb69f5d57f9"
  end

  def install
    bin.install "gitlike-darwin-arm64" => "gitlike" if Hardware::CPU.arm?
    bin.install "gitlike-darwin-amd64" => "gitlike" if Hardware::CPU.intel?
  end

  test do
    system "#{bin}/gitlike", "--help"
    assert_match "1.0.0", shell_output("#{bin}/gitlike --version 2>&1")
  end
end
