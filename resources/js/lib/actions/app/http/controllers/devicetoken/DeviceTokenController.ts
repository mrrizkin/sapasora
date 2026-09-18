import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
export const Destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

Destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/devicetoken/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
Destroy.url = (id: string, options?: RouteQueryOptions) => {
  return Destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
Destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
export const DestroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:167
 * @route /api/v1/devicetoken/:id
 */
DestroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});
/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
export const Get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

Get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/devicetoken/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
Get.url = (id: string, options?: RouteQueryOptions) => {
  return Get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
Get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
Get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Get.url(id, options),
  method: 'head',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
export const GetForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
GetForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:80
 * @route /api/v1/devicetoken/:id
 */
GetForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Get.url(id, options),
  method: 'head',
});
/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
export const List = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

List.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/devicetoken/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
List.url = (options?: RouteQueryOptions) => {
  return List.definition.url + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
List.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
List.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: List.url(options),
  method: 'head',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
export const ListForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
ListForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:50
 * @route /api/v1/devicetoken/
 */
ListForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: List.url(options),
  method: 'head',
});
/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
export const Store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

Store.definition = {
  methods: ['post'],
  url: '/api/v1/devicetoken/',
} satisfies RouteDefinition<['post']>;

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
Store.url = (options?: RouteQueryOptions) => {
  return Store.definition.url + queryParams(options);
};

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
Store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
export const StoreForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});

/**
 * @see devicetoken/internal/app/http/controllers/devicetoken/devicetoken.go:111
 * @route /api/v1/devicetoken/
 */
StoreForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});

export const DeviceTokenController = {
  Destroy: Object.assign(Destroy, Destroy),
  Get: Object.assign(Get, Get),
  List: Object.assign(List, List),
  Store: Object.assign(Store, Store),
};

export default DeviceTokenController;
