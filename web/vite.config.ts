import http from "node:http";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

// 启用 Keep-Alive 的 HTTP Agent，避免 http-proxy 默认发送 Connection: close
const keepAliveAgent = new http.Agent({ keepAlive: true });

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    host: true,
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8088",
        changeOrigin: true,
        // 允许跨域携带 Cookie
        cookieDomainRewrite: "localhost",
        // 后端路由也是 /api 开头，所以不需要 rewrite，直接原样转发
        agent: keepAliveAgent,
      },
    },
  },
});
