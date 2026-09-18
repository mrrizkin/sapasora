import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see 
 * @route /api/v1/health
 */
export const health = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: health.url(options),
  method: 'get',
});

health.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/health',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see 
 * @route /api/v1/health
 */
health.url = (options?: RouteQueryOptions) => {
  return health.definition.url + queryParams(options);
};

/**
 * @see 
 * @route /api/v1/health
 */
health.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: health.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /api/v1/health
 */
health.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: health.url(options),
  method: 'head',
});

/**
 * @see 
 * @route /api/v1/health
 */
export const healthForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: health.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /api/v1/health
 */
healthForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: health.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /api/v1/health
 */
healthForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: health.url(options),
  method: 'head',
});

export const api_v1 = {
  health: Object.assign(health, health),
};

export default api_v1;
