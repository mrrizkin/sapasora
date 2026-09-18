import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
export const Destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

Destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/api-key/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
Destroy.url = (id: string, options?: RouteQueryOptions) => {
  return Destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
Destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
export const DestroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
DestroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
export const Get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

Get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/api-key/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
Get.url = (id: string, options?: RouteQueryOptions) => {
  return Get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
Get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
Get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Get.url(id, options),
  method: 'head',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
export const GetForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
GetForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
GetForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Get.url(id, options),
  method: 'head',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
export const List = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

List.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/api-key/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
List.url = (options?: RouteQueryOptions) => {
  return List.definition.url + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
List.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
List.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: List.url(options),
  method: 'head',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
export const ListForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
ListForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
ListForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: List.url(options),
  method: 'head',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
export const Store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

Store.definition = {
  methods: ['post'],
  url: '/api/v1/api-key/',
} satisfies RouteDefinition<['post']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
Store.url = (options?: RouteQueryOptions) => {
  return Store.definition.url + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
Store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
export const StoreForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
StoreForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
export const Update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

Update.definition = {
  methods: ['put'],
  url: '/api/v1/api-key/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
Update.url = (id: string, options?: RouteQueryOptions) => {
  return Update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
Update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
export const UpdateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
UpdateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});

export const APIKeyController = {
  Destroy: Object.assign(Destroy, Destroy),
  Get: Object.assign(Get, Get),
  List: Object.assign(List, List),
  Store: Object.assign(Store, Store),
  Update: Object.assign(Update, Update),
};

export default APIKeyController;
