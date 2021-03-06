import env from "../environment";
import { getAccessToken } from "../utils/auth";

const get = async (url, options) =>
  await send(url, { method: "GET", ...options });

const post = async (url, body, options) =>
  await send(url, {
    method: "POST",
    body: JSON.stringify(body),
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  });

const postForm = async (url, formData, options) =>
  await send(url, {
    method: "POST",
    body: formData,
    ...options,
  });

const send = async (url, options) => {
  const dest = env.apiUrl + url;

  options.headers = {
    ...options.headers,
  };

  const accessToken = getAccessToken();
  if (accessToken) {
    options.headers.Authorization = "Bearer " + accessToken;
  }

  try {
    const res = await fetch(dest, options);

    switch (res.status) {
      case 200:
        let data = null;

        if (res.headers.get("Content-Type") === "application/json") {
          data = await res.json();
        }

        return {
          ok: true,
          data: data,
        };
      default:
        try {
          const error = await res.json();
          return {
            ok: false,
            error: error.message,
          };
        } catch (e) {
          console.error(
            `Failed to read response from ${dest} (status: ${res.status})`,
            e
          );
          return {
            ok: false,
            error: e.toString() || res.statusText,
          };
        }
    }
  } catch (e) {
    console.error(`Failed to make a ${options.method} request to ${dest}`, e);
    return {
      ok: false,
      error: e.toString(),
    };
  }
};

export { get, post, postForm };
