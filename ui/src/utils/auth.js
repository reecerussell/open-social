const AUTH_STORAGE_KEY = "oa_auth";

const getAccessToken = () => {
  const authJson = localStorage.getItem(AUTH_STORAGE_KEY);
  if (!authJson) {
    return null;
  }

  try {
    const authData = JSON.parse(authJson);
    if (authData.expiryDate < new Date().getTime()) {
      localStorage.removeItem(AUTH_STORAGE_KEY);
      return null;
    }

    return authData.accessToken;
  } catch {
    localStorage.removeItem(AUTH_STORAGE_KEY);
    return null;
  }
};

const setAccessToken = (accessToken, expiryTimestamp) => {
  const expiryDateUtc = new Date(expiryTimestamp * 1000).toUTCString();
  const expiryDate = new Date(expiryDateUtc).getTime();

  const data = {
    accessToken,
    expiryDate,
  };

  localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(data));
};

const isAuthenticated = () => getAccessToken() !== null;

export { getAccessToken, setAccessToken, isAuthenticated };
