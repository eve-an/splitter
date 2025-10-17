import createClient from "openapi-fetch";
import type { paths } from "./swagger/schema";

const apiBaseUrl =
  import.meta.env?.VITE_API_BASE_URL?.toString().replace(/\/$/, "") ??
  "http://localhost:9080";

export const apiClient = createClient<paths>({
  baseUrl: apiBaseUrl,
});

export { apiBaseUrl };
