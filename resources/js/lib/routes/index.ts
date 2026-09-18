import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
export const dashboard = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: dashboard.url(options),
  method: 'get',
});

dashboard.definition = {
  methods: ['get', 'head'],
  url: '/dashboard',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
dashboard.url = (options?: RouteQueryOptions) => {
  return dashboard.definition.url + queryParams(options);
};

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
dashboard.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: dashboard.url(options),
  method: 'get',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
dashboard.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: dashboard.url(options),
  method: 'head',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
export const dashboardForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: dashboard.url(options),
  method: 'get',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
dashboardForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: dashboard.url(options),
  method: 'get',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
dashboardForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: dashboard.url(options),
  method: 'head',
});
/**
 * @see 
 * @route /
 */
export const index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: index.url(options),
  method: 'get',
});

index.definition = {
  methods: ['get', 'head'],
  url: '/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see 
 * @route /
 */
index.url = (options?: RouteQueryOptions) => {
  return index.definition.url + queryParams(options);
};

/**
 * @see 
 * @route /
 */
index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: index.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /
 */
index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: index.url(options),
  method: 'head',
});

/**
 * @see 
 * @route /
 */
export const indexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: index.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /
 */
indexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: index.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /
 */
indexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: index.url(options),
  method: 'head',
});
