class Dbq < Formula
  desc "a database query tui"
  url "https://github.com/jshawl/dbq/archive/refs/tags/v2025.10.04.tar.gz"
  sha256 "266b7e6f7aa8fcad0058432cf5b4a43ecbf99dfac0784fead9a4db8112b92211"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"dbq"
  end
end