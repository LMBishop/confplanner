import { execSync } from "child_process";
let gitSha: string | null = null;
let version: string | null = null;
let revision: string | null = null;
try {
  gitSha = execSync("git rev-parse --short=7 HEAD").toString().trim();
  version = execSync("git log -1 --format=%cd --date=format:'%Y%m%d'").toString().trim();
  revision = execSync("git rev-list --count HEAD").toString().trim();
} catch (e) {
  gitSha = "unknown";
  revision = "0"
  version = new Date().toISOString().slice(0, 10).replace(/-/g, ".");
}

export default defineNuxtConfig({
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  ssr: true,
  css: ["~/assets/css/main.css"],

  runtimeConfig: {
    public: {
      baseURL: process.env.BASE_URL || "/api",
      gitSha: gitSha,
      version: version,
      revision: revision,
      versionString: gitSha === "unknown" ? "(unknown version)" : `v${version}.r${revision}.g${gitSha}`
    },
  },

  vite: {
    server: {
      proxy: {
        "/api": {
          target: "http://localhost:4000",
          changeOrigin: true,
        },
      },
    },
  },
  
  modules: ["@pinia/nuxt", "@vite-pwa/nuxt"],
});
