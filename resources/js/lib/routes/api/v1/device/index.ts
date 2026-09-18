import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
export const destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/device/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
destroy.url = (id: string, options?: RouteQueryOptions) => {
  return destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
export const destroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
destroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
export const get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/device/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
get.url = (id: string, options?: RouteQueryOptions) => {
  return get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get.url(id, options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
export const getForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
getForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
getForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get.url(id, options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
export const get_device_by_token = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get_device_by_token.url(options),
  method: 'get',
});

get_device_by_token.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/device/by-token',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
get_device_by_token.url = (options?: RouteQueryOptions) => {
  return get_device_by_token.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
get_device_by_token.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get_device_by_token.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
get_device_by_token.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get_device_by_token.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
export const get_device_by_tokenForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get_device_by_token.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
get_device_by_tokenForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get_device_by_token.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
get_device_by_tokenForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get_device_by_token.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
export const list = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

list.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/device/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
list.url = (options?: RouteQueryOptions) => {
  return list.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
list.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
list.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: list.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
export const listForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
listForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
listForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: list.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
export const store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

store.definition = {
  methods: ['post'],
  url: '/api/v1/device/',
} satisfies RouteDefinition<['post']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
store.url = (options?: RouteQueryOptions) => {
  return store.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
export const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
storeForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
export const update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

update.definition = {
  methods: ['put'],
  url: '/api/v1/device/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
update.url = (id: string, options?: RouteQueryOptions) => {
  return update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
export const updateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
updateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
export const update_status = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update_status.url(id, options),
  method: 'put',
});

update_status.definition = {
  methods: ['put'],
  url: '/api/v1/device/:id/status',
} satisfies RouteDefinition<['put']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
update_status.url = (id: string, options?: RouteQueryOptions) => {
  return update_status.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
update_status.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update_status.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
export const update_statusForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update_status.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
update_statusForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update_status.url(id, options),
  method: 'put',
});

export const api_v1_device = {
  destroy: Object.assign(destroy, destroy),
  get: Object.assign(get, get),
  get_device_by_token: Object.assign(get_device_by_token, get_device_by_token),
  list: Object.assign(list, list),
  store: Object.assign(store, store),
  update: Object.assign(update, update),
  update_status: Object.assign(update_status, update_status),
};

export default api_v1_device;
