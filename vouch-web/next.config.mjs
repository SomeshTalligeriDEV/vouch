/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  outputFileTracingRoot: new URL("../../", import.meta.url).pathname,
  images: { remotePatterns: [{ protocol: "https", hostname: "**" }] },
};

export default nextConfig;
