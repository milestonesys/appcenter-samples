// Client for the backend-for-frontend (BFF) served next to this MFE (see ../../app/src/main.py).
//
// No credentials are handled here. When the MFE runs inside the Smart Client shell, the shell
// adds `Authorization: Bearer <session access token>` to every fetch request, and the BFF
// forwards that token to the XProtect API Gateway.

declare const __webpack_public_path__: string;

// With `publicPath: 'auto'` this resolves to the directory remote-entry.js was loaded from,
// e.g. `/mfe/samples/restful-config-api/`, so the BFF is reachable at `<that>/api`.
const apiBase = `${__webpack_public_path__}api`;

export type ApiResult<T = unknown> = {
  ok: boolean;
  status: number;
  data: T;
};

export type UserDefinedEvent = {
  id: string;
  name: string;
  [key: string]: unknown;
};

async function request<T>(method: string, path: string, body?: unknown): Promise<ApiResult<T>> {
  const response = await fetch(`${apiBase}${path}`, {
    method,
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const text = await response.text();
  let data: unknown = text;
  try {
    data = text ? JSON.parse(text) : null;
  } catch {
    // Not JSON (e.g. an error page from a proxy); keep the raw text.
  }
  return { ok: response.ok, status: response.status, data: data as T };
}

export const getServerInfo = () => request<{ serverUrl: string }>('GET', '/server-info');

export const createEvent = (name: string) =>
  request<{ result: UserDefinedEvent }>('POST', '/user-defined-events', { name });

export const getEvent = (id: string) =>
  request<{ data: UserDefinedEvent }>('GET', `/user-defined-events/${encodeURIComponent(id)}`);

export const updateEvent = (id: string, name: string) =>
  request<{ data: UserDefinedEvent }>('PUT', `/user-defined-events/${encodeURIComponent(id)}`, { name });

export const deleteEvent = (id: string) =>
  request<unknown>('DELETE', `/user-defined-events/${encodeURIComponent(id)}`);
