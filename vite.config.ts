import { defineConfig, EnvironmentOptions } from "vite-plus";

// The UI ships two independent, self-contained scripts. In production the Go
// server inlines each bundle directly into a <script> tag (see ui/ui.go), so
// every entry must build to a standalone IIFE with no import/export or shared
// chunks. Rolldown only allows a single input per IIFE output, so each entry
// is its own build environment and `builder.buildApp` builds them in one pass.
export default defineConfig(({ mode }) => ({
  test: {
    globals: true,
    clearMocks: true,
    include: ["ui/assets/js/**/*.test.js"],
    coverage: {
      reportsDirectory: "coverage",
    },
  },
  fmt: {
    tabWidth: 4,
    ignorePatterns: ["ui/assets/js/JSXTransformer-0.12.2.js", "ui/assets/js/react-0.12.2.min.js"],
  },
  // JSX targets the global `React` (loaded separately via react-*.min.js),
  // using the classic runtime so no `react` import is injected.
  oxc: {
    jsx: {
      runtime: "classic",
      pragma: "React.createElement",
      pragmaFrag: "React.Fragment",
    },
  },
  environments: Object.fromEntries(
    ["hound", "excluded_files"].map((name) => [
      name,
      {
        build: {
          outDir: "ui/.build/ui/js",
          emptyOutDir: false,
          target: "es2015",
          minify: mode !== "development",
          lib: {
            entry: `ui/assets/js/${name}.jsx`,
            formats: ["iife"],
            name: "Hound",
            fileName: `${name}.js`,
          },
        },
      } satisfies EnvironmentOptions,
    ]),
  ),
  builder: {
    async buildApp(builder) {
      for (const [name, environment] of Object.entries(builder.environments)) {
        // Skip Vite's implicit default `client` environment; only our
        // explicitly-defined lib environments have an entry to build.
        if (name === "client") continue;
        await builder.build(environment);
      }
    },
  },
}));
