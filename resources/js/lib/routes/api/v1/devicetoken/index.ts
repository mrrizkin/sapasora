import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
export const destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/devicetoken/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
destroy.url = (id: string, options?: RouteQueryOptions) => {
  return destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
export const destroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
destroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});
/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
export const get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/devicetoken/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
get.url = (id: string, options?: RouteQueryOptions) => {
  return get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get.url(id, options),
  method: 'head',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
export const getForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
getForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
getForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get.url(id, options),
  method: 'head',
});
/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
export const list = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

list.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/devicetoken/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
list.url = (options?: RouteQueryOptions) => {
  return list.definition.url + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
list.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
list.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: list.url(options),
  method: 'head',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
export const listForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
listForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
listForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: list.url(options),
  method: 'head',
});
/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
export const store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

store.definition = {
  methods: ['post'],
  url: '/api/v1/devicetoken/',
} satisfies RouteDefinition<['post']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
store.url = (options?: RouteQueryOptions) => {
  return store.definition.url + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
export const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
storeForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});

export const api_v1_devicetoken = {
  destroy: Object.assign(destroy, destroy),
  get: Object.assign(get, get),
  list: Object.assign(list, list),
  store: Object.assign(store, store),
};

export default api_v1_devicetoken;
