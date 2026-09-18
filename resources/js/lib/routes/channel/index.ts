import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see 
 * @route /ws/ping
 */
export const ping = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: ping.url(options),
  method: 'get',
});

ping.definition = {
  methods: ['get', 'head'],
  url: '/ws/ping',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see 
 * @route /ws/ping
 */
ping.url = (options?: RouteQueryOptions) => {
  return ping.definition.url + queryParams(options);
};

/**
 * @see 
 * @route /ws/ping
 */
ping.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: ping.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /ws/ping
 */
ping.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: ping.url(options),
  method: 'head',
});

/**
 * @see 
 * @route /ws/ping
 */
export const pingForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: ping.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /ws/ping
 */
pingForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: ping.url(options),
  method: 'get',
});

/**
 * @see 
 * @route /ws/ping
 */
pingForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: ping.url(options),
  method: 'head',
});

export const channel = {
  ping: Object.assign(ping, ping),
};

export default channel;
