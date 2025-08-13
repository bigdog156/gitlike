class Gitlike < Formula
  desc "GitLike CLI with Git-like workflow for developers"
  homepage "https://github.com/bigdog156/gitlike"
  version "1.0.2"
  
  if Hardware::CPU.arm?
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.2/gitlike-darwin-arm64"
    sha256 "b735c7ff84540daaf8d5d748e80ba524199e8fa2b9110b0fe16b2b8d62c52a7a"
  else
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.2/gitlike-darwin-amd64"
    sha256 "2aa3093590ee62a5cd9fe451f8cde3b79e340edcf4d988508ee835352501f079"
  end

  def install
    bin.install "gitlike-darwin-arm64" => "gitlike" if Hardware::CPU.arm?
    bin.install "gitlike-darwin-amd64" => "gitlike" if Hardware::CPU.intel?
  end

  test do
    system "#{bin}/gitlike", "--help"
    assert_match "1.0.2", shell_output("#{bin}/gitlike --version 2>&1")
  end
end
