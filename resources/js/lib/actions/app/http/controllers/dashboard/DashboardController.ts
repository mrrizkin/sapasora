import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
export const Index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Index.url(options),
  method: 'get',
});

Index.definition = {
  methods: ['get', 'head'],
  url: '/dashboard',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
Index.url = (options?: RouteQueryOptions) => {
  return Index.definition.url + queryParams(options);
};

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
Index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Index.url(options),
  method: 'get',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
Index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Index.url(options),
  method: 'head',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
export const IndexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Index.url(options),
  method: 'get',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
IndexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Index.url(options),
  method: 'get',
});

/**
 * @see dashboard/internal/app/http/controllers/dashboard/dashboard.go:24
 * @route /dashboard
 */
IndexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Index.url(options),
  method: 'head',
});

export const DashboardController = {
  Index: Object.assign(Index, Index),
};

export default DashboardController;
