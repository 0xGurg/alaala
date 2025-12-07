class Alaala < Formula
  desc "ᜀᜎᜀᜎ - Semantic memory system for AI assistants"
  homepage "https://github.com/0xGurg/alaala"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/0xGurg/alaala/releases/download/v0.1.0/alaala_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    else
      url "https://github.com/0xGurg/alaala/releases/download/v0.1.0/alaala_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/0xGurg/alaala/releases/download/v0.1.0/alaala_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    else
      url "https://github.com/0xGurg/alaala/releases/download/v0.1.0/alaala_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  depends_on "docker" => :recommended

  def install
    bin.install "alaala"
  end

  def post_install
    ohai "ᜀᜎᜀᜎ (alaala) installed successfully!"
    puts ""
    puts "Setup Weaviate for vector search:"
    puts "  docker run -d --name weaviate -p 8080:8080 \\"
    puts "    -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \\"
    puts "    -e PERSISTENCE_DATA_PATH=/var/lib/weaviate \\"
    puts "    -e DEFAULT_VECTORIZER_MODULE=none \\"
    puts "    weaviate/weaviate:latest"
    puts ""
    puts "Quick start:"
    puts "  1. Get an API key from OpenRouter: https://openrouter.ai"
    puts "     (Free tier available with meta-llama/llama-3.1-8b-instruct:free)"
    puts "  2. Initialize your project: alaala init"
    puts "  3. Configure for Cursor - see: https://github.com/0xGurg/alaala#quick-start"
    puts ""
    puts "Documentation: https://github.com/0xGurg/alaala"
  end

  def caveats
    <<~EOS
      ᜀᜎᜀᜎ requires Weaviate for vector search.
      
      Setup Weaviate with Docker:
        docker run -d --name weaviate -p 8080:8080 \\
          -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \\
          weaviate/weaviate:latest
    EOS
  end

  test do
    assert_match "alaala version", shell_output("#{bin}/alaala version")
  end
end

