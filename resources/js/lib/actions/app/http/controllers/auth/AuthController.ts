import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
export const Index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Index.url(options),
  method: 'get',
});

Index.definition = {
  methods: ['get', 'head'],
  url: '/auth/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
Index.url = (options?: RouteQueryOptions) => {
  return Index.definition.url + queryParams(options);
};

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
Index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Index.url(options),
  method: 'get',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
Index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Index.url(options),
  method: 'head',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
export const IndexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Index.url(options),
  method: 'get',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
IndexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Index.url(options),
  method: 'get',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:36
 * @route /auth/
 */
IndexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Index.url(options),
  method: 'head',
});
/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
export const Login = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Login.url(options),
  method: 'post',
});

Login.definition = {
  methods: ['post'],
  url: '/auth/login',
} satisfies RouteDefinition<['post']>;

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
Login.url = (options?: RouteQueryOptions) => {
  return Login.definition.url + queryParams(options);
};

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
Login.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Login.url(options),
  method: 'post',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
export const LoginForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Login.url(options),
  method: 'post',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:40
 * @route /auth/login
 */
LoginForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Login.url(options),
  method: 'post',
});
/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
export const Logout = (options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Logout.url(options),
  method: 'delete',
});

Logout.definition = {
  methods: ['delete'],
  url: '/auth/logout',
} satisfies RouteDefinition<['delete']>;

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
Logout.url = (options?: RouteQueryOptions) => {
  return Logout.definition.url + queryParams(options);
};

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
Logout.delete = (options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Logout.url(options),
  method: 'delete',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
export const LogoutForm = (options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Logout.url(options),
  method: 'delete',
});

/**
 * @see auth/internal/app/http/controllers/auth/auth.go:84
 * @route /auth/logout
 */
LogoutForm.delete = (options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Logout.url(options),
  method: 'delete',
});

export const AuthController = {
  Index: Object.assign(Index, Index),
  Login: Object.assign(Login, Login),
  Logout: Object.assign(Logout, Logout),
};

export default AuthController;
