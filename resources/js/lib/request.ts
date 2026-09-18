import axios from 'axios';

const request = axios.create({
  baseURL: import.meta.env.VITE_APP_URL,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
  },
  withCredentials: true,
  timeout: 30000,
});

request.interceptors.response.use(null, requestErrorHandler);

export function json(data: unknown, init?: ResponseInit) {
  return new Response(JSON.stringify(data), {
    ...init,
    headers: {
      'Content-Type': 'application/json',
    },
  });
}

function requestErrorHandler(error: unknown) {
  if (axios.isAxiosError(error)) {
    const data = {
      title: 'Oops, something went wrong!',
      message: error.message,
      statusCode: error.response?.status,
    };

    if (error.response) {
      if (error.response.data.title) {
        data.title = error.response.data.title;
      }
      if (error.response.data.message) {
        data.message = error.response.data.message;
      }
    } else if (error.code === 'ERR_NETWORK') {
      data.title = 'Service unavailable';
      data.message = 'Service is unavailable at the moment, please try again later';
      data.statusCode = 503;
    } else if (error.code === 'ECONNABORTED') {
      data.title = 'Request timeout';
      data.message = 'Request timeout, please check your internet connection';
      data.statusCode = 408;
    }

    throw json(data, {
      status: data.statusCode,
    });
  }

  throw error;
}

export { request };
