import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
export const index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: index.url(options),
  method: 'get',
});

index.definition = {
  methods: ['get', 'head'],
  url: '/auth/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
index.url = (options?: RouteQueryOptions) => {
  return index.definition.url + queryParams(options);
};

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: index.url(options),
  method: 'get',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: index.url(options),
  method: 'head',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
export const indexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: index.url(options),
  method: 'get',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
indexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: index.url(options),
  method: 'get',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
indexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: index.url(options),
  method: 'head',
});
/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
export const login = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: login.url(options),
  method: 'post',
});

login.definition = {
  methods: ['post'],
  url: '/auth/login',
} satisfies RouteDefinition<['post']>;

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
login.url = (options?: RouteQueryOptions) => {
  return login.definition.url + queryParams(options);
};

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
login.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: login.url(options),
  method: 'post',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
export const loginForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: login.url(options),
  method: 'post',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
loginForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: login.url(options),
  method: 'post',
});
/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
export const logout = (options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: logout.url(options),
  method: 'delete',
});

logout.definition = {
  methods: ['delete'],
  url: '/auth/logout',
} satisfies RouteDefinition<['delete']>;

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
logout.url = (options?: RouteQueryOptions) => {
  return logout.definition.url + queryParams(options);
};

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
logout.delete = (options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: logout.url(options),
  method: 'delete',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
export const logoutForm = (options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: logout.url(options),
  method: 'delete',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
logoutForm.delete = (options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: logout.url(options),
  method: 'delete',
});

export const auth = {
  index: Object.assign(index, index),
  login: Object.assign(login, login),
  logout: Object.assign(logout, logout),
};

export default auth;
