import { Container, ThemeProvider, ViewProvider, setEnvVars } from "@airplane/views";
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./src/{{.Entrypoint}}";

setEnvVars(
  import.meta.env.AIRPLANE_API_HOST || "https://api.airplane.dev",
  import.meta.env.AIRPLANE_TOKEN,
  import.meta.env.AIRPLANE_API_KEY,
  import.meta.env.AIRPLANE_ENV_SLUG,
);
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <ViewProvider>
        <Container size="xl" py={96}>
          <App />
        </Container>
      </ViewProvider>
    </ThemeProvider>
  </React.StrictMode>
);
